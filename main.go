package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"syscall"

	fzf "github.com/junegunn/fzf/src"
	"gopkg.in/yaml.v3"
)

//go:embed version
var version string

type Config struct {
	Terminal struct {
		Command string   `yaml:"command"`
		Args    []string `yaml:"args"`
	} `yaml:"terminal"`
	TUIApps []string `yaml:"tui_apps"`
}

type UsageData struct {
	Counts map[string]int `yaml:"counts"`
}

var config *Config
var usage *UsageData

func main() {

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version", "-v":
			fmt.Printf("fling v%s\n", strings.TrimSpace(version))
			os.Exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
			os.Exit(1)
		}
	}

	loadConfig()
	loadUsage()
	runQuery()
}

func loadConfig() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		config = getDefaultConfig()
		return
	}

	configPath := filepath.Join(configDir, "fling", "config.yml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config = getDefaultConfig()
		createDefaultConfigFile(configPath)
		return
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		config = getDefaultConfig()
		return
	}

	config = &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		config = getDefaultConfig()
		return
	}
}

func getDefaultConfig() *Config {
	return &Config{
		Terminal: struct {
			Command string   `yaml:"command"`
			Args    []string `yaml:"args"`
		}{
			Command: "ghostty",
			Args:    []string{"-e"},
		},
		TUIApps: []string{
			"nvim",
			"vim",
			"nano",
			"htop",
			"lazygit",
		},
	}
}

func createDefaultConfigFile(configPath string) {
	configDir := filepath.Dir(configPath)
	os.MkdirAll(configDir, 0755)

	data, err := yaml.Marshal(getDefaultConfig())
	if err != nil {
		return
	}

	os.WriteFile(configPath, data, 0644)
}

func loadUsage() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		usage = &UsageData{Counts: make(map[string]int)}
		return
	}

	usagePath := filepath.Join(configDir, "fling", "usage.yml")

	if _, err := os.Stat(usagePath); os.IsNotExist(err) {
		usage = &UsageData{Counts: make(map[string]int)}
		return
	}

	data, err := os.ReadFile(usagePath)
	if err != nil {
		usage = &UsageData{Counts: make(map[string]int)}
		return
	}

	usage = &UsageData{}
	if err := yaml.Unmarshal(data, usage); err != nil {
		usage = &UsageData{Counts: make(map[string]int)}
		return
	}

	if usage.Counts == nil {
		usage.Counts = make(map[string]int)
	}
}

func saveUsage() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return
	}

	usagePath := filepath.Join(configDir, "fling", "usage.yml")
	configDir = filepath.Dir(usagePath)
	os.MkdirAll(configDir, 0755)

	data, err := yaml.Marshal(usage)
	if err != nil {
		return
	}

	os.WriteFile(usagePath, data, 0644)
}

func recordUsage(appName string) {
	if usage == nil {
		return
	}

	usage.Counts[appName]++
	saveUsage()
}

func isTUIApplication(appName string) bool {
	if config == nil {
		return false
	}

	return slices.Contains(config.TUIApps, appName)
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

	// Record usage and launch selected application
	recordUsage(selectedApp)
	launchApplication(selectedApp)
}

func runFzf(apps []string) (string, error) {
	inputChan := make(chan string)
	outputChan := make(chan string)

	// Send applications to fzf
	go func() {
		defer close(inputChan)
		for _, app := range apps {
			inputChan <- app
		}
	}()

	// Collect selected application
	var selected string
	go func() {
		for s := range outputChan {
			selected = s
		}
	}()

	// Build fzf.Options
	options, err := fzf.ParseOptions(
		false, // don't load defaults from environment
		[]string{"--height=20", "--reverse", "--prompt=> ", "--info=default", "--no-preview"},
	)
	if err != nil {
		return "", err
	}

	// Set up input and output channels
	options.Input = inputChan
	options.Output = outputChan

	// Run fzf
	code, err := fzf.Run(options)
	if err != nil {
		return "", err
	}

	if code != fzf.ExitOk {
		return "", fmt.Errorf("fzf exited with code %d", code)
	}

	return selected, nil
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

	// Sort applications by usage frequency, then alphabetically
	sort.Slice(apps, func(i, j int) bool {
		countI := 0
		countJ := 0

		if usage != nil {
			countI = usage.Counts[apps[i]]
			countJ = usage.Counts[apps[j]]
		}

		// If usage counts are different, sort by usage (most used first)
		if countI != countJ {
			return countI > countJ
		}

		// If usage counts are the same, sort alphabetically
		return apps[i] < apps[j]
	})

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
	var cmd *exec.Cmd

	if isTUIApplication(appName) {
		// Launch in terminal emulator
		args := append(config.Terminal.Args, appName)
		cmd = exec.Command(config.Terminal.Command, args...)
	} else {
		// Launch directly
		cmd = exec.Command(appName)
	}

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
