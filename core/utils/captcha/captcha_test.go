package captcha

import (
	"testing"
)

func TestCreateCaptcha(t *testing.T) {
	resp, err := CreateCaptcha()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil {
		t.Fatalf("expected response, got nil")
	}
	if resp.CaptchaID == "" {
		t.Errorf("expected non-empty CaptchaID")
	}
	if resp.ImagePath == "" {
		t.Errorf("expected non-empty ImagePath")
	}
}

func TestVerifyCode(t *testing.T) {
	// Test missing inputs
	if result := VerifyCode("", "test"); result != "ErrCaptchaCode" {
		t.Errorf("expected ErrCaptchaCode for empty codeID, got %v", result)
	}
	if result := VerifyCode("id", ""); result != "ErrCaptchaCode" {
		t.Errorf("expected ErrCaptchaCode for empty code, got %v", result)
	}

	codeID := "test_id"

	// Test correct code
	store.Set(codeID, "correct_code")
	if result := VerifyCode(codeID, "correct_code"); result != "" {
		t.Errorf("expected empty string for correct code, got %v", result)
	}

	// Test incorrect code
	store.Set(codeID, "correct_code")
	if result := VerifyCode(codeID, "wrong_code"); result != "ErrCaptchaCode" {
		t.Errorf("expected ErrCaptchaCode for wrong code, got %v", result)
	}

	// Test case-insensitivity
	store.Set(codeID, "CoDe")
	if result := VerifyCode(codeID, "cOdE"); result != "" {
		t.Errorf("expected empty string for case-insensitive match, got %v", result)
	}

	// Test whitespace trimming
	store.Set(codeID, " code ")
	if result := VerifyCode(codeID, "  code  "); result != "" {
		t.Errorf("expected empty string for whitespace trimmed match, got %v", result)
	}
}
