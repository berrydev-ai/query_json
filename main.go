package main

// Test comment for pre-commit hook
import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ohler55/ojg/jp"
)

// Version information - set via ldflags during build
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	var query string
	var pretty bool
	var raw bool
	var showVersion bool
	flag.StringVar(&query, "query", "", "JSONPath query (e.g., $.root[0], $.users[*].name)")
	flag.BoolVar(&pretty, "pretty", true, "Pretty print JSON output")
	flag.BoolVar(&raw, "raw", false, "Output raw values (no JSON formatting for strings)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	if showVersion {
		fmt.Printf("query_json version %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built: %s\n", date)
		os.Exit(0)
	}

	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <json-file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s --query '$.users[0].name' ./examples/data.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --query '$.products[?(@.price > 100)]' ./examples/data.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --query '$.users[*].email' --raw ./examples/data.json\n", os.Args[0])
		os.Exit(1)
	}

	filename := flag.Args()[0]

	if query == "" {
		fmt.Fprintf(os.Stderr, "Error: --query parameter is required\n")
		os.Exit(1)
	}

	// Validate JSONPath syntax
	if err := validateJSONPath(query); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid JSONPath query: %v\n", err)
		os.Exit(1)
	}

	// Read JSON file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse JSON
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	// Parse JSONPath expression
	expr, err := jp.ParseString(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSONPath: %v\n", err)
		os.Exit(1)
	}

	// Execute JSONPath query
	result := expr.Get(jsonData)
	if len(result) == 0 {
		result = []interface{}{nil}
	}

	// If single result, unwrap it
	var finalResult interface{}
	if len(result) == 1 {
		finalResult = result[0]
	} else {
		finalResult = result
	}

	// Handle output formatting
	if err := outputResult(finalResult, pretty, raw); err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		os.Exit(1)
	}
}

func outputResult(result interface{}, pretty, raw bool) error {
	if result == nil {
		fmt.Println("null")
		return nil
	}

	// Handle raw output for simple types
	if raw {
		switch v := result.(type) {
		case string:
			fmt.Println(v)
			return nil
		case float64:
			fmt.Printf("%.10g\n", v)
			return nil
		case bool:
			fmt.Printf("%t\n", v)
			return nil
		case []interface{}:
			// For arrays of simple types, output each on a new line
			for _, item := range v {
				if str, ok := item.(string); ok {
					fmt.Println(str)
				} else {
					// Fall back to JSON for complex types
					output, err := json.Marshal(item)
					if err != nil {
						return err
					}
					fmt.Println(string(output))
				}
			}
			return nil
		}
	}

	// JSON output
	var output []byte
	var err error

	if pretty {
		output, err = json.MarshalIndent(result, "", "  ")
	} else {
		output, err = json.Marshal(result)
	}

	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}

// Additional helper function to validate JSONPath syntax
func validateJSONPath(path string) error {
	// Basic validation - the library will do the real validation
	if path == "" {
		return fmt.Errorf("empty JSONPath")
	}
	if !strings.HasPrefix(path, "$") {
		return fmt.Errorf("JSONPath must start with '$'")
	}
	return nil
}
