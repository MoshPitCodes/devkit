package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"devkit/internal/agent/server"
	"devkit/internal/logging"
	pb "devkit/pkg/api/gen/agent"
	"devkit/pkg/config"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// socketPath defines the location for the Unix domain socket.
// Consider using XDG runtime dir for better compatibility: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
// const socketPath = "/tmp/devkitd.sock" // Removed: Path is determined by config

func main() {
	// Initialize Logger - Assume not debug for daemon unless configured otherwise
	logging.Init(false) // TODO: Make debug configurable for daemon?

	// Replace standard log output with zap
	// This redirects log.Printf etc. calls to zap, useful for dependencies
	// that might use the standard logger.
	_ = zap.RedirectStdLog(logging.L().Desugar()) // Use Desugar() to get *zap.Logger
	log.SetFlags(0)                               // Disable standard log flags (prefix, time)

	logging.L().Info("Starting devkitd agent...")

	// Flag to indicate if shutdown was initiated by a signal
	var shutdownInitiated bool = false
	shutdownMutex := &sync.Mutex{}

	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logging.L().Fatalf("Failed to load configuration: %v", err)
	}
	if cfg == nil {
		logging.L().Fatal("Configuration is nil after loading.")
	}
	socketPath := cfg.Agent.SocketPath
	logging.L().Infof("Agent configured to listen on: %s", socketPath)

	// Attempt to Bind Socket (Singleton Check) - Listen First
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		// If listening failed, check if it's because the address is in use
		if isAddrInUseError(err) {
			logging.L().Errorf("Socket %s already in use. Another devkitd instance might be running.", socketPath)
			// Optional: Could try to connect to the socket here to verify it's a live agent
			os.Exit(1)
		}

		// If it failed for another reason, check if a stale socket file exists
		if _, statErr := os.Stat(socketPath); statErr == nil {
			logging.L().Warnf("Listening failed (%v) and socket file %s exists. Attempting to remove stale socket...", err, socketPath)
			if removeErr := os.Remove(socketPath); removeErr != nil {
				logging.L().Fatalf("Failed to remove stale socket %s after listen error: %v (original listen error: %v)", socketPath, removeErr, err)
			}
			// Try listening again after removing the stale file
			logging.L().Info("Retrying listener after removing stale socket...")
			lis, err = net.Listen("unix", socketPath)
			if err != nil {
				// If it still fails after removing the stale file, exit
				logging.L().Fatalf("Failed to listen on socket %s even after removing stale file: %v", socketPath, err)
			}
		} else {
			// Listening failed, and it wasn't EADDRINUSE, and the file doesn't exist.
			logging.L().Fatalf("Failed to listen on socket %s: %v", socketPath, err)
		}
	}

	// If we reach here, listening was successful (either initially or after cleanup)
	defer lis.Close()
	logging.L().Infof("Agent listening on %s", socketPath)

	// Create Agent Server Implementation
	s, err := server.NewAgentServer()
	if err != nil {
		logging.L().Fatalf("Failed to create agent server: %v", err)
	}

	// Create gRPC Server and Register Services
	grpcServer := grpc.NewServer()
	s.SetGRPCServer(grpcServer)
	pbsrv := &agentServerWrapper{AgentServer: s}
	pb.RegisterAgentServiceServer(grpcServer, pbsrv)

	// Register gRPC Health Checking Protocol server.
	healthcheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthcheck)
	healthcheck.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	logging.L().Info("Registered gRPC services and health check.")

	// Start gRPC Server in a Goroutine
	go func() {
		logging.L().Info("Starting gRPC server...")
		if err := grpcServer.Serve(lis); err != nil {
			shutdownMutex.Lock()
			initiated := shutdownInitiated
			shutdownMutex.Unlock()

			if !initiated {
				logging.L().Fatalf("gRPC server failed unexpectedly: %v", err)
			}
			logging.L().Info("gRPC server stopped.")
		}
	}()

	// Wait for Shutdown Signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	logging.L().Infof("Agent started. Waiting for signals (PID: %d)... Press Ctrl+C to exit.", os.Getpid())

	sig := <-quit
	logging.L().Infof("Received signal: %v. Shutting down gRPC server gracefully...", sig)

	// Graceful Shutdown
	shutdownMutex.Lock()
	shutdownInitiated = true
	shutdownMutex.Unlock()

	grpcServer.GracefulStop()
	logging.L().Info("gRPC server gracefully stopped.")

	logging.L().Info("Agent shutdown complete.")
}

// isAddrInUseError checks if the error is due to the address being in use.
func isAddrInUseError(err error) bool {
	opErr, ok := err.(*net.OpError)
	if !ok {
		return false
	}
	syscallErr, ok := opErr.Err.(*os.SyscallError)
	if !ok {
		return false
	}
	if strings.Contains(strings.ToLower(syscallErr.Error()), "address already in use") {
		return true
	}
	switch syscallErr.Err {
	case syscall.EADDRINUSE:
		return true
	}
	return false
}

// agentServerWrapper just embeds the AgentServer implementation.
type agentServerWrapper struct {
	*server.AgentServer // Embed the original server implementation
	// pb.UnimplementedAgentServiceServer // REMOVED to avoid ambiguity
}

// Ensure it still implements the interface (optional but good practice)
var _ pb.AgentServiceServer = (*agentServerWrapper)(nil)

// All methods are handled by the embedded AgentServer.
