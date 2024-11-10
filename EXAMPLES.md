# ReqParser Examples

This document provides example commands for testing different features of ReqParser.

## Basic JSON Output

Shows compact JSON format without any additional formatting.

```bash
# Start server
./reqparser

# Test with curl
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"name":"test","age":30,"active":true}' \
  http://localhost:8080/api/data

# Expected output:
JSON-Body: {"name":"test","age":30,"active":true}
```

## Pretty JSON Output

Shows JSON with delimiters and proper indentation.

```bash
# Start server
./reqparser -pretty

# Test with curl
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"user":{"name":"test","email":"test@example.com"},"settings":{"theme":"dark","notifications":true}}' \
  http://localhost:8080/api/data

# Expected output:
==========
JSON START
==========
{
    "user": {
        "name": "test",
        "email": "test@example.com"
    },
    "settings": {
        "theme": "dark",
        "notifications": true
    }
}
========
JSON END
========
```

## Show HTTP Headers

Displays all HTTP headers along with the JSON body.

```bash
# Start server
./reqparser -headers

# Test with curl (with custom headers)
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -H "X-Custom-Header: test-value" \
  -d '{"message":"test","priority":"high"}' \
  http://localhost:8080/api/data

# Expected output:
Content-Type: application/json
Authorization: Bearer test-token
X-Custom-Header: test-value
Content-Length: 39

JSON-Body: {"message":"test","priority":"high"}
```

## Go Struct Generation

Generates Go struct definitions from JSON input.

```bash
# Start server
./reqparser -format go

# Test with complex JSON
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"user":{"id":123,"name":"test","roles":["admin","user"]},"metadata":{"created_at":"2024-01-01","tags":["important","urgent"]}}' \
  http://localhost:8080/api/data

# Expected output:
JSON-Body: {"user":{"id":123,"name":"test","roles":["admin","user"]},"metadata":{"created_at":"2024-01-01","tags":["important","urgent"]}}
Struct format:
type GeneratedStruct struct {
    user struct {
        id float64 `json:"id"`
        name string `json:"name"`
        roles []interface{} `json:"roles"`
    } `json:"user"`
    metadata struct {
        created_at string `json:"created_at"`
        tags []interface{} `json:"tags"`
    } `json:"metadata"`
}
```

## Rust Struct with Pretty Print

Generates Rust struct definitions and shows pretty-printed JSON.

```bash
# Start server
./reqparser -format rust -pretty

# Test with nested structures
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"config":{"database":{"host":"localhost","port":5432,"credentials":{"username":"admin","password":"secret"}}},"features":["logging","metrics"]}' \
  http://localhost:8080/api/data

# Expected output:
==========
JSON START
==========
{
    "config": {
        "database": {
            "host": "localhost",
            "port": 5432,
            "credentials": {
                "username": "admin",
                "password": "secret"
            }
        }
    },
    "features": [
        "logging",
        "metrics"
    ]
}
========
JSON END
========

Struct format:
#[derive(Debug, Serialize, Deserialize)]
struct GeneratedStruct {
    #[serde(rename = "config")]
    config: Config,
    #[serde(rename = "features")]
    features: Vec<String>,
}

#[derive(Debug, Serialize, Deserialize)]
struct Config {
    #[serde(rename = "database")]
    database: Database,
}

#[derive(Debug, Serialize, Deserialize)]
struct Database {
    #[serde(rename = "host")]
    host: String,
    #[serde(rename = "port")]
    port: f64,
    #[serde(rename = "credentials")]
    credentials: Credentials,
}

#[derive(Debug, Serialize, Deserialize)]
struct Credentials {
    #[serde(rename = "username")]
    username: String,
    #[serde(rename = "password")]
    password: String,
}
```

## Combined Flags

Shows headers, pretty JSON, and struct definitions.

```bash
# Start server
./reqparser -format go -pretty -headers

# Test with complex request
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -H "User-Agent: reqparser-test" \
  -H "Accept: application/json" \
  -d '{"request":{"method":"POST","path":"/api/data"},"payload":{"items":[{"id":1,"name":"item1"},{"id":2,"name":"item2"}]},"metadata":{"timestamp":"2024-01-01T12:00:00Z"}}' \
  http://localhost:8080/api/data

# Expected output:
Content-Type: application/json
Authorization: Bearer test-token
User-Agent: reqparser-test
Accept: application/json
Content-Length: 158

==========
JSON START
==========
{
    "request": {
        "method": "POST",
        "path": "/api/data"
    },
    "payload": {
        "items": [
            {
                "id": 1,
                "name": "item1"
            },
            {
                "id": 2,
                "name": "item2"
            }
        ]
    },
    "metadata": {
        "timestamp": "2024-01-01T12:00:00Z"
    }
}
========
JSON END
========

Struct format:
type GeneratedStruct struct {
    request struct {
        method string `json:"method"`
        path string `json:"path"`
    } `json:"request"`
    payload struct {
        items []struct {
            id float64 `json:"id"`
            name string `json:"name"`
        } `json:"items"`
    } `json:"payload"`
    metadata struct {
        timestamp string `json:"timestamp"`
    } `json:"metadata"`
}
