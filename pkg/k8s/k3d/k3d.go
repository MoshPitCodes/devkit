package k3d

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// runCommand executes a command and returns its stdout, stderr, and error.
func runCommand(name string, args ...string) (string, string, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// CreateCluster creates a new k3d cluster.
func CreateCluster(name, memory string, nodes int, wait time.Duration) error {
	args := []string{"cluster", "create", name}
	if nodes > 1 {
		// k3d creates 1 server by default. Additional nodes are agents.
		args = append(args, "--agents", fmt.Sprintf("%d", nodes-1))
	} else if nodes <= 0 {
		// Ensure at least one server node
		nodes = 1
		// No --agents flag needed
	}

	if memory != "" {
		args = append(args, "--servers-memory", memory)
		// Note: Applying to agents too might need a separate flag in devkit or logic here
		// if nodes > 1 {
		// 	 args = append(args, "--agents-memory", memory)
		// }
	}

	if wait > 0 {
		args = append(args, "--wait", "--timeout", wait.String())
	}

	fmt.Printf("Running k3d command: k3d %s\n", strings.Join(args, " ")) // Debug output
	stdout, stderr, err := runCommand("k3d", args...)

	if stdout != "" {
		fmt.Println("--- k3d stdout ---")
		fmt.Print(stdout)
		fmt.Println("------------------")
	}
	if stderr != "" {
		fmt.Println("--- k3d stderr ---")
		fmt.Print(stderr)
		fmt.Println("------------------")
	}

	if err != nil {
		return fmt.Errorf("k3d cluster create failed: %w", err)
	}
	return nil
}

// DeleteCluster deletes a k3d cluster.
func DeleteCluster(name string) error {
	args := []string{"cluster", "delete", name}
	fmt.Printf("Running k3d command: k3d %s\n", strings.Join(args, " ")) // Debug output
	stdout, stderr, err := runCommand("k3d", args...)

	if stdout != "" {
		fmt.Println("--- k3d stdout ---")
		fmt.Print(stdout)
		fmt.Println("------------------")
	}
	if stderr != "" {
		fmt.Println("--- k3d stderr ---")
		fmt.Print(stderr)
		fmt.Println("------------------")
	}

	if err != nil {
		return fmt.Errorf("k3d cluster delete failed: %w", err)
	}
	return nil
}

// ListClusters lists k3d clusters.
// Returns the raw output from `k3d cluster list` for now.
func ListClusters() (string, error) {
	args := []string{"cluster", "list"}
	// fmt.Printf("Running k3d command: k3d %s\n", strings.Join(args, " ")) // Debug output
	stdout, stderr, err := runCommand("k3d", args...)

	if stderr != "" {
		fmt.Println("--- k3d stderr ---")
		fmt.Print(stderr)
		fmt.Println("------------------")
	}

	if err != nil {
		return "", fmt.Errorf("k3d cluster list failed: %w", err)
	}
	return stdout, nil
}

// UseContext sets the kubectl context to the k3d cluster.
// Note: k3d automatically updates kubeconfig on create/delete.
// This function essentially runs `kubectl config use-context k3d-<name>`
func UseContext(name string) error {
	k3dContextName := fmt.Sprintf("k3d-%s", name)
	args := []string{"config", "use-context", k3dContextName}
	fmt.Printf("Running kubectl command: kubectl %s\n", strings.Join(args, " ")) // Debug output
	stdout, stderr, err := runCommand("kubectl", args...)

	if stdout != "" {
		fmt.Print(stdout)
	}
	if stderr != "" {
		fmt.Println("--- kubectl stderr ---")
		fmt.Print(stderr)
		fmt.Println("------------------")
	}

	if err != nil {
		// Check if the error is context not found?
		if strings.Contains(stderr, "no context exists with the name") {
			return fmt.Errorf("context '%s' not found. Does k3d cluster '%s' exist?", k3dContextName, name)
		}
		return fmt.Errorf("kubectl config use-context failed: %w", err)
	}
	return nil
}
