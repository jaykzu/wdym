package main

import (
	"context"
	"fmt"
	"strings"
)

// handleShellMode generates shell commands from natural language
func handleShellMode(ctx context.Context, ai *GeminiClient, query string) error {
	if query == "" {
		return fmt.Errorf("no query provided for shell mode")
	}

	fmt.Printf("🐚 %s\n", colorize("Generating shell command...", "blue"))
	
	command, err := ai.GenerateShellCommand(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to generate shell command: %w", err)
	}

	command = strings.TrimSpace(command)
	fmt.Printf("\n%s %s\n", colorize("$", "green"), colorize(command, "white"))

	// Ask user if they want to execute
	if confirmAction("Execute this command?") {
		fmt.Printf("\n%s\n", colorize("Executing...", "yellow"))
		output, err := executeCommand(command)
		if err != nil {
			fmt.Printf("%s Command failed: %v\n", colorize("❌", "red"), err)
		}
		if output != "" {
			fmt.Printf("\n%s\n", output)
		}
	}

	return nil
}

// handleWriteMode creates or edits files using AI
func handleWriteMode(ctx context.Context, ai *GeminiClient, query, filePath string) error {
	if filePath == "" {
		return fmt.Errorf("no file path provided for write mode")
	}
	if query == "" {
		return fmt.Errorf("no instruction provided for write mode")
	}

	fmt.Printf("✏️  %s %s\n", colorize("Processing file:", "blue"), filePath)
	
	var oldContent string
	isNewFile := !fileExists(filePath)
	
	if !isNewFile {
		var err error
		oldContent, err = readFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read existing file: %w", err)
		}
	}

	fmt.Printf("📝 %s\n", colorize("Generating content...", "blue"))
	
	newContent, err := ai.EditFile(ctx, filePath, oldContent, query)
	if err != nil {
		return fmt.Errorf("failed to generate file content: %w", err)
	}

	newContent = strings.TrimSpace(newContent)

	// Show diff if editing existing file
	if !isNewFile {
		printDiff(filePath, oldContent, newContent)
	} else {
		fmt.Printf("\n📄 %s\n", colorize("New file content:", "green"))
		fmt.Printf("─────────────────────────────────────────────────────────────\n")
		preview := truncateString(newContent, 500)
		fmt.Printf("%s\n", preview)
		if len(newContent) > 500 {
			fmt.Printf("... (truncated, full content will be written to file)\n")
		}
		fmt.Printf("─────────────────────────────────────────────────────────────\n")
	}

	// Ask for confirmation
	action := "Create"
	if !isNewFile {
		action = "Update"
	}
	
	if confirmAction(fmt.Sprintf("%s file %s?", action, filePath)) {
		if err := writeFile(filePath, newContent); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("✅ %s %s\n", colorize("File saved:", "green"), filePath)
	} else {
		fmt.Printf("❌ %s\n", colorize("Operation cancelled", "red"))
	}

	return nil
}

// handleCodeMode generates code from natural language
func handleCodeMode(ctx context.Context, ai *GeminiClient, query string) error {
	if query == "" {
		return fmt.Errorf("no query provided for code mode")
	}

	fmt.Printf("💻 %s\n", colorize("Generating code...", "blue"))
	
	code, err := ai.GenerateCode(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	code = strings.TrimSpace(code)
	
	fmt.Printf("\n%s\n", colorize("Generated code:", "green"))
	fmt.Printf("─────────────────────────────────────────────────────────────\n")
	fmt.Printf("%s\n", code)
	fmt.Printf("─────────────────────────────────────────────────────────────\n")

	return nil
}

// handleFileQuery reads a file and answers questions about it
func handleFileQuery(ctx context.Context, ai *GeminiClient, query, filePath string) error {
	if filePath == "" {
		return fmt.Errorf("no file path provided")
	}
	if query == "" {
		query = "What does this file do?"
	}

	fmt.Printf("📄 %s %s\n", colorize("Reading file:", "blue"), filePath)
	
	content, err := readFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fmt.Printf("🤔 %s\n", colorize("Analyzing file...", "blue"))
	
	response, err := ai.AnalyzeFile(ctx, content, query)
	if err != nil {
		return fmt.Errorf("failed to analyze file: %w", err)
	}

	response = strings.TrimSpace(response)
	
	fmt.Printf("\n%s\n", colorize("Analysis:", "green"))
	fmt.Printf("─────────────────────────────────────────────────────────────\n")
	fmt.Printf("%s\n", response)
	fmt.Printf("─────────────────────────────────────────────────────────────\n")

	return nil
}

// handleQuery handles general AI queries
func handleQuery(ctx context.Context, ai *GeminiClient, query string) error {
	if query == "" {
		return fmt.Errorf("no query provided")
	}

	fmt.Printf("🤔 %s\n", colorize("Processing query...", "blue"))
	
	response, err := ai.GenerateResponse(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to generate response: %w", err)
	}

	response = strings.TrimSpace(response)
	
	fmt.Printf("\n%s\n", colorize("Response:", "green"))
	fmt.Printf("─────────────────────────────────────────────────────────────\n")
	fmt.Printf("%s\n", response)
	fmt.Printf("─────────────────────────────────────────────────────────────\n")

	return nil
}

// handleLastCommand analyzes the user's last executed command
func handleLastCommand(ctx context.Context, ai *GeminiClient) error {
	fmt.Printf("🔍 %s\n", colorize("Finding your last command...", "blue"))
	
	lastCmd, err := getLastCommand()
	if err != nil {
		return fmt.Errorf("failed to get last command: %w", err)
	}

	if lastCmd == "" {
		return fmt.Errorf("no recent command found in shell history")
	}

	fmt.Printf("📝 %s %s\n", colorize("Last command:", "cyan"), colorize(lastCmd, "white"))
	
	// Try to get the output of the last command if possible
	// Note: This is challenging since we don't have access to the actual output
	// For now, we'll just analyze the command itself
	fmt.Printf("🤔 %s\n", colorize("Analyzing command...", "blue"))
	
	// Analyze just the command for now
	response, err := ai.AnalyzeCommand(ctx, lastCmd, "")
	if err != nil {
		return fmt.Errorf("failed to analyze command: %w", err)
	}

	response = strings.TrimSpace(response)
	
	fmt.Printf("\n%s\n", colorize("Analysis:", "green"))
	fmt.Printf("─────────────────────────────────────────────────────────────\n")
	fmt.Printf("%s\n", response)
	fmt.Printf("─────────────────────────────────────────────────────────────\n")

	return nil
} 