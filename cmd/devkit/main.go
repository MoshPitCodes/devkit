package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"devkit/internal/cliutil"     // Import the new helper package
	pb "devkit/pkg/api/gen/agent" // Need agent types
	"devkit/pkg/config"           // Import the new config package
	"devkit/pkg/k8s/k3d"          // Import the new k3d package

	"github.com/spf13/cobra"

	"devkit/internal/logging" // Import the new logging package
)

// Global flags
var debug bool

const devkitLogo = `
╔═══════════════════════════════╗
║        ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄        ║
║      ▄▀▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▀▄      ║
║    ▄▀▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▀▄    ║
║    █▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█    ║
║   █▒▒▒█   DEVKIT CLI  █▒▒▒█   ║
║   █▒▒▒█ VERSION 0.1.0 █▒▒▒█   ║
║   █▒▒▒█ @moshpitcodes █▒▒▒█   ║
║    █▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒█    ║
║     ▀▄▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▄▀     ║
║       ▀▄▄▀▀▀▀▀▀▀▀▀▀▀▄▄▀       ║
╚═══════════════════════════════╝
`

// Default config content
const defaultConfigContent = `# Default devkit Configuration
# This file allows you to customize devkit's behavior for this project.
# Settings here override global user settings (~/.config/devkit/config.yaml)
# and environment variables (DEVKIT_...).

# agent:
#   # Override the default Unix socket path for the agent
#   socketPath: /tmp/devkitd.sock
#   # Log level for the agent process
#   # Options: debug, info, warn, error
#   # logLevel: info

# cli:
#   # Default timeout in seconds for CLI commands connecting to the agent
#   defaultTimeoutSeconds: 60
#   # Log level for the CLI process
#   # Options: debug, info, warn, error
#   # logLevel: warn

# # Default location for container template YAML files
# templateDirectory: .devkit/templates

# --- Future Configuration Examples ---

# kubernetes:
#   # Default context to use for 'devkit k8s' commands
#   # context: docker-desktop
#   # Default namespace
#   # namespace: default
#   # Path to kubeconfig file (uses default resolution if not set)
#   # kubeconfig: ~/.kube/config

# database:
#   # Default connection details (can be overridden by templates or flags)
#   # defaultUser: myuser
#   # defaultPasswordEnvVar: MYAPP_DB_PASSWORD # Name of env var holding password

# observability:
#   # Default endpoint for pushing traces (e.g., Jaeger, Tempo)
#   # tracingEndpoint: http://localhost:4318/v1/traces
#   # Default endpoint for Prometheus metrics scraping
#   # metricsScrapeTarget: http://localhost:9090
`

// Simple function to ask for confirmation
func askForConfirmation(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/N]: ", prompt)
		response, err := reader.ReadString('\n')
		if err != nil {
			// Use logging for fatal errors if possible, but this might be too early
			// depending on when logging is initialized relative to init command.
			// Sticking to fmt for now as it's a direct user interaction.
			fmt.Fprintf(os.Stderr, "Error reading confirmation: %v\n", err)
			os.Exit(1) // Exit if we can't even read input
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" || response == "" {
			return false
		}
		// If input is invalid, loop again
	}
}

var rootCmd = &cobra.Command{
	Use:   "devkit",
	Short: "devkit - A Comprehensive Developer Toolkit",
	Long: `devkit provides a unified interface for various development tasks,
including container management, Kubernetes operations, database interactions,
and observability.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load configuration only if the command is NOT 'init'
		// because 'init' is responsible for creating the config
		if cmd.Name() != "init" {
			_, err := config.LoadConfig()
			if err != nil {
				// Only fail hard if config exists but is invalid.
				// If it doesn't exist, commands should handle config.Cfg being nil
				// or rely on defaults (though LoadConfig currently errors).
				// TODO: Refine LoadConfig to not error if file not found, maybe?
				// For now, init command bypasses this PersistentPreRun load. Consider only Warning?
				logging.L().Fatalf("Failed to load configuration: %v", err)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Print the logo when the root command is run without subcommands
		fmt.Print(devkitLogo)
		fmt.Println() // Add a newline after the logo
		cmd.Help()
	},
}

// --- initCmd Helper Functions ---

// ensureDir ensures a directory exists, creating it if necessary.
func ensureDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		logging.L().Infof("Creating directory: %s", dirPath)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dirPath, err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking directory %s: %w", dirPath, err)
	} else {
		logging.L().Infof("Directory %s already exists.", dirPath)
	}
	return nil
}

// handleConfigFile ensures the config file exists, creating or overwriting as needed.
func handleConfigFile(configFilePath string, forceOverwrite bool) error {
	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		logging.L().Infof("Creating default configuration file: %s", configFilePath)
		if err := os.WriteFile(configFilePath, []byte(defaultConfigContent), 0644); err != nil {
			return fmt.Errorf("error creating file %s: %w", configFilePath, err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking file %s: %w", configFilePath, err)
	} else { // Config file exists
		if forceOverwrite {
			if askForConfirmation(fmt.Sprintf("Overwrite existing configuration file %s?", configFilePath)) {
				logging.L().Infof("Overwriting default configuration file: %s", configFilePath)
				if err := os.WriteFile(configFilePath, []byte(defaultConfigContent), 0644); err != nil {
					return fmt.Errorf("error overwriting file %s: %w", configFilePath, err)
				}
			} else {
				logging.L().Info("Skipping overwrite of existing configuration file.")
			}
		} else {
			logging.L().Infof("Configuration file %s already exists.", configFilePath)
			logging.L().Infof("Hint: Use --force to overwrite with defaults.")
		}
	}
	return nil
}

// copyDevTemplates copies *.yaml files from devTemplatesDir to templatesDir.
func copyDevTemplates(devTemplatesDir, templatesDir string, forceOverwrite bool) error {
	devTemplatesExist := false
	if _, err := os.Stat(devTemplatesDir); err == nil {
		devTemplatesExist = true
		logging.L().Infof("Found development templates directory: %s", devTemplatesDir)
	} else if !os.IsNotExist(err) {
		logging.L().Warnf("Error checking directory %s: %v", devTemplatesDir, err)
		// Don't fail, just warn and proceed without copying
		return nil
	}

	if !devTemplatesExist {
		logging.L().Infof("Development templates directory %s not found, skipping copy.", devTemplatesDir)
		return nil
	}

	files, err := filepath.Glob(filepath.Join(devTemplatesDir, "*.yaml"))
	if err != nil {
		logging.L().Warnf("Error listing files in %s: %v", devTemplatesDir, err)
		return nil // Don't fail, just warn
	}

	if len(files) == 0 {
		logging.L().Infof("No *.yaml files found in %s to copy.", devTemplatesDir)
		return nil
	}

	logging.L().Infof("Found %d template(s) in %s.", len(files), devTemplatesDir)

	needsForceCheck := false
	templatesToCopy := make(map[string]string) // Map source path to dest path
	for _, srcPath := range files {
		fileName := filepath.Base(srcPath)
		destPath := filepath.Join(templatesDir, fileName)
		templatesToCopy[srcPath] = destPath
		if _, err := os.Stat(destPath); err == nil {
			needsForceCheck = true // At least one destination file exists
		}
	}

	proceedWithCopy := true
	if needsForceCheck && !forceOverwrite {
		logging.L().Infof("Existing templates found in %s.", templatesDir)
		logging.L().Infof("Hint: Use --force to overwrite.")
		proceedWithCopy = false
	} else if needsForceCheck && forceOverwrite {
		if !askForConfirmation(fmt.Sprintf("Overwrite existing templates in %s with content from %s?", templatesDir, devTemplatesDir)) {
			proceedWithCopy = false
			logging.L().Info("Skipping template overwrite.")
		}
	}

	if proceedWithCopy {
		logging.L().Infof("Copying templates from %s to %s...", devTemplatesDir, templatesDir)
		copiedCount := 0
		for srcPath, destPath := range templatesToCopy {
			content, err := os.ReadFile(srcPath)
			if err != nil {
				logging.L().Warnf("Failed to read %s: %v", srcPath, err)
				continue // Skip this file
			}
			if err := os.WriteFile(destPath, content, 0644); err != nil {
				logging.L().Warnf("Failed to write %s: %v", destPath, err)
				continue // Skip this file
			}
			copiedCount++
		}
		logging.L().Infof("Copied %d template(s).", copiedCount)
	}

	return nil
}

// checkUnknownEntries scans a directory for entries not in the expected list.
func checkUnknownEntries(dirPath string, expectedEntries map[string]bool) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		logging.L().Warnf("Error listing contents of %s: %v", dirPath, err)
		return nil // Don't fail the init process for this
	}

	for _, entry := range entries {
		if !expectedEntries[entry.Name()] {
			entryType := "file"
			if entry.IsDir() {
				entryType = "directory"
			}
			logging.L().Warnf("Found unexpected %s in %s: %s", entryType, dirPath, entry.Name())
		}
	}
	return nil
}

// initCmd represents the command to initialize the devkit structure
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize devkit structure in the current directory",
	Long: `Creates the .devkit directory and necessary subdirectories/configuration files if they don't exist.
Uses existing files/directories if found. Use --force to overwrite existing config/template files.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip the root PersistentPreRun config loading for the init command
	},
	Run: func(cmd *cobra.Command, args []string) {
		logging.L().Info("=== Initializing devkit ===")
		forceOverwrite, _ := cmd.Flags().GetBool("force")

		devkitDir := ".devkit"
		configFilePath := filepath.Join(devkitDir, "config.yaml")
		baseTemplatesDir := filepath.Join(devkitDir, "templates")
		volumesDir := filepath.Join(devkitDir, "volumes")

		// Define source and destination paths for each template domain
		templateDomains := map[string]struct {
			Source string
			Dest   string
		}{
			"containers": {Source: "templates/containers", Dest: filepath.Join(baseTemplatesDir, "containers")},
			"k8s":        {Source: "templates/k8s", Dest: filepath.Join(baseTemplatesDir, "k8s")},
			"db":         {Source: "templates/db", Dest: filepath.Join(baseTemplatesDir, "db")},
			"observe":    {Source: "templates/observe", Dest: filepath.Join(baseTemplatesDir, "observe")},
			"cicd":       {Source: "templates/cicd", Dest: filepath.Join(baseTemplatesDir, "cicd")},
		}

		// Step 1: Ensure .devkit directory exists
		if err := ensureDir(devkitDir); err != nil {
			logging.L().Fatalf("%v", err)
		}

		// Step 2: Ensure config file exists (create/overwrite if needed)
		if err := handleConfigFile(configFilePath, forceOverwrite); err != nil {
			logging.L().Fatalf("%v", err)
		}

		// Step 3: Ensure standard base subdirectories and domain template directories exist
		if err := ensureDir(baseTemplatesDir); err != nil { // Ensure .devkit/templates exists
			logging.L().Fatalf("%v", err)
		}
		for domain, paths := range templateDomains {
			if err := ensureDir(paths.Dest); err != nil { // Ensure .devkit/templates/<domain> exists
				logging.L().Fatalf("Failed to create template directory for %s: %v", domain, err)
			}
		}
		if err := ensureDir(volumesDir); err != nil {
			logging.L().Fatalf("%v", err)
		}

		// Step 4: Copy templates for each domain if source exists
		for domain, paths := range templateDomains {
			if err := copyDevTemplates(paths.Source, paths.Dest, forceOverwrite); err != nil {
				// copyDevTemplates only prints warnings, so no fatal error here
				logging.L().Errorf("Errors occurred during %s template copy: %v", domain, err) // Log if it somehow returns an error
			}
		}

		// Step 5: Check for unknown files/dirs in .devkit
		// We still expect 'templates' and 'volumes' directories directly under .devkit
		expected := map[string]bool{
			filepath.Base(configFilePath):   true,
			filepath.Base(baseTemplatesDir): true, // Check for 'templates'
			filepath.Base(volumesDir):       true,
		}
		if err := checkUnknownEntries(devkitDir, expected); err != nil {
			// checkUnknownEntries only prints warnings
		}

		logging.L().Info("devkit initialization complete.")

		// Step 6: Check and start the devkitd agent if not running
		logging.L().Info("Checking devkitd agent status...")
		if err := checkAndStartDaemon(); err != nil {
			// Provide a warning, not fatal, as core init succeeded
			logging.L().Warnf("Failed to start devkitd agent: %v", err)
			logging.L().Info("Hint: You may need to start the agent manually: run 'devkitd' in your terminal.")
		} else {
			logging.L().Info("devkitd agent is running or was started successfully.")
		}
	},
}

// checkAndStartDaemon attempts to connect to the agent. If it fails, it tries to start the daemon.
func checkAndStartDaemon() error {
	// Use a short timeout context for the check
	checkCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	agentConn, err := cliutil.GetAgentClient(checkCtx)
	if err == nil {
		// Connection successful, agent is running
		agentConn.Close() // Close the temporary connection
		logging.L().Info("devkitd agent is already running.")
		return nil
	}

	// Connection failed, log why (but proceed to attempt start)
	logging.L().Warnf("Failed to connect to agent (it may not be running): %v", err)
	logging.L().Info("Attempting to start the devkitd agent in the background...")

	// --- Find the devkitd executable ---
	var daemonPath string
	var findErr error
	// 1. Check PATH first
	daemonPath, findErr = exec.LookPath("devkitd")
	if findErr != nil {
		logging.L().Infof("'devkitd' not found in PATH, checking relative path...")
		// 2. If not in PATH, check relative to the current executable
		exePath, exeErr := os.Executable()
		if exeErr != nil {
			// If we can't even find our own path, give up here
			return fmt.Errorf("failed to get current executable path (%w), cannot find 'devkitd'", exeErr)
		}
		relDaemonPath := filepath.Join(filepath.Dir(exePath), "devkitd")
		if _, statErr := os.Stat(relDaemonPath); statErr == nil {
			// Found it relative to devkit executable
			daemonPath = relDaemonPath
			logging.L().Infof("Found 'devkitd' at %s", daemonPath)
			findErr = nil // Clear the LookPath error
		} else {
			// Still not found, return the original LookPath error plus context
			return fmt.Errorf("'devkitd' command not found in PATH (%v) or relative to the 'devkit' executable (%s: %v). Please install it or ensure it's accessible", findErr, relDaemonPath, statErr)
		}
	} else {
		logging.L().Infof("Found 'devkitd' in PATH at %s", daemonPath)
	}
	// --- End Find ---

	// Prepare the command to run devkitd
	cmd := exec.Command(daemonPath)
	// Set necessary attributes to run detached (basic approach)
	// TODO: Implement more robust platform-specific detachment (e.g., syscall.Setsid)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to execute '%s': %w", daemonPath, err)
	}

	logging.L().Infof("Started devkitd process with PID: %d", cmd.Process.Pid)
	logging.L().Infof("Hint: Allow a few seconds for the agent to initialize fully.")

	return nil
}

// refreshTemplatesCmd represents the command to signal the agent about template changes
var refreshTemplatesCmd = &cobra.Command{
	Use:   "refresh-templates",
	Short: "Notify the agent that container templates may have changed",
	Long:  `Signals the devkitd agent to acknowledge potential changes in container template files. Currently a no-op for the agent itself, as templates are loaded by the CLI on demand.`,
	Run: func(cmd *cobra.Command, args []string) {
		logging.L().Info("Attempting to signal agent to refresh templates...")

		// Use the command's context
		ctx := cmd.Context()

		agentConn, err := cliutil.GetAgentClient(ctx)
		if err != nil {
			// GetAgentClient provides a good error message, including the path.
			logging.L().Fatalf("%v\nHint: Is the devkitd agent running?", err)
		}
		defer agentConn.Close()
		logging.L().Info("Connected to devkitd agent.")

		client := agentConn.Client
		req := &pb.RefreshTemplatesRequest{}

		// Use the default CLI timeout for this operation
		opTimeout := time.Duration(config.Cfg.Cli.DefaultTimeoutSeconds) * time.Second
		opCtx, opCancel := context.WithTimeout(ctx, opTimeout)
		defer opCancel()

		logging.L().Info("Sending refresh request...")
		resp, err := client.RefreshTemplates(opCtx, req)
		if err != nil {
			logging.L().Fatalf("Error calling RefreshTemplates: %v", err)
		}

		logging.L().Infof("Agent response: %s", resp.Message)
	},
}

// k8sCmd represents the base command for Kubernetes operations
var k8sCmd = &cobra.Command{
	Use:     "k8s",
	Aliases: []string{"kubernetes"},
	Short:   "Manage local Kubernetes development clusters (using k3d)",
	Long: `Provides commands to create, delete, list, and manage local Kubernetes clusters
for development purposes, powered by k3d. Requires k3d and kubectl to be installed.
Docker must also be running.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Check prerequisites only if not the help command
		if cmd.Name() == "help" || cmd.Name() == "k8s" && len(args) == 0 { // Don't check for root k8s or help
			return nil
		}
		// Check if k3d exists
		if _, err := exec.LookPath("k3d"); err != nil {
			return fmt.Errorf("k3d command not found in PATH. Please install k3d: https://k3d.io")
		}
		// Check if kubectl exists (needed for 'config')
		if cmd.Name() == "config" {
			if _, err := exec.LookPath("kubectl"); err != nil {
				return fmt.Errorf("kubectl command not found in PATH: please install kubectl")
			}
		}
		// Check if Docker is running (implicitly checked by k3d, but a specific check might be good later)
		// TODO: Add Docker running check
		return nil
	},
}

// k8sCreateClusterCmd represents the k8s create cluster command
var k8sCreateClusterCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new local k3d cluster",
	Args:  cobra.MaximumNArgs(1), // Optional name, defaults later
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := "devkit-local" // Default name
		if len(args) > 0 {
			clusterName = args[0]
		}
		memory, _ := cmd.Flags().GetString("memory")
		nodes, _ := cmd.Flags().GetInt("nodes")
		wait, _ := cmd.Flags().GetDuration("wait")

		logging.L().Infof("Requesting to create k3d cluster '%s'...", clusterName)
		logging.L().Infof("Hint: Memory: %s, Nodes: %d, Wait: %v", memory, nodes, wait)

		if err := k3d.CreateCluster(clusterName, memory, nodes, wait); err != nil {
			logging.L().Fatalf("Failed to create cluster: %v", err)
		}
		logging.L().Infof("Successfully created k3d cluster '%s'.", clusterName)
		logging.L().Infof("Hint: Run 'devkit k8s config cluster %s' to set kubectl context.", clusterName)
	},
}

// k8sDeleteClusterCmd represents the k8s delete cluster command
var k8sDeleteClusterCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a local k3d cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]
		logging.L().Infof("Requesting to delete k3d cluster '%s'...", clusterName)

		if err := k3d.DeleteCluster(clusterName); err != nil {
			logging.L().Fatalf("Failed to delete cluster: %v", err)
		}
		logging.L().Infof("Successfully deleted k3d cluster '%s'.", clusterName)
	},
}

// k8sListClustersCmd represents the k8s list clusters command
var k8sListClustersCmd = &cobra.Command{
	Use:   "list",
	Short: "List local k3d clusters",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logging.L().Info("Requesting list of k3d clusters...")

		output, err := k3d.ListClusters()
		if err != nil {
			logging.L().Fatalf("Failed to list clusters: %v", err)
		}
		// Print the raw output from k3d for now
		fmt.Print(output)
	},
}

// k8sConfigClusterCmd represents the k8s config cluster command
var k8sConfigClusterCmd = &cobra.Command{
	Use:   "config <name>",
	Short: "Set kubectl context to use the specified local k3d cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]
		logging.L().Infof("Attempting to set kubectl context to k3d cluster '%s'...", clusterName)

		if err := k3d.UseContext(clusterName); err != nil {
			logging.L().Fatalf("Failed to set kubectl context: %v", err)
		}
		logging.L().Infof("kubectl context set to '%s'.", fmt.Sprintf("k3d-%s", clusterName))
	},
}

// agentCmd represents the command for agent operations
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage the devkitd agent",
	Long:  `Provides commands to manage the devkitd agent process.`,
}

// agentStopCmd represents the command to stop the devkitd agent
var agentStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running devkitd agent process",
	Long:  `Connects to the running devkitd agent via its socket and requests it to shut down.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		logging.L().Info("Attempting to stop the devkitd agent...")
		if force {
			logging.L().Warn("Using --force for immediate shutdown.")
		}

		ctx := cmd.Context()
		agentConn, err := cliutil.GetAgentClient(ctx)
		if err != nil {
			// If agent is not running, maybe consider this a success for 'stop'?
			// Or provide a specific hint.
			// GetAgentClient already provides a good hint about the agent not running.
			logging.L().Fatalf("%v", err)
		}
		defer agentConn.Close()
		logging.L().Info("Connected to agent.")

		client := agentConn.Client
		req := &pb.StopAgentRequest{Force: force}

		// Use a short timeout for the RPC call itself, the agent handles shutdown asynchronously
		opTimeout := 5 * time.Second // 5 seconds should be plenty for the RPC
		opCtx, opCancel := context.WithTimeout(ctx, opTimeout)
		defer opCancel()

		logging.L().Info("Sending stop request...")
		resp, err := client.StopAgent(opCtx, req)
		if err != nil {
			logging.L().Fatalf("Error calling StopAgent: %v", err)
		}

		logging.L().Infof("Agent response: %s", resp.Message)
		logging.L().Info("Hint: Agent process should terminate shortly.")
	},
}

func Execute() {
	// Initialize Logger based on global flags
	// Must be called before any logging occurs
	logging.Init(debug)

	if err := rootCmd.Execute(); err != nil {
		// Use the logger for fatal errors now
		logging.L().Fatalf("Error executing command: %v", err)
	}
}

func main() {
	Execute()
}

func init() {
	// Add global flags
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")

	// Add --force flag to initCmd
	initCmd.Flags().Bool("force", false, "Overwrite existing configuration and template files in .devkit with defaults/source templates.")

	// Add flags for k8s create cluster
	k8sCreateClusterCmd.Flags().String("memory", "", "Memory limit per server node (e.g., 512m, 1g). Passed to k3d.")
	k8sCreateClusterCmd.Flags().Int("nodes", 1, "Number of nodes (servers) in the cluster.")
	k8sCreateClusterCmd.Flags().Duration("wait", 0, "Wait duration for cluster components to become ready (e.g., 30s).")

	// Add flags for agent stop command
	agentStopCmd.Flags().Bool("force", false, "Force immediate shutdown instead of graceful shutdown.")

	// Add subcommands to the root command
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(containerCmd)
	rootCmd.AddCommand(refreshTemplatesCmd)
	rootCmd.AddCommand(agentCmd)
	rootCmd.AddCommand(k8sCmd)

	// Add agent subcommands
	// TODO: Add agent status command later if needed
	agentCmd.AddCommand(agentStopCmd)

	// Add container subcommands
	containerCmd.AddCommand(containerRunCmd)
	containerCmd.AddCommand(containerListCmd)
	containerCmd.AddCommand(containerStopCmd)
	containerCmd.AddCommand(containerRemoveCmd)
	containerCmd.AddCommand(containerLogsCmd)

	// Add k8s subcommands
	k8sCmd.AddCommand(k8sCreateClusterCmd)
	k8sCmd.AddCommand(k8sDeleteClusterCmd)
	k8sCmd.AddCommand(k8sListClustersCmd)
	k8sCmd.AddCommand(k8sConfigClusterCmd)

	// Add flags for container subcommands
	containerRemoveCmd.Flags().Bool("force", false, "Force the removal of a running container")
	containerLogsCmd.Flags().Bool("follow", false, "Follow log output")
	containerLogsCmd.Flags().String("tail", "all", "Number of lines to show from the end of the logs (e.g., '100')")
	containerLogsCmd.Flags().Bool("timestamps", false, "Show timestamps")
}
