package server

import (
	"bufio"
	"context"
	_ "errors"
	"fmt"
	"io"

	// "log" // Replace standard log
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	pb "devkit/pkg/api/gen/agent"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	units "github.com/docker/go-units"

	// Use zap logger
	"devkit/internal/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AgentServer implements the AgentService gRPC service.
type AgentServer struct {
	pb.UnimplementedAgentServiceServer                // Embed for forward compatibility
	dockerCli                          *client.Client // Add Docker client field
	grpcServer                         *grpc.Server   // Add field to hold the gRPC server instance
}

// SetGRPCServer allows the main function to inject the gRPC server instance.
func (s *AgentServer) SetGRPCServer(grpcSrv *grpc.Server) {
	s.grpcServer = grpcSrv
}

// NewAgentServer creates a new AgentServer.
func NewAgentServer() (*AgentServer, error) {
	logging.L().Info("Initializing Docker client...")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logging.L().Errorf("Failed to create Docker client: %v", err)
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Optional: Ping the Docker daemon to ensure connectivity
	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Short timeout for ping
	defer cancel()
	_, err = cli.Ping(pingCtx)
	if err != nil {
		logging.L().Errorf("Failed to ping Docker daemon: %v", err)
		cli.Close() // Close the client if ping fails
		return nil, fmt.Errorf("failed to connect to Docker daemon: %w", err)
	}

	logging.L().Info("Agent server initialized with Docker client.")
	return &AgentServer{dockerCli: cli}, nil
}

// Ping handles the Ping RPC call.
func (s *AgentServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	// Optionally, perform a Docker ping here as well
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second) // Quick ping check
	defer cancel()
	_, err := s.dockerCli.Ping(pingCtx)
	if err != nil {
		logging.L().Warnf("Ping RPC Warning: Docker daemon unreachable: %v", err)
		return nil, status.Errorf(codes.Unavailable, "Docker daemon unreachable: %v", err)
	}
	return &pb.PingResponse{Message: "Pong from devkitd"}, nil
}

// inspectAndVerifyDevkitContainer inspects the container and checks if it belongs to devkit.
// It returns the inspection details or a gRPC status error.
func (s *AgentServer) inspectAndVerifyDevkitContainer(ctx context.Context, containerIDOrName string) (*types.ContainerJSON, error) {
	inspectData, err := s.dockerCli.ContainerInspect(ctx, containerIDOrName)
	if err != nil {
		if client.IsErrNotFound(err) {
			logging.L().Infof("Container '%s' not found.", containerIDOrName)
			return nil, status.Errorf(codes.NotFound, "container '%s' not found", containerIDOrName)
		}
		logging.L().Errorf("Failed to inspect container '%s': %v", containerIDOrName, err)
		return nil, status.Errorf(codes.Internal, "failed to inspect container '%s': %v", containerIDOrName, err)
	}

	// Check the name prefix (remove leading slash)
	actualName := strings.TrimPrefix(inspectData.Name, "/")
	if !strings.HasPrefix(actualName, "devkit-") {
		logging.L().Warnf("Operation forbidden on non-devkit container '%s' (actual name: '%s')", containerIDOrName, actualName)
		return nil, status.Errorf(codes.PermissionDenied, "container '%s' is not managed by devkit (name: '%s')", containerIDOrName, actualName)
	}

	return &inspectData, nil // Return inspection details on success
}

// StartContainer handles the StartContainer RPC call using the Docker SDK.
func (s *AgentServer) StartContainer(ctx context.Context, req *pb.StartContainerRequest) (*pb.StartContainerResponse, error) {
	logging.L().Infof("Received StartContainer request for template: %s (Image: %s)", req.TemplateName, req.Image)

	containerName := fmt.Sprintf("devkit-%s", req.TemplateName)
	logging.L().Infof("Preparing Docker request for container: %s", containerName)

	// Environment variables (convert map to KEY=VALUE slice)
	envSlice := make([]string, 0, len(req.Environment))
	for k, v := range req.Environment {
		// Ensure keys are uppercase
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", strings.ToUpper(k), v))
	}
	logging.L().Debugf("Environment variables passed to Docker: %v", envSlice)

	// Port mappings (convert template format to Docker SDK format)
	portBindings := nat.PortMap{}
	exposedPortsSet := nat.PortSet{}
	for _, portMapping := range req.Ports {
		parts := strings.SplitN(portMapping, ":", 2)
		if len(parts) != 2 {
			logging.L().Warnf("Invalid port mapping format: %s, skipping", portMapping)
			continue
		}
		hostPort := parts[0]
		containerPortProto := parts[1] // e.g., "6379/tcp"

		// Split container port and protocol
		portParts := strings.SplitN(containerPortProto, "/", 2)
		if len(portParts) != 2 {
			logging.L().Warnf("Invalid container port/protocol format '%s', skipping", containerPortProto)
			continue
		}
		containerPortStr := portParts[0]
		proto := portParts[1]

		containerNatPort, err := nat.NewPort(proto, containerPortStr)
		if err != nil {
			// This should be less likely now, but keep the check
			logging.L().Warnf("Error creating nat.Port for '%s/%s': %v, skipping", proto, containerPortStr, err)
			continue
		}

		exposedPortsSet[containerNatPort] = struct{}{}
		portBindings[containerNatPort] = []nat.PortBinding{
			{HostIP: "0.0.0.0", HostPort: hostPort},
		}
	}

	// Volume mounts (Bind Mounts)
	binds := make([]string, 0, len(req.Volumes))
	for _, volumeMapping := range req.Volumes {
		parts := strings.SplitN(volumeMapping, ":", 2)
		if len(parts) != 2 {
			logging.L().Warnf("Invalid volume mapping format: %s, skipping", volumeMapping)
			continue
		}

		hostConfigPath := parts[0]
		containerPath := parts[1]

		wd, err := os.Getwd()
		if err != nil {
			logging.L().Warnf("Failed to get working directory: %v, skipping volume %s", err, volumeMapping)
			continue
		}
		devkitVolumesBaseDir := filepath.Join(wd, ".devkit", "volumes")
		volumeSubDir := filepath.Base(hostConfigPath)
		finalHostPath := filepath.Join(devkitVolumesBaseDir, volumeSubDir)

		if err := os.MkdirAll(finalHostPath, 0755); err != nil {
			logging.L().Warnf("Failed to create host directory '%s' for volume: %v, skipping", finalHostPath, err)
			continue
		}
		binds = append(binds, fmt.Sprintf("%s:%s", finalHostPath, containerPath))
	}

	// Ensure image is present
	logging.L().Infof("Ensuring image '%s' is present...", req.Image)
	_, _, err := s.dockerCli.ImageInspectWithRaw(ctx, req.Image)
	if err != nil {
		if client.IsErrNotFound(err) {
			logging.L().Infof("Image '%s' not found locally, pulling...", req.Image)
			reader, err := s.dockerCli.ImagePull(ctx, req.Image, image.PullOptions{})
			if err != nil {
				logging.L().Errorf("Failed to pull image '%s': %v", req.Image, err)
				return nil, status.Errorf(codes.Internal, "failed to pull image '%s': %v", req.Image, err)
			}
			defer reader.Close()
			io.Copy(io.Discard, reader) // Drain the pull output to ensure completion
			logging.L().Infof("Image '%s' pulled successfully.", req.Image)
		} else {
			// Other error inspecting image
			logging.L().Errorf("Failed to inspect image '%s': %v", req.Image, err)
			return nil, status.Errorf(codes.Internal, "failed to inspect image '%s': %v", req.Image, err)
		}
	} else {
		logging.L().Infof("Image '%s' found locally.", req.Image)
	}

	// Create container
	logging.L().Infof("Creating container '%s'...", containerName)
	contConfig := &container.Config{
		Image:        req.Image,
		Env:          envSlice,
		Cmd:          req.Command,
		ExposedPorts: exposedPortsSet,
		// TODO: Add other config: User, WorkingDir, Labels etc. from template
	}
	hostConfig := &container.HostConfig{
		Binds:        binds,
		PortBindings: portBindings,
		NetworkMode:  container.NetworkMode(req.GetNetworkMode()),
		RestartPolicy: container.RestartPolicy{
			Name: container.RestartPolicyMode(req.GetRestartPolicy()),
		},
		Resources: container.Resources{},
		// TODO: Add other host config if needed
	}

	// Apply resource limits if specified
	if req.Resources != nil {
		// Memory (needs parsing)
		if req.Resources.Memory != "" {
			memBytes, err := units.RAMInBytes(req.Resources.Memory)
			if err != nil {
				logging.L().Warnf("Warning: Invalid memory format '%s': %v, skipping memory limit.", req.Resources.Memory, err)
			} else {
				hostConfig.Resources.Memory = memBytes
			}
		}
		// CPU (needs parsing)
		if req.Resources.Cpu != "" {
			// Docker API expects CPU shares (relative weight) or specific core counts.
			// The API structure container.Resources uses NanoCPUs (1e9 NanoCPUs = 1 CPU core)
			// Let's parse simple float values like "0.5", "1.0", "2"
			cpuFloat, err := strconv.ParseFloat(req.Resources.Cpu, 64)
			if err != nil {
				logging.L().Warnf("Warning: Invalid CPU format '%s': %v, skipping CPU limit.", req.Resources.Cpu, err)
			} else {
				hostConfig.Resources.NanoCPUs = int64(cpuFloat * 1e9)
			}
		}
	}

	resp, err := s.dockerCli.ContainerCreate(ctx, contConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		// TODO: Check for conflict errors (container name already exists)
		logging.L().Errorf("Failed to create container '%s': %v", containerName, err)
		return nil, status.Errorf(codes.Internal, "failed to create container '%s': %v", containerName, err)
	}

	containerID := resp.ID
	logging.L().Infof("Container '%s' created successfully (ID: %s).", containerName, containerID[:12])

	// Start container
	logging.L().Infof("Starting container '%s' (ID: %s)...", containerName, containerID[:12])
	if err := s.dockerCli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		logging.L().Errorf("Failed to start container '%s': %v", containerName, err)
		// Attempt to clean up the created container if start fails?
		return nil, status.Errorf(codes.Internal, "failed to start container '%s': %v", containerName, err)
	}

	logging.L().Infof("Container '%s' (ID: %s) started successfully.", containerName, containerID[:12])

	// Return Response
	return &pb.StartContainerResponse{
		ContainerId: containerID,
		Message:     fmt.Sprintf("Container %s started successfully", containerName),
	}, nil
}

// ListContainers handles the ListContainers RPC call using the Docker client.
func (s *AgentServer) ListContainers(ctx context.Context, req *pb.ListContainersRequest) (*pb.ListContainersResponse, error) {
	logging.L().Info("Received ListContainers request")

	// List containers (including non-running ones for now, might filter later)
	containers, err := s.dockerCli.ContainerList(ctx, container.ListOptions{All: true}) // Use container.ListOptions
	if err != nil {
		logging.L().Errorf("Failed to list Docker containers: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list containers: %v", err)
	}

	logging.L().Infof("Found %d containers total", len(containers))

	respContainers := make([]*pb.ContainerInfo, 0)
	for _, container := range containers {
		// Filter containers by the 'devkit-' prefix in their name
		// Note: container.Names usually looks like ["/devkit-redis"]
		var containerName string
		if len(container.Names) > 0 {
			// Remove the leading slash
			containerName = strings.TrimPrefix(container.Names[0], "/")
		} else {
			// Should not happen for named containers, but handle defensively
			containerName = container.ID[:12] // Use short ID if no name
		}

		if !strings.HasPrefix(containerName, "devkit-") {
			continue // Skip containers not managed by us
		}

		// Format ports
		ports := make([]string, 0, len(container.Ports))
		for _, p := range container.Ports {
			// Format: "hostIP:hostPort->containerPort/protocol"
			// We want: "hostPort:containerPort/protocol" or just "containerPort/protocol" if not published
			portStr := fmt.Sprintf("%d/%s", p.PrivatePort, p.Type)
			if p.PublicPort > 0 {
				// Use PublicPort if the port is exposed to the host
				portStr = fmt.Sprintf("%d:%s", p.PublicPort, portStr)
			}
			ports = append(ports, portStr)
		}

		info := &pb.ContainerInfo{
			Id:     container.ID[:12], // Use short ID
			Name:   containerName,
			Image:  container.Image,
			Status: container.State, // e.g., "running", "exited"
			Ports:  ports,
		}
		respContainers = append(respContainers, info)
	}

	logging.L().Infof("Returning %d devkit-managed containers", len(respContainers))

	return &pb.ListContainersResponse{
		Containers: respContainers,
	}, nil
}

// StopContainer handles the StopContainer RPC call using the Docker client.
func (s *AgentServer) StopContainer(ctx context.Context, req *pb.StopContainerRequest) (*pb.StopContainerResponse, error) {
	containerIDOrName := req.GetId()
	logging.L().Infof("Received StopContainer request for: %s", containerIDOrName)

	if containerIDOrName == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID or name must be provided")
	}

	// Use Helper for Safety Check
	inspectData, err := s.inspectAndVerifyDevkitContainer(ctx, containerIDOrName)
	if err != nil {
		return nil, err // Return the gRPC status error from the helper
	}
	// --- End Safety Check ---

	actualName := strings.TrimPrefix(inspectData.Name, "/") // Get name from inspect data
	logging.L().Infof("Attempting to stop container '%s' (ID: %s)...", actualName, inspectData.ID[:12])

	// Stop the container (using default timeout)
	// Note: ContainerStopOptions can be used for custom timeout
	stopOptions := container.StopOptions{}
	if err := s.dockerCli.ContainerStop(ctx, inspectData.ID, stopOptions); err != nil {
		// Check if it's already stopped? Docker API might return an error or might not.
		// For simplicity, we report the error for now.
		logging.L().Warnf("Failed to stop container '%s': %v", actualName, err)
		return nil, status.Errorf(codes.Internal, "failed to stop container '%s': %v", actualName, err)
	}

	logging.L().Infof("Container '%s' stopped successfully.", actualName)

	return &pb.StopContainerResponse{
		Message: fmt.Sprintf("Container %s stopped successfully", actualName),
	}, nil
}

// RemoveContainer handles the RemoveContainer RPC call using the Docker client.
func (s *AgentServer) RemoveContainer(ctx context.Context, req *pb.RemoveContainerRequest) (*pb.RemoveContainerResponse, error) {
	containerIDOrName := req.GetId()
	forceRemove := req.GetForce()
	logging.L().Infof("Received RemoveContainer request for: %s (Force: %t)", containerIDOrName, forceRemove)

	if containerIDOrName == "" {
		return nil, status.Error(codes.InvalidArgument, "container ID or name must be provided")
	}

	// Use Helper for Safety Check
	inspectData, err := s.inspectAndVerifyDevkitContainer(ctx, containerIDOrName)
	if err != nil {
		// Special case for remove: if NotFound, treat as success? Or keep NotFound?
		// Let's keep NotFound for now for clarity.
		return nil, err // Return the gRPC status error from the helper
	}
	// --- End Safety Check ---

	actualName := strings.TrimPrefix(inspectData.Name, "/")
	logging.L().Infof("Attempting to remove container '%s' (ID: %s, Force: %t)...", actualName, inspectData.ID[:12], forceRemove)

	// Remove the container
	removeOptions := container.RemoveOptions{
		Force: forceRemove,
		// RemoveVolumes: false, // Default is false, consider adding to proto req
		// RemoveLinks: false, // Usually not needed
	}
	if err := s.dockerCli.ContainerRemove(ctx, inspectData.ID, removeOptions); err != nil {
		// Handle potential errors (e.g., trying to remove a running container without force)
		logging.L().Errorf("Failed to remove container '%s': %v", actualName, err)
		return nil, status.Errorf(codes.Internal, "failed to remove container '%s': %v", actualName, err)
	}

	logging.L().Infof("Container '%s' removed successfully.", actualName)

	return &pb.RemoveContainerResponse{
		Message: fmt.Sprintf("Container %s removed successfully", actualName),
	}, nil
}

// LogsContainer handles the LogsContainer RPC call using the Docker client.
func (s *AgentServer) LogsContainer(req *pb.LogsContainerRequest, stream pb.AgentService_LogsContainerServer) error {
	containerIDOrName := req.GetId()
	logging.L().Infof("Received LogsContainer request for: %s (Follow: %t, Tail: %s)",
		containerIDOrName, req.GetFollow(), req.GetTail())

	if containerIDOrName == "" {
		return status.Error(codes.InvalidArgument, "container ID or name must be provided")
	}

	ctx := stream.Context() // Use the stream's context

	// Use Helper for Safety Check
	inspectData, err := s.inspectAndVerifyDevkitContainer(ctx, containerIDOrName)
	if err != nil {
		return err // Return the gRPC status error from the helper
	}
	// --- End Safety Check ---

	actualName := strings.TrimPrefix(inspectData.Name, "/") // Get name from inspect data

	logOptions := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     req.GetFollow(),
		Tail:       req.GetTail(),       // e.g., "100" or "all"
		Timestamps: req.GetTimestamps(), // Show timestamps
	}

	logReader, err := s.dockerCli.ContainerLogs(ctx, inspectData.ID, logOptions)
	if err != nil {
		logging.L().Errorf("Failed to get logs for container '%s': %v", actualName, err)
		return status.Errorf(codes.Internal, "failed to get logs for '%s': %v", actualName, err)
	}
	defer logReader.Close()

	logging.L().Infof("Streaming logs for container '%s'...", actualName)

	// Stream the logs back line by line
	// Note: Docker log stream might multiplex stdout/stderr with a header.
	// For simplicity here, we read line by line, which might interleave them
	// but won't show the source stream (stdout/stderr) explicitly.
	// A more robust solution uses docker/pkg/stdcopy.StdCopy.
	scanner := bufio.NewScanner(logReader)
	for scanner.Scan() {
		if err := stream.Send(&pb.LogsContainerResponse{Line: scanner.Bytes()}); err != nil {
			logging.L().Warnf("Error sending log line to client for '%s': %v", actualName, err)
			// Client likely disconnected
			return status.Errorf(codes.Unavailable, "log stream failed: %v", err)
		}
		// Check if context is cancelled (client disconnected)
		if ctx.Err() != nil {
			logging.L().Infof("Client disconnected while streaming logs for '%s'.", actualName)
			return status.Error(codes.Canceled, "client disconnected")
		}
	}

	if err := scanner.Err(); err != nil {
		if err == io.EOF || ctx.Err() == context.Canceled { // Check for EOF or context cancellation
			logging.L().Infof("Finished streaming logs for container '%s'.", actualName)
		} else {
			logging.L().Errorf("Error reading logs for container '%s': %v", actualName, err)
			return status.Errorf(codes.Internal, "error reading log stream for '%s': %v", actualName, err)
		}
	}

	return nil // Success
}

// RefreshTemplates handles the RefreshTemplates RPC call.
// Currently a no-op as templates are loaded by the CLI, but useful for future agent-side caching.
func (s *AgentServer) RefreshTemplates(ctx context.Context, req *pb.RefreshTemplatesRequest) (*pb.RefreshTemplatesResponse, error) {
	logging.L().Info("Received RefreshTemplates request. (Currently a no-op)")
	// In the future, this could trigger the agent to clear any cached template info.
	return &pb.RefreshTemplatesResponse{
		Message: "Template refresh acknowledged by agent (currently no action taken).",
	}, nil
}

// Shutdown handles the Shutdown RPC call.
func (s *AgentServer) Shutdown(ctx context.Context, req *pb.ShutdownRequest) (*pb.ShutdownResponse, error) {
	logging.L().Info("Received Shutdown request via RPC.")

	// Initiate graceful stop of the gRPC server in a goroutine
	if s.grpcServer != nil {
		logging.L().Info("Initiating graceful shutdown of the gRPC server...")
		go s.grpcServer.GracefulStop()
	}

	// Signal the main process to terminate itself
	pid := os.Getpid()
	proc, err := os.FindProcess(pid)
	if err != nil {
		logging.L().Errorf("Error finding current process (PID: %d): %v", pid, err)
		// Proceed with shutdown response anyway, but log the error
		return &pb.ShutdownResponse{
			Message: "Shutdown initiated, but failed to signal process for termination.",
		}, status.Errorf(codes.Internal, "failed to find process %d: %v", pid, err)
	}

	logging.L().Infof("Sending SIGTERM to process PID: %d...", pid)
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		logging.L().Errorf("Error sending SIGTERM to process (PID: %d): %v", pid, err)
		// Proceed with shutdown response anyway
		return &pb.ShutdownResponse{
			Message: "Shutdown initiated, but failed to send termination signal.",
		}, status.Errorf(codes.Internal, "failed to send SIGTERM to process %d: %v", pid, err)
	}

	// Return success response immediately, the signal will cause the main loop to exit.
	return &pb.ShutdownResponse{
		Message: "Shutdown initiated successfully via RPC.",
	}, nil
}

// StopAgent handles the StopAgent RPC call.
func (s *AgentServer) StopAgent(ctx context.Context, req *pb.StopAgentRequest) (*pb.StopAgentResponse, error) {
	if s.grpcServer == nil {
		logging.L().Error("StopAgent called, but gRPC server instance is nil")
		return nil, status.Errorf(codes.Internal, "gRPC server instance not available")
	}

	go func() {
		// Run the stop/gracefulstop in a goroutine so the RPC can return immediately.
		if req.GetForce() {
			logging.L().Info("StopAgent requested: Performing immediate shutdown (grpcServer.Stop())...")
			s.grpcServer.Stop()
		} else {
			logging.L().Info("StopAgent requested: Performing graceful shutdown (grpcServer.GracefulStop())...")
			s.grpcServer.GracefulStop()
		}

		// After stopping/gracefully stopping the server, signal the main process to exit.
		pid := os.Getpid()
		logging.L().Infof("StopAgent: Signaling PID %d with SIGTERM to exit main process...", pid)
		proc, err := os.FindProcess(pid)
		if err != nil {
			logging.L().Errorf("StopAgent Error: Failed to find own process (PID: %d): %v", pid, err)
			// Cannot signal, main process might hang. Log error.
		} else {
			if err := proc.Signal(syscall.SIGTERM); err != nil {
				logging.L().Errorf("StopAgent Error: Failed to send SIGTERM to own process (PID: %d): %v", pid, err)
				// Main process might hang. Log error.
			}
		}
	}()

	// Return immediately
	msg := "Graceful shutdown initiated."
	if req.GetForce() {
		msg = "Immediate shutdown initiated."
	}
	return &pb.StopAgentResponse{
		Message: msg,
	}, nil
}
