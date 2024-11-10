# ReqParser

ReqParser is a HTTP request parsing and formatting tool that accepts and logs HTTP requests, parses JSON data, and optionally converts it into Go or Rust struct definitions.

## Features

- Accepts and logs all HTTP methods (GET, POST, PUT, DELETE, etc.)
- Parses and displays JSON request bodies
- Optional conversion to programming language formats:
  - Go structs
  - Rust structs (with serde attributes)
- Pretty print JSON with delimiters
- Optional HTTP headers display
- Configurable port

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

## Usage

Start the server:

```bash
# Basic usage - shows compact JSON only
./reqparser

# Show pretty-printed JSON with delimiters
./reqparser -pretty

# Show HTTP headers
./reqparser -headers

# Generate Go struct and show JSON
./reqparser -format go

# Generate Rust struct and show pretty JSON
./reqparser -format rust -pretty

# Show headers and generate Go struct
./reqparser -headers -format go

# Specify a different port
./reqparser -port 3000
```

### Flag Behavior

- Without `-format`: Shows only JSON output
- With `-format`: Shows both struct definition and JSON output
- With `-pretty`: Shows JSON with delimiters
- Without `-pretty`: Shows compact JSON-Body format
- With `-headers`: Shows HTTP headers
- Without `-headers`: Headers are hidden

### Example Outputs

1. Basic usage (no flags):
```
JSON-Body: {"name":"test","value":123}
```

2. With `-pretty`:
```
==================
JSON START
{
    "name": "test",
    "value": 123
}
==================
JSON END
==================
```

3. With `-headers`:
```
Content-Type: application/json
User-Agent: curl/7.79.1
Accept: */*
Content-Length: 31

JSON-Body: {"name":"test","value":123}
```

4. With `-format go`:
```
JSON-Body: {"name":"test","value":123}
Struct format:
type GeneratedStruct struct {
    name string `json:"name"`
    value float64 `json:"value"`
}
```

5. With `-format rust`:
```
JSON-Body: {"name":"test","value":123}
Struct format:
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

reqparser is a HTTP request parsing and formatting tool

Options:
  -port int
        Port to run the server on (default 8080)
  -format string
        Output format type (go, rust) - if not provided, no struct will be generated
  -pretty
        Pretty print JSON with delimiters (if not provided, shows compact JSON-Body)
  -headers
        Show HTTP headers in output
  -version
        Show version information

Behavior:
  - Without -format: Shows only JSON (pretty or compact)
  - With -format: Shows struct and JSON (pretty or compact)
  - With -pretty: Shows JSON with delimiters
  - Without -pretty: Shows compact JSON-Body
  - With -headers: Shows HTTP headers
  - Without -headers: Headers are hidden
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
