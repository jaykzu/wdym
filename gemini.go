package main

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiClient wraps the Gemini AI client
type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient(ctx context.Context, apiKey string) (*GeminiClient, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Use Gemini 1.5 Flash model
	model := client.GenerativeModel("gemini-1.5-flash")
	
	// Configure the model for better CLI assistance
	model.SetTemperature(0.1) // Lower temperature for more focused responses
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(2048)

	// Set safety settings to be more permissive for code generation
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}

	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

// Close closes the Gemini client
func (g *GeminiClient) Close() error {
	return g.client.Close()
}

// GenerateResponse generates a response from the AI
func (g *GeminiClient) GenerateResponse(ctx context.Context, prompt string) (string, error) {
	response, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(response.Candidates) == 0 {
		return "", fmt.Errorf("no response candidates generated")
	}

	candidate := response.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from AI")
	}

	// Extract text from the response
	var result string
	for _, part := range candidate.Content.Parts {
		if textPart, ok := part.(genai.Text); ok {
			result += string(textPart)
		}
	}

	return result, nil
}

// GenerateShellCommand generates a shell command from natural language
func (g *GeminiClient) GenerateShellCommand(ctx context.Context, query string) (string, error) {
	prompt := fmt.Sprintf(`You are a command-line expert. Generate a shell command for the following request.

Operating System: macOS (Darwin)
Shell: zsh

Request: %s

Provide ONLY the shell command, no explanation or markdown formatting. The command should be safe to execute.`, query)

	return g.GenerateResponse(ctx, prompt)
}

// GenerateCode generates code from natural language
func (g *GeminiClient) GenerateCode(ctx context.Context, query string) (string, error) {
	prompt := fmt.Sprintf(`You are a coding expert. Generate code for the following request.

Request: %s

Provide ONLY the code, no explanation or markdown formatting unless specifically requested.`, query)

	return g.GenerateResponse(ctx, prompt)
}

// AnalyzeFile analyzes a file and answers questions about it
func (g *GeminiClient) AnalyzeFile(ctx context.Context, content, query string) (string, error) {
	prompt := fmt.Sprintf(`You are a code and file analysis expert. Analyze the following file content and answer the question.

File Content:
%s

Question: %s

Provide a clear and concise answer.`, content, query)

	return g.GenerateResponse(ctx, prompt)
}

// AnalyzeCommand analyzes a command and its output
func (g *GeminiClient) AnalyzeCommand(ctx context.Context, command, output string) (string, error) {
	prompt := fmt.Sprintf(`You are a command-line expert. Analyze the following command and its output, then provide helpful insights, explanations, or suggestions.

Command: %s
Output: %s

Provide helpful insights about what this command does, any potential issues, or suggestions for next steps.`, command, output)

	return g.GenerateResponse(ctx, prompt)
}

// EditFile generates instructions or code to edit a file
func (g *GeminiClient) EditFile(ctx context.Context, filePath, content, instruction string) (string, error) {
	prompt := fmt.Sprintf(`You are a file editing expert. You need to modify the following file according to the given instruction.

File Path: %s
Current Content:
%s

Instruction: %s

Provide the complete modified file content. If creating a new file, provide the full file content.`, filePath, content, instruction)

	return g.GenerateResponse(ctx, prompt)
} 