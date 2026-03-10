package utils

import "testing"

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"CamelCase", "camel_case"},
		{"camelCase", "camel_case"},
		{"HTTPServer", "httpserver"},
		{"myHTTPServer", "my_httpserver"},
		{"already_snake", "already_snake"},
		{"with-hyphen", "with_hyphen"},
		{"Simple", "simple"},
		{"", ""},
		{"A", "a"},
		{"AB", "ab"},
		{"aB", "a_b"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := SnakeCase(tc.input)
			if result != tc.expected {
				t.Fatalf("SnakeCase(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestCapitalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"snake_case", "SnakeCase"},
		{"with-hyphen", "WithHyphen"},
		{"already", "Already"},
		{"UPPER", "UPPER"},
		{"a_b_c", "ABC"},
		{"", ""},
		{"a", "A"},
		{"hello_world", "HelloWorld"},
		{"multi_word_string", "MultiWordString"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := CapitalCase(tc.input)
			if result != tc.expected {
				t.Fatalf("CapitalCase(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}
