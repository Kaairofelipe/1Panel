package captcha

import (
	"testing"
)

func TestVerifyCode(t *testing.T) {
	tests := []struct {
		name       string
		setupStore func()
		codeID     string
		code       string
		want       string
	}{
		{
			name:       "empty codeID",
			setupStore: func() {},
			codeID:     "",
			code:       "1234",
			want:       "ErrCaptchaCode",
		},
		{
			name: "empty code",
			setupStore: func() {
				store.Set("id1", "1234")
			},
			codeID: "id1",
			code:   "",
			want:   "ErrCaptchaCode",
		},
		{
			name: "correct code",
			setupStore: func() {
				store.Set("id2", "1234")
			},
			codeID: "id2",
			code:   "1234",
			want:   "",
		},
		{
			name: "correct code with spaces",
			setupStore: func() {
				store.Set("id3", "  1234  ")
			},
			codeID: "id3",
			code:   " 1234 ",
			want:   "",
		},
		{
			name: "correct code case insensitive",
			setupStore: func() {
				store.Set("id4", "aBcD")
			},
			codeID: "id4",
			code:   "AbCd",
			want:   "",
		},
		{
			name: "incorrect code",
			setupStore: func() {
				store.Set("id5", "1234")
			},
			codeID: "id5",
			code:   "4321",
			want:   "ErrCaptchaCode",
		},
		{
			name:       "non-existent codeID",
			setupStore: func() {},
			codeID:     "id6",
			code:       "1234",
			want:       "ErrCaptchaCode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupStore()
			if got := VerifyCode(tt.codeID, tt.code); got != tt.want {
				t.Errorf("VerifyCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
