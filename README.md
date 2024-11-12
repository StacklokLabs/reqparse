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

## Installation

```bash
go install github.com/stackloklabs/reqparser@latest
```

Or clone and build manually:

```bash
git clone https://github.com/stackloklabs/reqparser.git
cd reqparser
make build
```

### Flag Behavior

- With `-pretty`: Shows JSON with delimiters
- Without `-pretty`: Shows compact JSON-Body format
- With `-headers`: Shows HTTP headers
- With `-format go|rust` Generates a struct

* Note: Only parent structs, need to code up child struct generation 

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
```

### Prerequisites

- Go 1.20 or later
- Make
- golangci-lint (installed automatically via Makefile)
- gosec (installed automatically via Makefile)
