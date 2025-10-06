package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

type AppItem struct {
	Name      string
	Command   string
	Desktop   string
	IsDesktop bool
}

func main() {
	// if len(os.Args) < 2 {
	// 	fmt.Fprintf(os.Stderr, "Usage: %s <command> [script]\n", os.Args[0])
	// 	fmt.Fprintf(os.Stderr, "Commands: query\n")
	// 	os.Exit(1)
	// }

	// command := os.Args[1]

	// switch command {
	// case "query":
	// 	script := "./apps.sh"
	// 	if len(os.Args) > 2 {
	// 		script = os.Args[2]
	// 	}
	// 	runQuery(script)
	// default:
	// 	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
	// 	os.Exit(1)
	// }

	script := "./apps.sh"
	runQuery(script)
}

func runQuery(script string) {
	// Get applications list
	apps, err := getApplications(script)
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

func runFzf(apps []AppItem) (AppItem, error) {
	// Start fzf process
	fzfCmd := exec.Command("fzf",
		"--height=20",
		"--reverse",
		"--prompt=Apps> ",
		"--info=default",
		"--no-preview")

	// Set up pipes
	stdin, err := fzfCmd.StdinPipe()
	if err != nil {
		return AppItem{}, err
	}

	stdout, err := fzfCmd.StdoutPipe()
	if err != nil {
		return AppItem{}, err
	}

	// Start fzf
	if err := fzfCmd.Start(); err != nil {
		return AppItem{}, err
	}

	// Send application names to fzf
	go func() {
		defer stdin.Close()
		for _, app := range apps {
			fmt.Fprintln(stdin, app.Name)
		}
	}()

	// Read selected item
	selectedBytes, err := io.ReadAll(stdout)
	if err != nil {
		return AppItem{}, err
	}

	// Wait for fzf to finish
	if err := fzfCmd.Wait(); err != nil {
		return AppItem{}, err
	}

	// Find the selected application
	selectedName := strings.TrimSpace(string(selectedBytes))
	for _, app := range apps {
		if app.Name == selectedName {
			return app, nil
		}
	}

	return AppItem{}, fmt.Errorf("selected application not found")
}

func getApplications(script string) ([]AppItem, error) {
	// Make script path absolute
	absScript, err := filepath.Abs(script)
	if err != nil {
		return nil, err
	}

	// Run the script to get applications
	cmd := exec.Command(absScript, "executables")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run script: %v", err)
	}

	// Parse output
	var apps []AppItem
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Check if it's a desktop app (contains |)
		if strings.Contains(line, "|") {
			parts := strings.SplitN(line, "|", 3)
			if len(parts) >= 2 {
				apps = append(apps, AppItem{
					Name:      parts[0],
					Command:   parts[1],
					Desktop:   parts[2],
					IsDesktop: true,
				})
			}
		} else {
			// Regular executable
			apps = append(apps, AppItem{
				Name:      line,
				Command:   line,
				IsDesktop: false,
			})
		}
	}

	return apps, scanner.Err()
}

func launchApplication(app AppItem) {
	var cmd *exec.Cmd

	if app.IsDesktop {
		// Use the Exec line from .desktop file
		// Remove field codes like %f, %F, %u, %U
		execLine := app.Command
		execLine = strings.ReplaceAll(execLine, " %f", "")
		execLine = strings.ReplaceAll(execLine, " %F", "")
		execLine = strings.ReplaceAll(execLine, " %u", "")
		execLine = strings.ReplaceAll(execLine, " %U", "")
		execLine = strings.TrimSpace(execLine)

		// Split command and args
		parts := strings.Fields(execLine)
		if len(parts) == 0 {
			fmt.Fprintf(os.Stderr, "Empty command\n")
			return
		}

		cmd = exec.Command(parts[0], parts[1:]...)
	} else {
		// Regular executable
		cmd = exec.Command(app.Command)
	}

	// Detach from parent process
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// Start the application
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to launch %s: %v\n", app.Name, err)
		os.Exit(1)
	}

	// Don't wait for the process to finish
	go func() {
		cmd.Wait()
	}()
}
