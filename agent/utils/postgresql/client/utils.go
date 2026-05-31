package client

import (
	"strings"

	"github.com/jackc/pgx/v5"
)

// EscapeIdentifier safely escapes identifiers like table, database, or role names.
func EscapeIdentifier(s string) string {
	return pgx.Identifier{s}.Sanitize()
}

// EscapeString safely escapes literal strings, mostly used for passwords in SQL queries.
func EscapeString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// EscapeBashString securely escapes a string for bash execution by wrapping it in single quotes
// and replacing single quotes with '\”
func EscapeBashString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
