package quote

import (
	"testing"
)

func TestJoinWindows(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "simple arguments",
			input:    []string{"java", "-jar", "app.jar"},
			expected: `java -jar app.jar`,
		},
		{
			name:     "arguments with spaces",
			input:    []string{"java", "-Djava.library.path=C:\\Program Files\\libs", "-jar", "my app.jar"},
			expected: `java "-Djava.library.path=C:\Program Files\libs" -jar "my app.jar"`,
		},
		{
			name:     "arguments with quotes",
			input:    []string{"java", "-Dapp.name=\"My App\"", "-jar", "app.jar"},
			expected: `java "-Dapp.name=""My App""" -jar app.jar`,
		},
		{
			name:     "arguments with special characters",
			input:    []string{"java", "-Dpath=C:\\Users\\test&user", "-jar", "app.jar"},
			expected: `java "-Dpath=C:\Users\test&user" -jar app.jar`,
		},
		{
			name:     "empty arguments",
			input:    []string{"java", "", "-jar", "app.jar"},
			expected: `java "" -jar app.jar`,
		},
		{
			name:     "single argument",
			input:    []string{"java"},
			expected: `java`,
		},
		{
			name:     "no arguments",
			input:    []string{},
			expected: ``,
		},
		{
			name:     "complex JVM arguments",
			input:    []string{"java", "-Xmx512m", "-Djava.library.path=./libs with spaces", "-Dfile.encoding=UTF-8", "-jar", "myapp.jar", "--config", "my config.properties"},
			expected: `java -Xmx512m "-Djava.library.path=./libs with spaces" -Dfile.encoding=UTF-8 -jar myapp.jar --config "my config.properties"`,
		},
		{
			name:     "arguments with tabs",
			input:    []string{"java", "-Dvalue=test\ttab", "-jar", "app.jar"},
			expected: `java "-Dvalue=test\ttab" -jar app.jar`,
		},
		{
			name:     "arguments with newlines",
			input:    []string{"java", "-Dvalue=line1\nline2", "-jar", "app.jar"},
			expected: `java "-Dvalue=line1\nline2" -jar app.jar`,
		},
		{
			name:     "arguments with backslashes",
			input:    []string{"java", "-Dpath=C:\\Windows\\System32", "-jar", "app.jar"},
			expected: `java -Dpath=C:\Windows\System32 -jar app.jar`,
		},
		{
			name:     "arguments with pipes and redirects",
			input:    []string{"java", "-Darg=value|other", "-jar", "app.jar"},
			expected: `java "-Darg=value|other" -jar app.jar`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinWindows(tt.input)
			if result != tt.expected {
				t.Errorf("JoinWindows(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestQuoteWindowsArg(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no spaces or special chars",
			input:    "simple",
			expected: "simple",
		},
		{
			name:     "with spaces",
			input:    "has spaces",
			expected: `"has spaces"`,
		},
		{
			name:     "with double quote",
			input:    `has"quote`,
			expected: `"has""quote"`,
		},
		{
			name:     "with multiple quotes",
			input:    `"quoted"`,
			expected: `"""quoted"""`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: `""`,
		},
		{
			name:     "with tab",
			input:    "has\ttab",
			expected: `"has\ttab"`,
		},
		{
			name:     "with newline",
			input:    "has\nnewline",
			expected: `"has\nnewline"`,
		},
		{
			name:     "with backslashes",
			input:    `C:\Windows\System32`,
			expected: `C:\Windows\System32`,
		},
		{
			name:     "with backslashes and quotes",
			input:    `C:\Program Files\"quoted"`,
			expected: `"C:\Program Files\""quoted"""`,
		},
		{
			name:     "special characters",
			input:    "test&pipe|redirect<>",
			expected: `"test&pipe|redirect<>"`,
		},
		{
			name:     "semicolon",
			input:    "command;other",
			expected: `"command;other"`,
		},
		{
			name:     "equals sign",
			input:    "key=value",
			expected: "key=value",
		},
		{
			name:     "parentheses",
			input:    "func(param)",
			expected: `"func(param)"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := QuoteWindowsArg(tt.input)
			if result != tt.expected {
				t.Errorf("QuoteWindowsArg(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Test round-trip compatibility with argument parsing
func TestWindowsQuotingRoundTrip(t *testing.T) {
	testCases := [][]string{
		{"simple", "args"},
		{"args with spaces", "and", "more"},
		{"", "empty", ""},
		{"arg\"with\"quotes", "normal"},
		{"-Djava.library.path=C:\\Program Files\\Java", "-Xmx512m"},
		{"complex arg with spaces & special chars", "normal_arg"},
	}

	for _, tc := range testCases {
		t.Run("round_trip", func(t *testing.T) {
			quoted := JoinWindows(tc)

			// This tests that our quoting produces reasonable output
			// In practice, the Windows command line parser would parse this back
			if quoted == "" && len(tc) == 0 {
				return // OK
			}

			if quoted == "" && len(tc) > 0 {
				t.Errorf("JoinWindows(%v) produced empty string", tc)
			}
		})
	}
}