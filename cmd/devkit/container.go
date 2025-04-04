package main

import (
	"context"
	"fmt"
	"io"

	"strings"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"devkit/internal/cliutil"
	"devkit/internal/logging"
	pb "devkit/pkg/api/gen/agent"
	"devkit/pkg/config"
)

// containerCmd represents the base command when called without any subcommands
var containerCmd = &cobra.Command{
	Use:     "container",
	Aliases: []string{"c", "cont"},
	Short:   "Manage development containers via the devkit agent",
	Long: `Provides commands to run, list, stop, remove, and view logs for
development containers defined in template files. Requires the devkitd agent to be running.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Root PersistentPreRun already handles config loading
		// We could add a check here if config.Cfg is nil, but rely on commands
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help() // Show help if no subcommand is given
	},
}

// containerRunCmd represents the container run command
var containerRunCmd = &cobra.Command{
	Use:   "run [template-name]",
	Short: "Run a container based on a template",
	Args:  cobra.ExactArgs(1), // Requires exactly one argument: the template name
	Run: func(cmd *cobra.Command, args []string) {
		templateName := args[0]
		logging.L().Infof("Requesting to run container from template: %s", templateName)

		// Ensure config is loaded before trying to load template
		if config.Cfg == nil {
			_, err := config.LoadConfig()
			if err != nil || config.Cfg == nil {
				logging.L().Fatalf("Failed to load configuration before loading template: %v", err)
			}
		}

		// 1. Load the specific container template
		logging.L().Infof("Loading template '%s' definition...", templateName)
		tmpl, err := config.LoadContainerTemplate(templateName)
		if err != nil {
			logging.L().Fatalf("Failed to load template '%s': %v", templateName, err)
		}
		logging.L().Info("Template loaded successfully.")
		logging.L().Infof("Hint: Image: %s", tmpl.Image)
		if len(tmpl.Ports) > 0 {
			logging.L().Infof("Hint: Ports: %v", tmpl.Ports)
		}
		if len(tmpl.Volumes) > 0 {
			logging.L().Infof("Hint: Volumes: %v", tmpl.Volumes)
		}
		if len(tmpl.Environment) > 0 {
			// Format env vars nicely
			envPairs := make([]string, 0, len(tmpl.Environment))
			for k, v := range tmpl.Environment {
				envPairs = append(envPairs, fmt.Sprintf("%s=%s", k, v))
			}
			logging.L().Infof("Hint: Environment: [%s]", strings.Join(envPairs, ", "))
		}

		// 2. Connect to the agent
		logging.L().Info("Connecting to devkit agent...")
		// Use the command's context which might have its own timeout if cobra supports it,
		// otherwise background context is fine as GetAgentClient applies its own timeout.
		ctx := cmd.Context()
		agentConn, err := cliutil.GetAgentClient(ctx)
		if err != nil {
			logging.L().Fatalf("%v", err) // GetAgentClient formats the error well
		}
		defer agentConn.Close()
		logging.L().Info("Connected to agent.")

		// 3. Prepare and send the gRPC request
		client := agentConn.Client
		req := &pb.StartContainerRequest{
			TemplateName:  templateName,
			Image:         tmpl.Image,
			Ports:         tmpl.Ports,
			Volumes:       tmpl.Volumes,
			Environment:   tmpl.Environment,
			Command:       tmpl.Command,
			NetworkMode:   tmpl.NetworkMode,
			RestartPolicy: tmpl.RestartPolicy,
		}
		if tmpl.Resources != nil {
			req.Resources = &pb.Resources{
				Memory: tmpl.Resources.Memory,
				Cpu:    tmpl.Resources.CPU, // Note the field name difference (Cpu vs CPU)
			}
		}

		// Use a longer timeout for the actual operation
		// TODO: Make this configurable? For now, use a fixed longer duration.
		opTimeout := 5 * time.Minute // Example: 5 minutes for potentially long pulls/starts
		opCtx, opCancel := context.WithTimeout(ctx, opTimeout)
		defer opCancel()

		logging.L().Infof("Sending start request to agent for '%s'...", templateName)
		resp, err := client.StartContainer(opCtx, req)
		if err != nil {
			logging.L().Fatalf("Error calling StartContainer: %v", err)
		}
		logging.L().Infof("Agent response: %s (ID: %s)", resp.Message, resp.ContainerId[:12])
	},
}

// containerListCmd represents the container list command
var containerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List devkit-managed containers",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logging.L().Info("Requesting list of devkit containers...")

		ctx := cmd.Context()
		agentConn, err := cliutil.GetAgentClient(ctx)
		if err != nil {
			logging.L().Fatalf("%v", err)
		}
		defer agentConn.Close()
		logging.L().Info("Connected to agent.")

		client := agentConn.Client
		req := &pb.ListContainersRequest{}

		// Use default CLI timeout for list operations
		opTimeout := time.Duration(config.Cfg.Cli.DefaultTimeoutSeconds) * time.Second
		opCtx, opCancel := context.WithTimeout(ctx, opTimeout)
		defer opCancel()

		resp, err := client.ListContainers(opCtx, req)
		if err != nil {
			logging.L().Fatalf("Error calling ListContainers: %v", err)
		}

		if len(resp.Containers) == 0 {
			logging.L().Info("No devkit-managed containers found.")
			return
		}

		logging.L().Infof("=== Found %d Devkit Container(s) ===", len(resp.Containers))
		// Simple table-like output
		fmt.Printf("%-15s %-25s %-30s %-15s %s\n", "ID", "NAME", "IMAGE", "STATUS", "PORTS")
		fmt.Printf("%-15s %-25s %-30s %-15s %s\n",
			strings.Repeat("-", 15),
			strings.Repeat("-", 25),
			strings.Repeat("-", 30),
			strings.Repeat("-", 15),
			strings.Repeat("-", 20)) // Adjust port width

		for _, c := range resp.Containers {
			statusStr := c.Status
			if strings.ToLower(c.Status) == "running" {
				// Apply color codes directly if needed, or omit
			} else if strings.Contains(strings.ToLower(c.Status), "exit") {
				// Apply color codes directly if needed, or omit
			}
			fmt.Printf("%-15s %-25s %-30s %-15s %s\n",
				c.Id,
				c.Name,
				c.Image,
				statusStr,
				strings.Join(c.Ports, ", "),
			)
		}
	},
}

// containerStopCmd represents the container stop command
var containerStopCmd = &cobra.Command{
	Use:   "stop [container-id or name]",
	Short: "Stop a running devkit container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		containerIDOrName := args[0]
		logging.L().Infof("Requesting to stop container: %s", containerIDOrName)

		ctx := cmd.Context()
		agentConn, err := cliutil.GetAgentClient(ctx)
		if err != nil {
			logging.L().Fatalf("%v", err)
		}
		defer agentConn.Close()
		logging.L().Info("Connected to agent.")

		client := agentConn.Client
		req := &pb.StopContainerRequest{Id: containerIDOrName}

		// Use default CLI timeout for stop operations
		opTimeout := time.Duration(config.Cfg.Cli.DefaultTimeoutSeconds) * time.Second
		opCtx, opCancel := context.WithTimeout(ctx, opTimeout)
		defer opCancel()

		logging.L().Infof("Sending stop request for '%s'...", containerIDOrName)
		resp, err := client.StopContainer(opCtx, req)
		if err != nil {
			logging.L().Fatalf("Error calling StopContainer: %v", err)
		}
		logging.L().Infof("Agent response: %s", resp.Message)
	},
}

// containerRemoveCmd represents the container remove command
var containerRemoveCmd = &cobra.Command{
	Use:     "remove [container-id or name]",
	Aliases: []string{"rm"},
	Short:   "Remove a stopped devkit container",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		containerIDOrName := args[0]
		force, _ := cmd.Flags().GetBool("force")
		logging.L().Infof("Requesting to remove container: %s (Force: %t)", containerIDOrName, force)

		ctx := cmd.Context()
		agentConn, err := cliutil.GetAgentClient(ctx)
		if err != nil {
			logging.L().Fatalf("%v", err)
		}
		defer agentConn.Close()
		logging.L().Info("Connected to agent.")

		client := agentConn.Client
		req := &pb.RemoveContainerRequest{Id: containerIDOrName, Force: force}

		// Use default CLI timeout for remove operations
		opTimeout := time.Duration(config.Cfg.Cli.DefaultTimeoutSeconds) * time.Second
		opCtx, opCancel := context.WithTimeout(ctx, opTimeout)
		defer opCancel()

		logging.L().Infof("Sending remove request for '%s'...", containerIDOrName)
		resp, err := client.RemoveContainer(opCtx, req)
		if err != nil {
			logging.L().Fatalf("Error calling RemoveContainer: %v", err)
		}
		logging.L().Infof("Agent response: %s", resp.Message)
	},
}

// containerLogsCmd represents the container logs command
var containerLogsCmd = &cobra.Command{
	Use:   "logs [container-id or name]",
	Short: "Fetch logs from a devkit container",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		containerIDOrName := args[0]
		follow, _ := cmd.Flags().GetBool("follow")
		tail, _ := cmd.Flags().GetString("tail")
		timestamps, _ := cmd.Flags().GetBool("timestamps")
		logging.L().Infof("Requesting logs for container: %s (Follow: %t, Tail: %s)", containerIDOrName, follow, tail)

		// Use background context as the stream itself controls the lifetime
		ctx := context.Background() // Streaming doesn't use command context timeout directly

		agentConn, err := cliutil.GetAgentClient(ctx)
		if err != nil {
			logging.L().Fatalf("%v", err)
		}

		logging.L().Info("Connected to agent.")

		client := agentConn.Client
		req := &pb.LogsContainerRequest{
			Id:         containerIDOrName,
			Follow:     follow,
			Tail:       tail,
			Timestamps: timestamps,
		}

		// The LogsContainer call itself doesn't need a short timeout,
		// it establishes a stream. Use the stream's context.
		stream, err := client.LogsContainer(ctx, req)
		if err != nil {
			agentConn.Close() // Close connection on error
			logging.L().Fatalf("Error starting log stream: %v", err)
		}

		logging.L().Infof("Streaming logs from agent for '%s'...", containerIDOrName)
		for {
			resp, err := stream.Recv()
			if err != nil {
				// Check for specific error types instead of string matching
				if err == io.EOF {
					logging.L().Infof("Log stream finished gracefully for %s.", containerIDOrName)
					break // End of stream
				}
				st, ok := status.FromError(err)
				if ok {
					// Check for common client-side or server-side stream termination codes
					if st.Code() == codes.Canceled || st.Code() == codes.Unavailable || st.Code() == codes.DeadlineExceeded {
						logging.L().Infof("Log stream terminated for %s (%s).", containerIDOrName, st.Code())
						break
					}
				}
				// Log other unexpected errors
				logging.L().Errorf("Error receiving log stream data: %v", err)
				break // Exit loop on other errors
			}
			// Print the raw bytes received
			fmt.Print(string(resp.Line))
			// Ensure newline if the log line didn't include one (Docker logs usually do)
			if len(resp.Line) > 0 && resp.Line[len(resp.Line)-1] != '\n' {
				fmt.Println()
			}
		}

		// Close the connection now that the stream is done
		agentConn.Close()
	},
}
