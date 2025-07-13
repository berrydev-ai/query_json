# JSON Query Tool

A fast, cross-platform command-line tool for querying JSON files using JSONPath expressions. Built with Go for excellent performance and zero runtime dependencies.

## Features

- 🚀 **Fast & Lightweight**: Single binary with no dependencies
- 🌐 **Cross-Platform**: Works on Windows, macOS, and Linux
- 📝 **Full JSONPath Support**: Powered by `github.com/PaesslerAG/jsonpath`
- 🎨 **Flexible Output**: Pretty-printed JSON or raw values
- 🔍 **Advanced Queries**: Filters, array slicing, recursive descent, and more
- ⚡ **Easy to Use**: Simple command-line interface

## Installation

### Download Pre-built Binaries

Download the latest binary for your platform from the [releases page](https://github.com/berrydev-ai/query_json/releases).

### Build from Source

Requirements:
- Go 1.21 or higher

```bash
git clone https://github.com/berrydev-ai/query_json.git
cd query_json
go build -o query_json
```

### Install via Go

```bash
go install github.com/berrydev-ai/query_json@latest
```

## Usage

```bash
query_json [options] <json-file>
```

### Options

- `--query`: JSONPath query expression (required)
- `--pretty`: Pretty print JSON output (default: true)
- `--raw`: Output raw values without JSON formatting (default: false)

### Basic Examples

```bash
# Get a specific field
query_json --query '$.name' data.json

# Get array element
query_json --query '$.users[0]' data.json

# Get nested field
query_json --query '$.users[0].email' data.json

# Get all elements from array
query_json --query '$.users[*].name' data.json

# Raw output (no JSON formatting)
query_json --query '$.users[*].email' --raw data.json
```

### Advanced JSONPath Queries

#### Filtering

```bash
# Users older than 25
query_json --query '$.users[?(@.age > 25)]' data.json

# Users with specific name
query_json --query '$.users[?(@.name == "Alice")]' data.json

# Products in electronics category
query_json --query '$.products[?(@.category == "electronics")]' data.json

# Multiple conditions with logical operators
query_json --query '$.products[?(@.price > 10 && @.category == "books")]' data.json

# Advanced filtering with various operators
query_json --query '$.users[?(@.age >= 30)]' data.json
query_json --query '$.users[?(@.active == true)]' data.json
```

#### Array Operations

```bash
# Array slicing (elements 1-3)
query_json --query '$.users[1:3]' data.json

# Last element
query_json --query '$.users[-1]' data.json

# Multiple fields
query_json --query '$.users[*]["name","email"]' data.json
```

#### Recursive Descent

```bash
# All name fields at any level
query_json --query '$..name' data.json

# All email addresses in the entire document
query_json --query '$..email' data.json
```

#### Regular Expressions

```bash
# Names starting with 'A' (regex support)
query_json --query '$.users[?(@.name =~ /^A.*/)]' data.json

# Email addresses containing 'gmail'
query_json --query '$.users[?(@.email =~ /gmail/)]' data.json

# Case-insensitive matching
query_json --query '$.users[?(@.name =~ /alice/i)]' data.json
```

## Example Data

Create a `test.json` file to experiment with:

```json
{
  "company": "Tech Corp",
  "employees": [
    {
      "name": "Alice Johnson",
      "age": 30,
      "email": "alice@techcorp.com",
      "department": "Engineering",
      "salary": 95000,
      "skills": ["Go", "Python", "Docker"]
    },
    {
      "name": "Bob Smith",
      "age": 25,
      "email": "bob@techcorp.com",
      "department": "Marketing",
      "salary": 65000,
      "skills": ["SEO", "Analytics", "Social Media"]
    },
    {
      "name": "Carol Davis",
      "age": 35,
      "email": "carol@techcorp.com",
      "department": "Engineering",
      "salary": 110000,
      "skills": ["JavaScript", "React", "AWS"]
    }
  ],
  "projects": [
    {
      "name": "Project Alpha",
      "status": "active",
      "budget": 50000,
      "team": ["Alice Johnson", "Carol Davis"]
    },
    {
      "name": "Project Beta",
      "status": "completed",
      "budget": 30000,
      "team": ["Bob Smith"]
    }
  ]
}
```

### Query Examples with Test Data

```bash
# Get all employee names
query_json --query '$.employees[*].name' test.json

# Get high-salary employees (>80k)
query_json --query '$.employees[?(@.salary > 80000)]' test.json

# Get engineering department employees
query_json --query '$.employees[?(@.department == "Engineering")]' test.json

# Get all skills across all employees
query_json --query '$.employees[*].skills[*]' --raw test.json

# Get active projects
query_json --query '$.projects[?(@.status == "active")]' test.json

# Get project names and budgets
query_json --query '$.projects[*]["name","budget"]' test.json
```

## Building for Multiple Platforms

Use the provided build script to create binaries for all major platforms:

```bash
#!/bin/bash

APP_NAME="query_json"

# Create builds directory
mkdir -p builds

echo "Building ${APP_NAME} for multiple platforms..."

# Windows
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o builds/${APP_NAME}-windows-amd64.exe
GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o builds/${APP_NAME}-windows-386.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o builds/${APP_NAME}-macos-amd64
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o builds/${APP_NAME}-macos-arm64

# Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o builds/${APP_NAME}-linux-amd64
GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o builds/${APP_NAME}-linux-386

echo "Build completed! Binaries are in ./builds/"
```

## JSONPath Syntax Reference

### Basic Syntax

| Expression | Description |
|------------|-------------|
| `$` | Root element |
| `.` | Child element |
| `..` | Recursive descent |
| `*` | Wildcard |
| `[n]` | Array index |
| `[start:end]` | Array slice |
| `[*]` | All array elements |
| `[?()]` | Filter expression |

### Filter Operators

✅ **Full JSONPath Support**: This library supports all standard JSONPath operators and functions.

| Operator | Description | Example | Supported |
|----------|-------------|---------|-----------|
| `==` | Equal | `$.users[?(@.age == 30)]` | ✅ |
| `!=` | Not equal | `$.users[?(@.status != "inactive")]` | ✅ |
| `>` | Greater than | `$.products[?(@.price > 100)]` | ✅ |
| `>=` | Greater than or equal | `$.users[?(@.age >= 18)]` | ✅ |
| `<` | Less than | `$.items[?(@.quantity < 10)]` | ✅ |
| `<=` | Less than or equal | `$.scores[?(@.value <= 80)]` | ✅ |
| `=~` | Regular expression | `$.users[?(@.name =~ /^A.*/)]` | ✅ |
| `&&` | Logical AND | `$.users[?(@.age > 25 && @.active == true)]` | ✅ |
| `\|\|` | Logical OR | `$.users[?(@.role == "admin" \|\| @.role == "moderator")]` | ✅ |

### Functions

⚠️ **Limited Function Support**: This JSONPath library supports basic operations but has limited function support.

| Function | Description | Example | Supported |
|----------|-------------|---------|-----------|
| `length()` | Array/string length | `$.users[?(@.skills.length() > 3)]` | ❌ |
| `size()` | Alias for length | `$.items[?(@.tags.size() > 0)]` | ❌ |
| `keys()` | Object keys | `$.users[?(@.keys().length() > 5)]` | ❌ |
| `values()` | Object values | `$.users[?(@.values().length() > 0)]` | ❌ |

**Note**: For array/string length comparisons, consider using array indexing or nested queries instead.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Dependencies

- [github.com/ohler55/ojg](https://github.com/ohler55/ojg) - Optimized JSON processing library with full JSONPath support

## Acknowledgments

- Thanks to the [ohler55](https://github.com/ohler55) team for the excellent OJG library with full JSONPath support
- Inspired by the need for a fast, cross-platform JSON querying tool
- Built with ❤️ in Go