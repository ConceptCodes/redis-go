package resp_test

import (
	"bytes"
	"testing"

	"github.com/conceptcodes/redis-go/internal/resp"
)

func TestInvalidFormatError(t *testing.T) {
	// Arrange: Set up the parser with an invalid RESP format
	input := "invalid format"
	reader := bytes.NewReader([]byte(input))
	parser := resp.NewParser(reader)

	// Act: Attempt to parse the invalid input
	_, err := parser.Parse()

	// Assert: Check if the error is of type ErrInvalidFormat
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
	if _, ok := err.(*resp.ErrInvalidFormat); !ok {
		t.Fatalf("Expected ErrInvalidFormat, got %T", err)
	}
}

func TestGenericError(t *testing.T) {
	// Arrange: Set up the parser with an invalid RESP format
	input := "-ERR invalid command\r\n"
	reader := bytes.NewReader([]byte(input))
	parser := resp.NewParser(reader)

	// Act: Attempt to parse the invalid input
	result, err := parser.Parse()

	// Assert: Check if the result is nil and error is of type GenericError
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if _, ok := result.(*resp.GenericError); !ok {
		t.Fatalf("Expected GenericError, got %T", result)
	}
}

func TestSimpleStringParser(t *testing.T) {
	// Arrange: Set up the parser and the expected output
	tests := []struct {
		input    string
		expected string
	}{
		{"+OK\r\n", "OK"},
		{"+Hello World\r\n", "Hello World"},
	}

	// Act: Run the parser on the test inputs
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			reader := bytes.NewReader([]byte(test.input))
			parser := resp.NewParser(reader)
			result, err := parser.Parse()
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// Assert: Check if the result matches the expected output
			out, ok := result.(*resp.SimpleString)
			if !ok {
				t.Fatalf("Expected result to be a string, got %T", result)
			}
			if out.Value != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, out)
			}
		})
	}
}

func TestIntegerParser(t *testing.T) {
	// Arrange: Set up the parser and the expected output
	tests := []struct {
		input    string
		expected int64
	}{
		{":1000\r\n", 1000},
		{":-1000\r\n", -1000},
	}

	// Act: Run the parser on the test inputs
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			reader := bytes.NewReader([]byte(test.input))
			parser := resp.NewParser(reader)
			result, err := parser.Parse()
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// Assert: Check if the result matches the expected output
			out, ok := result.(*resp.Integer)
			if !ok {
				t.Fatalf("Expected result to be an int64, got %T", result)
			}
			if int64(out.Value) != test.expected {
				t.Errorf("Expected %d, got %d", test.expected, out)
			}
		})
	}
}

func TestBulkStringParser(t *testing.T) {
	// Arrange: Set up the parser and the expected output
	tests := []struct {
		input    string
		expected string
	}{
		{"$5\r\nHello\r\n", "Hello"},
		{"$0\r\n\r\n", ""},
	}

	// Act: Run the parser on the test inputs
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			reader := bytes.NewReader([]byte(test.input))
			parser := resp.NewParser(reader)
			result, err := parser.Parse()
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// Assert: Check if the result matches the expected output
			out, ok := result.(*resp.BulkString)
			if !ok {
				t.Fatalf("Expected result to be a string, got %T", result)
			}
			if out.Value != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, out.Value)
			}
		})
	}
}

func TestParseArray(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"*3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n", []string{"foo", "bar", "baz"}},
		{"*2\r\n$0\r\n\r\n$0\r\n\r\n", []string{"", ""}},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tc.input))
			parser := resp.NewParser(reader)
			result, err := parser.Parse()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			arr, ok := result.(*resp.Array)
			if !ok {
				t.Fatalf("Expected result to be *resp.Array, got %T", result)
			}
			if len(arr.Elements) != len(tc.expected) {
				t.Fatalf("Expected %d elements, got %d", len(tc.expected), len(arr.Elements))
			}
			for i, expectedVal := range tc.expected {
				bs, ok := arr.Elements[i].(*resp.BulkString)
				if !ok {
					t.Fatalf("Expected element %d to be *resp.BulkString, got %T", i, arr.Elements[i])
				}
				if bs.Value != expectedVal {
					t.Errorf("Expected element %d to be %q, got %q", i, expectedVal, bs.Value)
				}
			}
		})
	}
}
