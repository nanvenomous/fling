package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

func main() {
	runQuery()
}

func runQuery() {
	// Get applications list
	apps, err := getExecutableApps()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting applications: %v\n", err)
		os.Exit(1)
	}

	if len(apps) == 0 {
		fmt.Fprintf(os.Stderr, "No applications found\n")
		os.Exit(1)
	}

	// Run fzf with the applications
	selectedApp, err := runFzf(apps)
	if err != nil {
		// User cancelled or error occurred
		os.Exit(1)
	}

	// Launch selected application
	launchApplication(selectedApp)
}

func runFzf(apps []string) (string, error) {
	// Start fzf process
	fzfCmd := exec.Command("fzf",
		"--height=20",
		"--reverse",
		"--prompt=> ",
		"--info=default",
		"--no-preview")

	// Set up pipes
	stdin, err := fzfCmd.StdinPipe()
	if err != nil {
		return "", err
	}

	stdout, err := fzfCmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	// Start fzf
	if err := fzfCmd.Start(); err != nil {
		return "", err
	}

	// Send application names to fzf
	go func() {
		defer stdin.Close()
		for _, app := range apps {
			fmt.Fprintln(stdin, app)
		}
	}()

	// Read selected item
	selectedBytes, err := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	// Wait for fzf to finish
	if err := fzfCmd.Wait(); err != nil {
		return "", err
	}

	return strings.TrimSpace(string(selectedBytes)), nil
}

func getExecutableApps() ([]string, error) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return nil, fmt.Errorf("PATH environment variable not set")
	}

	pathDirs := strings.Split(pathEnv, ":")
	appSet := make(map[string]bool)
	var apps []string

	for _, dir := range pathDirs {
		if dir == "" {
			continue
		}

		// Check if directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		// Read directory contents
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue // Skip directories we can't read
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			name := entry.Name()
			fullPath := filepath.Join(dir, name)

			// Check if file is executable
			if isExecutable(fullPath) && !appSet[name] {
				appSet[name] = true
				apps = append(apps, name)
			}
		}
	}

	// Sort applications alphabetically
	sort.Strings(apps)
	return apps, nil
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if it's a regular file and has execute permission
	mode := info.Mode()
	return mode.IsRegular() && (mode&0111) != 0
}

func launchApplication(appName string) {
	// Create command to execute the application
	cmd := exec.Command(appName)

	// Detach from parent process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Start the application
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to launch %s: %v\n", appName, err)
		os.Exit(1)
	}

	// Don't wait for the process to finish
	go func() {
		cmd.Wait()
	}()
}
