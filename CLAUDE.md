# Claude Code Project Instructions

## Project Overview

This is a cross-platform JSON query tool built in Go that allows users to query JSON files using JSONPath expressions. The tool compiles to a single binary with no runtime dependencies and works on Windows, macOS, and Linux.

## Project Structure

```
query_json/
├── main.go       # Main application code
├── go.mod        # Go module definition
├── go.sum        # Go module checksums
├── README.md     # Project documentation
├── .gitignore    # Git ignore rules
├── LICENSE       # MIT license
├── build.sh      # Cross-platform build script
├── builds/       # Build output directory (ignored by git)
├── examples/     # Example JSON files for testing
└── CLAUDE.md     # This file
```

## Key Dependencies

- **Primary**: `github.com/PaesslerAG/jsonpath` - JSONPath implementation
- **Indirect**: `github.com/PaesslerAG/gval` - Expression evaluation library
- **Go Version**: 1.21 or higher

## Build Instructions

### Development Build
```bash
go build -o query_json
```

### Cross-Platform Build
Use the provided `build.sh` script:
```bash
chmod +x build.sh
./build.sh
```

This creates binaries for:
- Windows (amd64, 386)
- macOS (amd64, arm64)
- Linux (amd64, 386)

### Build Flags
For release builds, use these flags to reduce binary size:
```bash
go build -ldflags="-s -w" -o query_json
```

## Testing

### Manual Testing
Create test JSON files in the `examples/` directory and test with:
```bash
./query_json examples/test.json --query '$.users[0].name'
```

### Common Test Cases
- Basic field access: `$.name`
- Array indexing: `$.users[0]`
- Filtering: `$.users[?(@.age > 25)]`
- Wildcard selection: `$.users[*].name`
- Raw output: `$.users[*].email --raw`

## Development Guidelines

### Code Style
- Follow standard Go conventions
- Use `gofmt` for formatting
- Keep functions focused and testable
- Add error handling for all operations

### Error Handling
- Always handle file I/O errors
- Provide helpful error messages to users
- Use `fmt.Fprintf(os.Stderr, ...)` for errors
- Exit with appropriate codes (0 for success, 1 for errors)

### Output Formats
- Default: Pretty-printed JSON
- `--raw` flag: Raw string values (useful for shell scripting)
- `--pretty=false`: Compact JSON output

## Key Features to Maintain

1. **Zero Dependencies**: Binary should be self-contained
2. **Cross-Platform**: Must work on Windows, macOS, Linux
3. **JSONPath Compliance**: Full JSONPath syntax support
4. **Performance**: Fast execution for large JSON files
5. **User-Friendly**: Clear help text and error messages

## Common Tasks

### Adding New Command-Line Options
- Add flag definition in `main()`
- Update help text and examples
- Handle the new option in the appropriate function
- Update README.md with new usage examples

### Improving Performance
- Profile with `go tool pprof` for large files
- Consider streaming JSON parsing for very large files
- Optimize memory usage for large result sets

### Adding Features
- JSONPath syntax validation
- Multiple output formats (CSV, TSV, etc.)
- Batch processing of multiple files
- Configuration file support

## Release Process

1. **Update Version**: Update version in build script
2. **Test**: Run comprehensive tests on all platforms
3. **Build**: Generate all platform binaries with `build.sh`
4. **Document**: Update README.md and CHANGELOG.md
5. **Tag**: Create git tag with version number
6. **Release**: Create GitHub release with binaries

## Troubleshooting

### Common Issues
- **Module not found**: Run `go mod download`
- **Build fails**: Check Go version (requires 1.21+)
- **JSONPath errors**: Validate syntax with online JSONPath tester
- **Large files**: Consider memory usage and streaming options

### Debug Mode
Add verbose logging for debugging:
```go
import "log"

// Add debug flag
var debug bool
flag.BoolVar(&debug, "debug", false, "Enable debug output")

// Use in code
if debug {
    log.Printf("Processing query: %s", query)
}
```

## Example Usage Patterns

### Shell Integration
```bash
# Get all email addresses
./query_json users.json --query '$.users[*].email' --raw

# Filter and pipe to other tools
./query_json data.json --query '$.products[?(@.price > 100)]' | jq '.[] | .name'

# Use in scripts
ADMIN_EMAIL=$(./query_json config.json --query '$.admin.email' --raw)
```

### Batch Processing
```bash
# Process multiple files
for file in *.json; do
    echo "Processing $file"
    ./query_json "$file" --query '$.summary.total'
done
```

## Security Considerations

- **File Access**: Only read files explicitly provided by user
- **Input Validation**: Sanitize JSONPath expressions
- **Memory Usage**: Handle large files gracefully
- **Error Information**: Don't expose sensitive file paths in errors

## Performance Benchmarks

Target performance metrics:
- **Small files** (<1MB): Sub-second response
- **Medium files** (1-10MB): Under 5 seconds
- **Large files** (>10MB): Graceful handling with progress indication

## Future Enhancements

Potential areas for improvement:
- Interactive mode with REPL
- JSON schema validation
- Plugin system for custom output formats
- Integration with popular JSON tools
- Web interface for complex queries

## Notes for Claude Code

- Always test changes with the provided example files
- Maintain backward compatibility with existing command-line interface
- Focus on performance and reliability over feature bloat
- Keep the binary size reasonable (target <10MB)
- Ensure all error messages are helpful and actionable