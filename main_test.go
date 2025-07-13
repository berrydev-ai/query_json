package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ohler55/ojg/jp"
)

// Test data for various test cases
const testJSON = `{
  "users": [
    {
      "name": "Alice",
      "age": 30,
      "email": "alice@example.com",
      "city": "New York",
      "active": true
    },
    {
      "name": "Bob",
      "age": 25,
      "email": "bob@example.com",
      "city": "Los Angeles",
      "active": false
    },
    {
      "name": "Charlie",
      "age": 35,
      "email": "charlie@example.com",
      "city": "Chicago",
      "active": true
    }
  ],
  "products": [
    {
      "name": "Laptop",
      "price": 999.99,
      "category": "electronics",
      "inStock": true
    },
    {
      "name": "Book",
      "price": 15.99,
      "category": "books",
      "inStock": false
    }
  ],
  "company": {
    "name": "Tech Corp",
    "founded": 2010,
    "employees": 150
  }
}`

const emptyJSON = `{}`
const arrayJSON = `[1, 2, 3, 4, 5]`
const primitiveJSON = `"hello world"`

// Helper function to create a temporary JSON file
func createTempJSONFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.json")

	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	return tmpFile
}

// Test outputResult function with various data types
func TestOutputResult(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		pretty   bool
		raw      bool
		expected string
	}{
		{
			name:     "null value",
			input:    nil,
			pretty:   true,
			raw:      false,
			expected: "null\n",
		},
		{
			name:     "string value with raw output",
			input:    "hello",
			pretty:   false,
			raw:      true,
			expected: "hello\n",
		},
		{
			name:     "string value with JSON output",
			input:    "hello",
			pretty:   false,
			raw:      false,
			expected: "\"hello\"\n",
		},
		{
			name:     "number value with raw output",
			input:    42.5,
			pretty:   false,
			raw:      true,
			expected: "42.5\n",
		},
		{
			name:     "boolean value with raw output",
			input:    true,
			pretty:   false,
			raw:      true,
			expected: "true\n",
		},
		{
			name:     "array of strings with raw output",
			input:    []interface{}{"alice@example.com", "bob@example.com"},
			pretty:   false,
			raw:      true,
			expected: "alice@example.com\nbob@example.com\n",
		},
		{
			name:     "object with pretty JSON",
			input:    map[string]interface{}{"name": "Alice", "age": 30},
			pretty:   true,
			raw:      false,
			expected: "{\n  \"age\": 30,\n  \"name\": \"Alice\"\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output by redirecting stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := outputResult(tt.input, tt.pretty, tt.raw)

			w.Close()
			os.Stdout = oldStdout

			if err != nil {
				t.Fatalf("outputResult returned error: %v", err)
			}

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			output := string(buf[:n])

			if output != tt.expected {
				t.Errorf("Expected output %q, got %q", tt.expected, output)
			}
		})
	}
}

// Test validateJSONPath function
func TestValidateJSONPath(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "valid root path",
			path:        "$",
			expectError: false,
		},
		{
			name:        "valid field access",
			path:        "$.users",
			expectError: false,
		},
		{
			name:        "valid array access",
			path:        "$.users[0]",
			expectError: false,
		},
		{
			name:        "valid complex query",
			path:        "$.users[?(@.age > 25)]",
			expectError: false,
		},
		{
			name:        "empty path",
			path:        "",
			expectError: true,
		},
		{
			name:        "path without dollar sign",
			path:        "users[0]",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateJSONPath(tt.path)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for path %q, but got none", tt.path)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for path %q, but got: %v", tt.path, err)
			}
		})
	}
}

// Integration tests for various JSONPath queries
func TestJSONPathQueries(t *testing.T) {
	tempFile := createTempJSONFile(t, testJSON)
	defer os.Remove(tempFile)

	// Parse the test JSON for expected results
	var testData interface{}
	err := json.Unmarshal([]byte(testJSON), &testData)
	if err != nil {
		t.Fatalf("Failed to parse test JSON: %v", err)
	}

	tests := []struct {
		name          string
		query         string
		expectedCount int // For arrays, count of elements
		expectedType  string
	}{
		{
			name:          "root access",
			query:         "$",
			expectedCount: 1,
			expectedType:  "object",
		},
		{
			name:          "get all users",
			query:         "$.users",
			expectedCount: 3,
			expectedType:  "array",
		},
		{
			name:          "get first user",
			query:         "$.users[0]",
			expectedCount: 1,
			expectedType:  "object",
		},
		{
			name:          "get all user names",
			query:         "$.users[*].name",
			expectedCount: 3,
			expectedType:  "array",
		},
		{
			name:          "filter users by age (greater than)",
			query:         "$.users[?(@.age > 25)]", // Now supports all comparison operators
			expectedCount: 2,
			expectedType:  "array",
		},
		{
			name:          "get company name",
			query:         "$.company.name",
			expectedCount: 1,
			expectedType:  "string",
		},
		{
			name:          "recursive descent for all names",
			query:         "$..name",
			expectedCount: 6, // 3 users + 2 products + 1 company
			expectedType:  "array",
		},
		{
			name:          "filter users greater than or equal to age",
			query:         "$.users[?(@.age >= 30)]",
			expectedCount: 2,
			expectedType:  "array",
		},
		{
			name:          "filter users less than age",
			query:         "$.users[?(@.age < 30)]",
			expectedCount: 1,
			expectedType:  "object", // Single result returns object, not array
		},
		{
			name:          "filter active users",
			query:         "$.users[?(@.active == true)]",
			expectedCount: 2,
			expectedType:  "array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read the temp file
			file, err := os.Open(tempFile)
			if err != nil {
				t.Fatalf("Failed to open temp file: %v", err)
			}
			defer file.Close()

			data, err := os.ReadFile(tempFile)
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}

			var jsonData interface{}
			if err := json.Unmarshal(data, &jsonData); err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			// Execute JSONPath query (this tests the actual jsonpath library integration)
			result, err := executeJSONPathQuery(tt.query, jsonData)
			if err != nil {
				t.Fatalf("Query %q failed: %v", tt.query, err)
			}

			// Validate result type and count
			switch tt.expectedType {
			case "array":
				if arr, ok := result.([]interface{}); ok {
					if len(arr) != tt.expectedCount {
						t.Errorf("Query %q expected %d elements, got %d", tt.query, tt.expectedCount, len(arr))
					}
				} else {
					t.Errorf("Query %q expected array result, got %T", tt.query, result)
				}
			case "object":
				if _, ok := result.(map[string]interface{}); !ok {
					t.Errorf("Query %q expected object result, got %T", tt.query, result)
				}
			case "string":
				if _, ok := result.(string); !ok {
					t.Errorf("Query %q expected string result, got %T", tt.query, result)
				}
			}
		})
	}
}

// Test error cases
func TestErrorCases(t *testing.T) {
	tests := []struct {
		name          string
		jsonContent   string
		query         string
		expectError   bool
		errorContains string
	}{
		{
			name:          "invalid JSON",
			jsonContent:   `{"invalid": json}`,
			query:         "$.invalid",
			expectError:   true,
			errorContains: "JSON",
		},
		{
			name:          "invalid JSONPath query",
			jsonContent:   testJSON,
			query:         "$.users[invalid]",
			expectError:   true,
			errorContains: "invalid",
		},
		{
			name:        "non-existent field",
			jsonContent: testJSON,
			query:       "$.nonexistent",
			expectError: false, // This library returns null for non-existent fields
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile := createTempJSONFile(t, tt.jsonContent)
			defer os.Remove(tempFile)

			data, err := os.ReadFile(tempFile)
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}

			var jsonData interface{}
			jsonErr := json.Unmarshal(data, &jsonData)

			if tt.errorContains == "JSON" && jsonErr == nil {
				t.Errorf("Expected JSON parsing error but got none")
				return
			}

			if jsonErr != nil && tt.errorContains != "JSON" {
				t.Fatalf("Unexpected JSON parsing error: %v", jsonErr)
			}

			if jsonErr == nil {
				_, queryErr := executeJSONPathQuery(tt.query, jsonData)

				if tt.expectError && queryErr == nil {
					t.Errorf("Expected query error but got none")
				}

				if !tt.expectError && queryErr != nil {
					t.Errorf("Unexpected query error: %v", queryErr)
				}

				if queryErr != nil && tt.errorContains != "" && !strings.Contains(queryErr.Error(), tt.errorContains) {
					t.Errorf("Error message should contain %q, got: %v", tt.errorContains, queryErr)
				}
			}
		})
	}
}

// Test with different JSON structures
func TestDifferentJSONStructures(t *testing.T) {
	tests := []struct {
		name        string
		jsonContent string
		query       string
		expectError bool
	}{
		{
			name:        "empty object",
			jsonContent: emptyJSON,
			query:       "$",
			expectError: false,
		},
		{
			name:        "array root",
			jsonContent: arrayJSON,
			query:       "$[0]",
			expectError: false,
		},
		{
			name:        "primitive root",
			jsonContent: primitiveJSON,
			query:       "$",
			expectError: false,
		},
		{
			name:        "nested arrays",
			jsonContent: `{"matrix": [[1,2], [3,4]]}`,
			query:       "$.matrix[0][1]",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempFile := createTempJSONFile(t, tt.jsonContent)
			defer os.Remove(tempFile)

			data, err := os.ReadFile(tempFile)
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}

			var jsonData interface{}
			if err := json.Unmarshal(data, &jsonData); err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			_, err = executeJSONPathQuery(tt.query, jsonData)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for query %q but got none", tt.query)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for query %q: %v", tt.query, err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkJSONPathQuery(b *testing.B) {
	var jsonData interface{}
	err := json.Unmarshal([]byte(testJSON), &jsonData)
	if err != nil {
		b.Fatalf("Failed to parse test JSON: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := executeJSONPathQuery("$.users[*].name", jsonData)
		if err != nil {
			b.Fatalf("Query failed: %v", err)
		}
	}
}

func BenchmarkOutputResult(b *testing.B) {
	data := map[string]interface{}{
		"name": "Alice",
		"age":  30,
		"city": "New York",
	}

	// Redirect stdout to discard output during benchmark
	oldStdout := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = oldStdout }()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := outputResult(data, true, false)
		if err != nil {
			b.Fatalf("outputResult failed: %v", err)
		}
	}
}

// Helper function to execute JSONPath queries (wrapper around the library)
func executeJSONPathQuery(query string, data interface{}) (interface{}, error) {
	// Parse JSONPath expression
	expr, err := jp.ParseString(query)
	if err != nil {
		return nil, err
	}

	// Execute JSONPath query
	result := expr.Get(data)
	if len(result) == 0 {
		return nil, nil
	}

	// If single result, unwrap it
	if len(result) == 1 {
		return result[0], nil
	}

	return result, nil
}
