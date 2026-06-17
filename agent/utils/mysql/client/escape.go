package client

import (
	"fmt"
	"strings"
)

// Escape escapes single quotes and backslashes in a string for use in SQL.
func Escape(sql string) string {
	sql = strings.ReplaceAll(sql, "\\", "\\\\")
	sql = strings.ReplaceAll(sql, "'", "''")
	return sql
}

// EscapeIdentifier escapes backticks in a string for use as a MySQL identifier.
func EscapeIdentifier(identifier string) string {
	return strings.ReplaceAll(identifier, "`", "``")
}


// FormatQuery replaces ? placeholders with safely escaped string literals or numeric values.
func FormatQuery(query string, args ...any) string {
	if len(args) == 0 {
		return query
	}
	parts := strings.Split(query, "?")
	if len(parts) != len(args)+1 {
		return query
	}
	var b strings.Builder
	for i, arg := range args {
		b.WriteString(parts[i])
		switch v := arg.(type) {
		case string:
			b.WriteString("'" + Escape(v) + "'")
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
			b.WriteString(fmt.Sprintf("%v", v))
		default:
			b.WriteString("'" + Escape(fmt.Sprintf("%v", v)) + "'")
		}
	}
	b.WriteString(parts[len(args)])
	return b.String()
}
