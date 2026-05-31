package client

import "strings"

// Escape string for SQL (prevents SQL injection in identifiers and strings)
func Escape(sql string) string {
	dest := make([]byte, 0, 2*len(sql))
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		switch c {
		case '\x00':
			dest = append(dest, '\\', '0')
		case '\n':
			dest = append(dest, '\\', 'n')
		case '\r':
			dest = append(dest, '\\', 'r')
		case '\x1a':
			dest = append(dest, '\\', 'Z')
		case '\'':
			dest = append(dest, '\\', '\'')
		case '"':
			dest = append(dest, '\\', '"')
		case '\\':
			dest = append(dest, '\\', '\\')
		default:
			dest = append(dest, c)
		}
	}
	return string(dest)
}

func EscapeIdentifier(sql string) string {
	return strings.ReplaceAll(sql, "`", "``")
}
