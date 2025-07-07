# wdym

**wdym** (What Do You Mean?) is an AI-powered command-line assistant built with Go that helps you understand, generate, and work with terminal commands using natural language.

## Features

**AI-Powered Command Analysis**
- Analyze your recently executed commands
- Get contextual suggestions and explanations
- Understand what commands do and potential issues

**Shell Command Generation**
- Generate shell commands from natural language
- Safe command suggestions with execution confirmation
- OS and shell-aware recommendations

**File Operations**
- Create new files from natural language descriptions
- Edit existing files with AI assistance
- Read and analyze file contents with questions

**Code Generation**
- Generate code snippets from descriptions
- Language-agnostic code assistance
- Quick prototyping and examples

## Installation

### Prerequisites

1. **Go 1.19+** - [Install Go](https://golang.org/doc/install)
2. **Gemini API Key** - Get one from [Google AI Studio](https://aistudio.google.com/app/apikey)

### Install from Source

```bash
# Clone the repository
git clone https://github.com/jacobmortera/wdym.git
cd wdym

# Build the binary
go build -o wdym .

# Move to your PATH (optional)
sudo mv wdym /usr/local/bin/
```

### Set up API Key

```bash
export GEMINI_API_KEY="your-api-key-here"

# Add to your shell profile for persistence
echo 'export GEMINI_API_KEY="your-api-key-here"' >> ~/.zshrc
source ~/.zshrc
```

## Usage

### Basic Command Analysis

Analyze your last executed command:
```bash
$ ls -la
$ wdym
```

### Shell Command Generation

Generate shell commands from natural language:
```bash
$ wdym -s "find all python files larger than 1MB"
$ wdym --shell "kill all processes using port 8080"
```

### File Operations

Read and analyze files:
```bash
$ wdym @config.json "What does this configuration do?"
$ wdym @script.py "How can I optimize this code?"
```

Create or edit files:
```bash
$ wdym -w @hello.py "Create a Python script that prints hello world"
$ wdym --write @server.js "Add error handling to this Express server"
```

### Code Generation

Generate code snippets:
```bash
$ wdym -c "Create a function to validate email addresses in JavaScript"
$ wdym --code "Write a Python class for handling API requests"
```

### General Queries

Ask any question:
```bash
$ wdym "What's the difference between TCP and UDP?"
$ wdym "How do I set up a reverse proxy with nginx?"
```

## Command Reference

| Flag | Description | Example |
|------|-------------|---------|
| `-s, --shell` | Generate shell commands | `wdym -s "compress this directory"` |
| `-w, --write` | Write or edit files | `wdym -w @file.txt "add comments"` |
| `-c, --code` | Generate code only | `wdym -c "bubble sort in Python"` |
| `-q, --query` | Explicit query string | `wdym -q "explain this error"` |
| `-f, --file` | Specify file path | `wdym -f config.yaml "validate this"` |
| `-h, --help` | Show help | `wdym --help` |
| `-v, --version` | Show version | `wdym --version` |

## File Syntax

Use the `@filename` syntax to reference files:

```bash
# Read and analyze
wdym @package.json "What dependencies are used?"

# Edit with write mode
wdym -w @script.sh "add error checking"

# Works with paths
wdym @src/main.go "explain this function"
```

## Examples

### Shell Commands
```bash
$ wdym -s "show disk usage of current directory"
🐚 Generating shell command...

$ du -sh .
Execute this command? (y/N): y

Executing...
1.2G    .
```

### File Creation
```bash
$ wdym -w @fibonacci.py "create a function to generate fibonacci sequence"
✏️ Processing file: fibonacci.py
📝 Generating content...

📄 New file content:
─────────────────────────────────────────────────────────────
def fibonacci(n):
    """Generate fibonacci sequence up to n terms"""
    if n <= 0:
        return []
    elif n == 1:
        return [0]
    
    fib = [0, 1]
    for i in range(2, n):
        fib.append(fib[i-1] + fib[i-2])
    
    return fib
...
─────────────────────────────────────────────────────────────
Create file fibonacci.py? (y/N): y
✅ File saved: fibonacci.py
```

### Command Analysis
```bash
$ find . -name "*.go" -type f
$ wdym
🔍 Finding your last command...
📝 Last command: find . -name "*.go" -type f
🤔 Analyzing command...

Analysis:
─────────────────────────────────────────────────────────────
This command searches for all Go source files (.go extension) in the current 
directory and its subdirectories. The `-type f` flag ensures it only finds 
regular files, not directories. This is commonly used in Go projects to:

- Get an overview of all Go source files
- Count lines of code with `| xargs wc -l`
- Perform batch operations on Go files
- Check project structure
─────────────────────────────────────────────────────────────
```

## Configuration

**wdym** reads from environment variables:

- `GEMINI_API_KEY` - Your Gemini API key (required)
- `SHELL` - Current shell (auto-detected)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Powered by Google's [Gemini AI](https://ai.google.dev/)

## Roadmap

- [ ] Chat mode with conversation history
- [ ] Plugin system for custom commands
- [ ] Integration with popular development tools
- [ ] Shell integration with hotkeys
- [ ] Configuration file support
- [ ] Multiple AI provider support 
