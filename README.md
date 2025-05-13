# Termpilot

Termpilot is a terminal-based AI chat agent powered by Ollama models.

## Features

- Chat with AI models via terminal interface
- Save and manage conversations
- TUI mode for easy browsing and chatting
- Command-line interface for quick interactions

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/termpilot.git
cd termpilot

# Install dependencies
go mod tidy

# Build the application
go build -o termpilot
```

## Usage

```bash
# Start a new chat
./termpilot chat "Your message here"

# List saved conversations
./termpilot chat --list

# Continue a conversation
./termpilot chat --continue <conversation-id> "Your follow-up message"

# Show a conversation
./termpilot chat --show <conversation-id>

# Launch the TUI
./termpilot
```

## Testing

The project includes a comprehensive test suite covering:

- Database operations
- Ollama client API
- CLI commands
- Data models

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make coverage

# Run specific package tests
make test-db
make test-models
make test-cmd
make test-ollama
```

### Test Structure

- `db/database_test.go` - Tests for database operations
- `models/models_test.go` - Tests for data models and GORM functionality
- `ollamaclient/ollamaclient_test.go` - Tests for Ollama API client
- `cmd/commands_test.go` - Tests for CLI commands
- `testutils/testutils.go` - Common testing utilities

### Adding Tests

When adding new features, please also add corresponding tests. Follow these guidelines:

1. Unit tests should be added for each package
2. Integration tests should verify component interactions
3. Use mocks for external dependencies (e.g., Ollama API)
4. Use the testutils package for common test functionality

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [GORM](https://gorm.io/) - ORM library
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [testify](https://github.com/stretchr/testify) - Testing assertions

## License

MIT