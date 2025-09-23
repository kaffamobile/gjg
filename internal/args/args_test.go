package args

import (
	"reflect"
	"testing"
)

func TestExtractSpecial(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		wantDebug   bool
		wantDryRun  bool
		wantForward []string
	}{
		{
			name:        "no special flags",
			input:       []string{"--verbose", "arg1", "arg2"},
			wantDebug:   false,
			wantDryRun:  false,
			wantForward: []string{"--verbose", "arg1", "arg2"},
		},
		{
			name:        "debug flag only",
			input:       []string{"--gjg-debug", "--verbose", "arg1"},
			wantDebug:   true,
			wantDryRun:  false,
			wantForward: []string{"--verbose", "arg1"},
		},
		{
			name:        "dry-run flag only",
			input:       []string{"--gjg-dry-run", "--verbose", "arg1"},
			wantDebug:   false,
			wantDryRun:  true,
			wantForward: []string{"--verbose", "arg1"},
		},
		{
			name:        "both special flags",
			input:       []string{"--gjg-debug", "--gjg-dry-run", "--verbose"},
			wantDebug:   true,
			wantDryRun:  true,
			wantForward: []string{"--verbose"},
		},
		{
			name:        "special flags mixed with regular args",
			input:       []string{"arg1", "--gjg-debug", "arg2", "--gjg-dry-run", "arg3"},
			wantDebug:   true,
			wantDryRun:  true,
			wantForward: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:        "empty input",
			input:       []string{},
			wantDebug:   false,
			wantDryRun:  false,
			wantForward: []string{},
		},
		{
			name:        "special flags at end",
			input:       []string{"arg1", "arg2", "--gjg-debug", "--gjg-dry-run"},
			wantDebug:   true,
			wantDryRun:  true,
			wantForward: []string{"arg1", "arg2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDebug, gotDryRun, gotForward := ExtractSpecial(tt.input)

			if gotDebug != tt.wantDebug {
				t.Errorf("ExtractSpecial() debug = %v, want %v", gotDebug, tt.wantDebug)
			}

			if gotDryRun != tt.wantDryRun {
				t.Errorf("ExtractSpecial() dryRun = %v, want %v", gotDryRun, tt.wantDryRun)
			}

			if !reflect.DeepEqual(gotForward, tt.wantForward) {
				t.Errorf("ExtractSpecial() forward = %v, want %v", gotForward, tt.wantForward)
			}
		})
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "simple tokens",
			input:    "arg1 arg2 arg3",
			expected: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "double quoted string",
			input:    `arg1 "arg with spaces" arg3`,
			expected: []string{"arg1", "arg with spaces", "arg3"},
		},
		{
			name:     "single quoted string",
			input:    `arg1 'arg with spaces' arg3`,
			expected: []string{"arg1", "arg with spaces", "arg3"},
		},
		{
			name:     "mixed quotes",
			input:    `arg1 "double quoted" 'single quoted' normal`,
			expected: []string{"arg1", "double quoted", "single quoted", "normal"},
		},
		{
			name:     "escaped quotes in double quotes",
			input:    `"escaped \"quote\" inside"`,
			expected: []string{`escaped "quote" inside`},
		},
		{
			name:     "escaped quotes in single quotes",
			input:    `'escaped \'quote\' inside'`,
			expected: []string{`escaped 'quote' inside`},
		},
		{
			name:     "multiple spaces",
			input:    "arg1    arg2     arg3",
			expected: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "tabs and spaces",
			input:    "arg1\t\targ2  \t arg3",
			expected: []string{"arg1", "arg2", "arg3"},
		},
		{
			name:     "empty quoted strings",
			input:    `arg1 "" '' arg4`,
			expected: []string{"arg1", "", "", "arg4"},
		},
		{
			name:     "complex JVM args",
			input:    `-Xmx512m "-Djava.library.path=./libs with spaces" -Dfile.encoding=UTF-8`,
			expected: []string{"-Xmx512m", "-Djava.library.path=./libs with spaces", "-Dfile.encoding=UTF-8"},
		},
		{
			name:     "quoted paths with backslashes",
			input:    `"-Djava.library.path=C:\Program Files\Java\libs"`,
			expected: []string{`-Djava.library.path=C:\Program Files\Java\libs`},
		},
		{
			name:     "only whitespace",
			input:    "   \t  \n  ",
			expected: []string{},
		},
		{
			name:     "trailing and leading spaces",
			input:    "  arg1  arg2  ",
			expected: []string{"arg1", "arg2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Tokenize(tt.input)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Tokenize(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTokenizeEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		desc     string
	}{
		{
			name:     "unterminated double quote",
			input:    `arg1 "unterminated quote`,
			expected: []string{"arg1", "unterminated quote"},
			desc:     "should handle unterminated quotes gracefully",
		},
		{
			name:     "unterminated single quote",
			input:    `arg1 'unterminated quote`,
			expected: []string{"arg1", "unterminated quote"},
			desc:     "should handle unterminated quotes gracefully",
		},
		{
			name:     "nested different quotes",
			input:    `"outer 'inner' quote"`,
			expected: []string{"outer 'inner' quote"},
			desc:     "should handle nested different quote types",
		},
		{
			name:     "backslash before quote",
			input:    `"path\\with\"quote"`,
			expected: []string{`path\with"quote`},
			desc:     "should handle backslash before quote correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Tokenize(tt.input)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Tokenize(%q) = %v, want %v (%s)", tt.input, result, tt.expected, tt.desc)
			}
		})
	}
}
