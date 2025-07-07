package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:     "wdym",
	Short:   "What Do You Mean? - AI-powered CLI assistant",
	Long: `wdym (What Do You Mean?) is an AI-powered command-line assistant that helps you:
- Analyze your recently executed commands
- Generate shell commands from natural language
- Read and modify files with AI assistance
- Get contextual help and suggestions

Built with Go for speed and reliability.`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	Run:     runDefault,
}

var (
	shellMode   bool
	writeMode   bool
	codeMode    bool
	queryString string
	filePath    string
)

func init() {
	rootCmd.Flags().BoolVarP(&shellMode, "shell", "s", false, "Generate shell commands")
	rootCmd.Flags().BoolVarP(&writeMode, "write", "w", false, "Write or edit files")
	rootCmd.Flags().BoolVarP(&codeMode, "code", "c", false, "Generate code only")
	rootCmd.Flags().StringVarP(&queryString, "query", "q", "", "Query string for AI analysis")
	rootCmd.Flags().StringVarP(&filePath, "file", "f", "", "File to read or write (use @filename syntax)")
}

func runDefault(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	
	// Check if API key is set
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "Error: GEMINI_API_KEY environment variable is not set\n")
		fmt.Fprintf(os.Stderr, "Please set your Gemini API key: export GEMINI_API_KEY=\"your-api-key\"\n")
		os.Exit(1)
	}

	// Parse arguments and determine mode
	query := queryString
	if len(args) > 0 {
		query = args[0]
	}

	// Handle file syntax (@filename)
	if query != "" && query[0] == '@' {
		filePath = query[1:]
		if len(args) > 1 {
			query = args[1]
		} else {
			query = "What does this file do?"
		}
	}

	// Initialize AI client
	ai, err := NewGeminiClient(ctx, apiKey)
	if err != nil {
		log.Fatalf("Failed to initialize AI client: %v", err)
	}
	defer ai.Close()

	// Execute based on mode
	switch {
	case shellMode:
		err = handleShellMode(ctx, ai, query)
	case writeMode:
		err = handleWriteMode(ctx, ai, query, filePath)
	case codeMode:
		err = handleCodeMode(ctx, ai, query)
	case filePath != "":
		err = handleFileQuery(ctx, ai, query, filePath)
	case query != "":
		err = handleQuery(ctx, ai, query)
	default:
		err = handleLastCommand(ctx, ai)
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
} 