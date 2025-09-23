package quote

import (
    "strings"
)

// QuoteWindowsArg quotes a single argument per Windows command-line parsing rules for display.
// It is intended for accurate display (debug/dry-run), not for shell invocation.
func QuoteWindowsArg(s string) string {
    if s == "" {
        return "\"\""
    }
    // If no spaces, tabs, quotes, or special characters, return as-is
    if !strings.ContainsAny(s, " \t\"\n\r&|<>();") {
        return s
    }
    var b strings.Builder
    b.WriteByte('"')
    for i := 0; i < len(s); i++ {
        c := s[i]
        if c == '"' {
            // Double the quote for Windows
            b.WriteByte('"')
            b.WriteByte('"')
        } else if c == '\t' {
            b.WriteString("\\t")
        } else if c == '\n' {
            b.WriteString("\\n")
        } else if c == '\r' {
            b.WriteString("\\r")
        } else {
            b.WriteByte(c)
        }
    }
    b.WriteByte('"')
    return b.String()
}

// JoinWindows joins argv into a single command string suitable for display on Windows.
func JoinWindows(argv []string) string {
    quoted := make([]string, len(argv))
    for i, a := range argv {
        quoted[i] = QuoteWindowsArg(a)
    }
    return strings.Join(quoted, " ")
}

