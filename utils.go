package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// readFile reads the content of a file
func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return string(content), nil
}

// writeFile writes content to a file
func writeFile(filePath, content string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write the file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}
	return nil
}

// fileExists checks if a file exists
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// getLastCommand attempts to get the last executed command from shell history
func getLastCommand() (string, error) {
	// Try to get from different shell history files
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check different shell history files in order of preference
	historyFiles := []string{
		filepath.Join(homeDir, ".zsh_history"),
		filepath.Join(homeDir, ".bash_history"),
		filepath.Join(homeDir, ".history"),
	}

	for _, histFile := range historyFiles {
		if fileExists(histFile) {
			cmd, err := getLastCommandFromFile(histFile)
			if err == nil && cmd != "" {
				return cmd, nil
			}
		}
	}

	return "", fmt.Errorf("could not find shell history or last command")
}

// getLastCommandFromFile reads the last command from a history file
func getLastCommandFromFile(historyFile string) (string, error) {
	file, err := os.Open(historyFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lastLine string
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			// Handle zsh history format (timestamp:duration;command)
			if strings.Contains(line, ";") && strings.HasPrefix(line, ":") {
				parts := strings.SplitN(line, ";", 2)
				if len(parts) == 2 {
					line = parts[1]
				}
			}
			// Skip our own wdym commands
			if !strings.HasPrefix(line, "wdym") {
				lastLine = line
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return lastLine, nil
}

// executeCommand executes a shell command and returns the output
func executeCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// promptUser prompts the user for input and returns the response
func promptUser(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// confirmAction asks the user to confirm an action
func confirmAction(message string) bool {
	response, err := promptUser(fmt.Sprintf("%s (y/N): ", message))
	if err != nil {
		return false
	}
	response = strings.ToLower(response)
	return response == "y" || response == "yes"
}

// printDiff shows a simple diff between old and new content
func printDiff(filePath, oldContent, newContent string) {
	fmt.Printf("\n📝 Changes for %s:\n", filePath)
	fmt.Println("─────────────────────────────────────────────────────────────")

	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")

	additions := len(newLines) - len(oldLines)
	if additions > 0 {
		fmt.Printf("  %d additions (+)", additions)
	}
	if additions < 0 {
		fmt.Printf("  %d deletions (-)", -additions)
	}
	if additions == 0 {
		fmt.Printf("  Content modified")
	}
	
	fmt.Println()
	
	// Show a few sample changes if content is different
	if oldContent != newContent {
		fmt.Println("\nKey changes:")
		
		// Simple diff - show first few different lines
		maxLines := 5
		shown := 0
		
		for i := 0; i < len(newLines) && shown < maxLines; i++ {
			if i >= len(oldLines) {
				fmt.Printf("  + %s\n", newLines[i])
				shown++
			} else if oldLines[i] != newLines[i] {
				if oldLines[i] != "" {
					fmt.Printf("  - %s\n", oldLines[i])
				}
				fmt.Printf("  + %s\n", newLines[i])
				shown++
			}
		}
		
		if len(newLines) > maxLines || len(oldLines) > maxLines {
			fmt.Printf("  ... and more changes\n")
		}
	}
	
	fmt.Println("─────────────────────────────────────────────────────────────")
}

// getShellInfo returns information about the current shell and OS
func getShellInfo() (string, string) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "unknown"
	} else {
		// Extract just the shell name
		parts := strings.Split(shell, "/")
		shell = parts[len(parts)-1]
	}

	// Get OS info
	osInfo := "unknown"
	if cmd := exec.Command("uname", "-s"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			osInfo = strings.TrimSpace(string(output))
		}
	}

	return shell, osInfo
}

// colorize adds ANSI color codes to text
func colorize(text, color string) string {
	colors := map[string]string{
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"white":   "\033[37m",
		"reset":   "\033[0m",
	}

	if colorCode, exists := colors[color]; exists {
		return colorCode + text + colors["reset"]
	}
	return text
}

// truncateString truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
} 