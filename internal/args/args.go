package args

import "strings"

// ExtractSpecial parses special flags and returns debug, dryRun, and remaining args.
func ExtractSpecial(in []string) (debug bool, dryRun bool, rest []string) {
	rest = make([]string, 0, len(in))
	for _, a := range in {
		switch a {
		case "--gjg-debug":
			debug = true
			continue
		case "--gjg-dry-run":
			dryRun = true
			debug = true
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
			switch r {
			case ' ', '\t', '\n', '\r':
				// Skip whitespace and flush token
				if b.Len() > 0 || inQuotes {
					flush(false)
					inQuotes = false
				}
				// Skip multiple whitespace
				for i+1 < len(in) && (in[i+1] == ' ' || in[i+1] == '\t' || in[i+1] == '\n' || in[i+1] == '\r') {
					i++
				}
			case '\'':
				mode = sq
				inQuotes = true
			case '"':
				mode = dq
				inQuotes = true
			case '\\':
				// backslash outside quotes is literal
				b.WriteRune(r)
			default:
				b.WriteRune(r)
			}
		case sq:
			if r == '\'' {
				mode = none
				flush(true) // Force flush for empty quoted strings
				inQuotes = false
			} else if r == '\\' && i+1 < len(in) && in[i+1] == '\'' {
				// Handle escaped single quote
				i++
				b.WriteRune('\'')
			} else {
				b.WriteRune(r)
			}
		case dq:
			if r == '"' {
				mode = none
				flush(true) // Force flush for empty quoted strings
				inQuotes = false
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
		}
		i++
	}

	// Flush last token
	if b.Len() > 0 || inQuotes {
		flush(false)
	}

	return out
}
