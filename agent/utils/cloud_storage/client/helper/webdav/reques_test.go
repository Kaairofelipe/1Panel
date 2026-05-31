package webdav

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Propfind_Compression(t *testing.T) {
	xmlResponse := `<?xml version="1.0" encoding="utf-8" ?>
<d:multistatus xmlns:d="DAV:">
  <d:response>
    <d:href>/test/</d:href>
    <d:propstat>
      <d:prop>
        <d:displayname>test</d:displayname>
        <d:resourcetype><d:collection/></d:resourcetype>
      </d:prop>
      <d:status>HTTP/1.1 200 OK</d:status>
    </d:propstat>
  </d:response>
</d:multistatus>`

	tests := []struct {
		name            string
		contentEncoding string
		compress        func(string) []byte
	}{
		{
			name:            "No compression",
			contentEncoding: "",
			compress: func(s string) []byte {
				return []byte(s)
			},
		},
		{
			name:            "Gzip compression",
			contentEncoding: "gzip",
			compress: func(s string) []byte {
				var buf bytes.Buffer
				gw := gzip.NewWriter(&buf)
				gw.Write([]byte(s))
				gw.Close()
				return buf.Bytes()
			},
		},
		{
			name:            "Deflate compression",
			contentEncoding: "deflate",
			compress: func(s string) []byte {
				var buf bytes.Buffer
				zw := zlib.NewWriter(&buf)
				zw.Write([]byte(s))
				zw.Close()
				return buf.Bytes()
			},
		},
		{
			name:            "Gzip compression (case insensitive)",
			contentEncoding: "GZIP",
			compress: func(s string) []byte {
				var buf bytes.Buffer
				gw := gzip.NewWriter(&buf)
				gw.Write([]byte(s))
				gw.Close()
				return buf.Bytes()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PROPFIND" {
					t.Errorf("Expected method PROPFIND, got %s", r.Method)
				}
				if r.Header.Get("Accept-Encoding") != "gzip, deflate;q=0.8,q=0.7" {
					t.Errorf("Expected Accept-Encoding gzip, deflate;q=0.8,q=0.7, got %s", r.Header.Get("Accept-Encoding"))
				}

				w.Header().Set("Content-Type", "application/xml;charset=UTF-8")
				if tt.contentEncoding != "" {
					w.Header().Set("Content-Encoding", tt.contentEncoding)
				}
				w.WriteHeader(http.StatusMultiStatus)
				w.Write(tt.compress(xmlResponse))
			}))
			defer ts.Close()

			c := NewClient(ts.URL, "user", "pass")

			var foundName string
			parse := func(resp interface{}) error {
				r := resp.(*response)
				if p := getProps(r, "200"); p != nil {
					foundName = p.Name
				}
				return nil
			}

			err := c.propfind("/", true, template, &response{}, parse)
			if err != nil {
				t.Fatalf("propfind failed: %v", err)
			}

			if foundName != "test" {
				t.Errorf("Expected displayname 'test', got '%s'", foundName)
			}
		})
	}
}

func TestClient_Propfind_Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "pass")
	err := c.propfind("/", true, template, &response{}, func(resp interface{}) error { return nil })
	if err == nil {
		t.Fatal("Expected error for 500 status code, got nil")
	}

	expectedErr := fmt.Sprintf("PROPFIND /: 500")
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%v'", expectedErr, err)
	}
}
