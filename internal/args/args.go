// Package args provides command-line argument parsing and tokenization for the GJG launcher.
package args

import "strings"

// ExtractSpecial parses special flags and returns debug, dryRun, and remaining args.
func ExtractSpecial(in []string) (debug, dryRun bool, rest []string) {
	rest = make([]string, 0, len(in))
	for _, a := range in {
		switch a {
		case "--gjg-debug":
			debug = true
			continue
		case "--gjg-dry-run":
			dryRun = true
			continue
		}
		rest = append(rest, a)
	}
	return
}

// Tokenize splits a command-line string into arguments, supporting quotes and escapes.
// Supports single ('), double (") quotes, and backslash escaping within quoted sections.
func Tokenize(s string) []string {
	if strings.TrimSpace(s) == "" {
		return []string{}
	}

	var out []string
	var b strings.Builder
	in := []rune(s)
	i := 0
	const (
		none = 0
		sq   = 1
		dq   = 2
	)
	mode := none
	inQuotes := false

	flush := func(force bool) {
		if b.Len() > 0 || force {
			out = append(out, b.String())
			b.Reset()
		}
	}

	for i < len(in) {
		r := in[i]
		switch mode {
		case none:
			i = handleNoneMode(r, i, in, &mode, &inQuotes, &b, flush)
		case sq:
			i = handleSingleQuoteMode(r, i, in, &mode, &inQuotes, &b, flush)
		case dq:
			i = handleDoubleQuoteMode(r, i, in, &mode, &inQuotes, &b, flush)
		}
		i++
	}

	// Flush last token
	if b.Len() > 0 || inQuotes {
		flush(false)
	}

	return out
}

// handleNoneMode processes characters when not inside quotes
func handleNoneMode(r rune, i int, in []rune, mode *int, inQuotes *bool, b *strings.Builder, flush func(bool)) int {
	const (
		none = 0
		sq   = 1
		dq   = 2
	)

	switch r {
	case ' ', '\t', '\n', '\r':
		// Skip whitespace and flush token
		if b.Len() > 0 || *inQuotes {
			flush(false)
			*inQuotes = false
		}
		// Skip multiple whitespace
		for i+1 < len(in) && isWhitespace(in[i+1]) {
			i++
		}
	case '\'':
		*mode = sq
		*inQuotes = true
	case '"':
		*mode = dq
		*inQuotes = true
	case '\\':
		// backslash outside quotes is literal
		b.WriteRune(r)
	default:
		b.WriteRune(r)
	}
	return i
}

// handleSingleQuoteMode processes characters inside single quotes
func handleSingleQuoteMode(
	r rune, i int, in []rune, mode *int, inQuotes *bool, b *strings.Builder, flush func(bool),
) int {
	const none = 0

	if r == '\'' {
		*mode = none
		flush(true) // Force flush for empty quoted strings
		*inQuotes = false
	} else if r == '\\' && i+1 < len(in) && in[i+1] == '\'' {
		// Handle escaped single quote
		i++
		b.WriteByte('\'')
	} else {
		b.WriteRune(r)
	}
	return i
}

// handleDoubleQuoteMode processes characters inside double quotes
func handleDoubleQuoteMode(
	r rune, i int, in []rune, mode *int, inQuotes *bool, b *strings.Builder, flush func(bool),
) int {
	const none = 0

	if r == '"' {
		*mode = none
		flush(true) // Force flush for empty quoted strings
		*inQuotes = false
	} else if r == '\\' {
		// handle escape inside double quotes
		if i+1 < len(in) {
			next := in[i+1]
			if next == '"' || next == '\\' {
				i++
				b.WriteRune(next)
			} else {
				b.WriteRune(r)
			}
		} else {
			b.WriteRune(r)
		}
	} else {
		b.WriteRune(r)
	}
	return i
}

// isWhitespace checks if a rune is whitespace
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
