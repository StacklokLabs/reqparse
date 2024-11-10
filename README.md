# ReqParser

ReqParser is a HTTP request parsing and formatting tool that accepts and logs HTTP requests, parses JSON data, and converts it into Go or Rust struct definitions.

## Features

- Accepts and logs all HTTP methods (GET, POST, PUT, DELETE, etc.)
- Parses JSON request bodies into structured data
- Converts JSON data into different programming language formats:
  - Go structs
  - Rust structs (with serde attributes)
- Verbose logging option for detailed request information
- Raw request logging showing complete headers and body
- Configurable port and output format

## Installation

```bash
go install github.com/lukehinds/reqparser@latest
```

Or clone and build manually:

```bash
git clone https://github.com/lukehinds/reqparser.git
cd reqparser
make build
```

## Development

### Prerequisites

- Go 1.20 or later
- Make
- golangci-lint (installed automatically via Makefile)
- gosec (installed automatically via Makefile)

### Available Make Commands

```bash
# Show all available commands
make help

# Build the binary
make build

# Run tests
make test

# Generate test coverage report
make coverage

# Run linter
make lint

# Format code
make fmt

# Run security check
make sec

# Clean up build artifacts
make clean

# Run all checks (format, lint, security, tests)
make check

# Build and run the binary
make run
```

## Usage

Start the server:

```bash
# Basic usage (defaults to port 8080 and Go format)
./reqparser

# Specify a different port
./reqparser -port 3000

# Enable verbose logging
./reqparser -verbose

# Change output format (go or rust)
./reqparser -format rust

# Show version
./reqparser -version

# Show help
./reqparser -help
```

### Example Requests

Send a POST request with JSON data:

```bash
# Using curl
curl -X POST -H "Content-Type: application/json" \
     -d '{"name": "test", "value": 123}' \
     http://localhost:8080/api/data
```

Example output with Go format:

```go
type GeneratedStruct struct {
    name string `json:"name"`
    value float64 `json:"value"`
}
```

Example output with Rust format:

```rust
#[derive(Debug, Serialize, Deserialize)]
struct GeneratedStruct {
    #[serde(rename = "name")]
    name: String,
    #[serde(rename = "value")]
    value: f64,
}
```

## Command Line Options

```
Usage of reqparser:

reqparser is a HTTP request parsing and formatting tool that converts JSON data into Go or Rust structs

Options:
  -port int
        Port to run the server on (default 8080)
  -verbose
        Enable verbose logging
  -format string
        Output format type (go, rust) (default "go")
  -version
        Show version information
```

## Project Structure

- `main.go`: Entry point and CLI interface
- `server/server.go`: Core server implementation
- `server/server_test.go`: Server tests
- `Makefile`: Build and development commands
- `.github/workflows/ci.yml`: CI/CD pipeline configuration

## Continuous Integration

The project uses GitHub Actions for CI/CD, which automatically:

- Runs tests on multiple Go versions
- Performs code linting
- Checks code formatting
- Runs security scanning
- Generates test coverage reports
- Creates releases (when tags are pushed)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests and checks (`make check`)
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
