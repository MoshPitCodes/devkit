package cliutil

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "devkit/pkg/api/gen/agent"
	"devkit/pkg/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AgentClientConnection holds the client and the connection itself for cleanup.
type AgentClientConnection struct {
	Client pb.AgentServiceClient
	Conn   *grpc.ClientConn
}

// Close closes the underlying gRPC connection.
func (acc *AgentClientConnection) Close() {
	if acc.Conn != nil {
		acc.Conn.Close()
	}
}

// GetAgentClient connects to the devkit agent via gRPC over its Unix socket.
// It handles loading configuration if necessary, applying a timeout to the connection attempt,
// and returning a client connection wrapper or an error.
// The caller is responsible for calling Close() on the returned AgentClientConnection.
func GetAgentClient(ctx context.Context) (*AgentClientConnection, error) {
	// Ensure configuration is loaded
	if config.Cfg == nil {
		_, err := config.LoadConfig()
		if err != nil || config.Cfg == nil {
			return nil, fmt.Errorf("devkit configuration not loaded or failed to load: %w", err)
		}
	}

	// Get socket path and default timeout from config
	sockPath := config.Cfg.Agent.SocketPath
	if sockPath == "" {
		// Should ideally use the constant from config package, but define here for now
		// to avoid dependency cycle if cliutil needed by config in future.
		sockPath = "/tmp/devkitd.sock"
	}
	connectTimeout := time.Duration(config.Cfg.Cli.DefaultTimeoutSeconds) * time.Second

	// Use the provided context, but apply a specific timeout for the Dial operation
	dialCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx,
		sockPath, // Use the socket path directly as the address for the custom dialer
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // Block until the connection is established or fails
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			// addr will be the sockPath passed to DialContext
			d := net.Dialer{}
			// Explicitly use "unix" network type
			return d.DialContext(ctx, "unix", addr)
		}),
	)

	if err != nil {
		// Provide a more helpful error message including the path
		return nil, fmt.Errorf("failed to connect to devkitd agent at %s: %w", sockPath, err)
	}

	// Connection successful
	client := pb.NewAgentServiceClient(conn)

	return &AgentClientConnection{
		Client: client,
		Conn:   conn,
	}, nil
}
