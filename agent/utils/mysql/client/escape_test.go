package client

import (
	"testing"
)

func TestEscape(t *testing.T) {
	str := "test'\\\"\x00\n\r\x1a"
	expected := "test\\'\\\\\\\"\\0\\n\\r\\Z"
	if escaped := Escape(str); escaped != expected {
		t.Errorf("expected %q, got %q", expected, escaped)
	}
}

func TestEscapeIdentifier(t *testing.T) {
	str := "my`db"
	expected := "my``db"
	if escaped := EscapeIdentifier(str); escaped != expected {
		t.Errorf("expected %q, got %q", expected, escaped)
	}
}
