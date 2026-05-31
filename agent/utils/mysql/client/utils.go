package client

import "strings"

// EscapeIdentifier escapes MySQL identifiers such as database names, table names, usernames.
func EscapeIdentifier(s string) string {
	return strings.ReplaceAll(s, "`", "``")
}

// EscapeString escapes MySQL string literals such as passwords.
// It also escapes backslashes to prevent escaping the surrounding quotes.
func EscapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	return strings.ReplaceAll(s, "'", "''")
}
