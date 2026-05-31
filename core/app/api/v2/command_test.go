package v2

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"

	"github.com/1Panel-dev/1Panel/core/app/model"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/1Panel-dev/1Panel/core/i18n"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	global.DB = db

	db.AutoMigrate(&model.Group{})
	db.AutoMigrate(&model.Command{})

	db.Create(&model.Group{
		BaseModel: model.BaseModel{ID: 1},
		Name:      "default",
		Type:      "command",
		IsDefault: true,
	})

    // Init i18n
    i18n.Init()
}

func TestBaseApi_UploadCommandCsv(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		setupRequest func() (*http.Request, error)
		expectedCode int
	}{
		{
			name: "Success",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("file", "test.csv")
				if err != nil {
					return nil, err
				}
				part.Write([]byte("name,command\ntest_cmd,ls -l\n"))
				writer.Close()

				req, err := http.NewRequest("POST", "/core/commands/import", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, err
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "NoFile",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.Close()

				req, err := http.NewRequest("POST", "/core/commands/import", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, err
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "InvalidCSV_Empty",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("file", "test.csv")
				if err != nil {
					return nil, err
				}
				part.Write([]byte(""))
				writer.Close()

				req, err := http.NewRequest("POST", "/core/commands/import", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, err
			},
			expectedCode: http.StatusOK,
		},
        {
			name: "InvalidCSV_ContentError",
			setupRequest: func() (*http.Request, error) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("file", "test.csv")
				if err != nil {
					return nil, err
				}
				part.Write([]byte("name,command\ntest_cmd"))
				writer.Close()

				req, err := http.NewRequest("POST", "/core/commands/import", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				return req, err
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, err := tt.setupRequest()
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			c.Request = req

			api := &BaseApi{}
			api.UploadCommandCsv(c)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if tt.name == "Success" {
				if response["code"].(float64) != float64(200) {
					t.Errorf("Expected code 200, got %v", response["code"])
				}
			} else if tt.name == "NoFile" {
				if response["code"].(float64) != float64(400) {
					t.Errorf("Expected code 400, got %v", response["code"])
				}
			} else if tt.name == "InvalidCSV_Empty" {
				if response["code"].(float64) != float64(400) {
					t.Errorf("Expected code 400, got %v", response["code"])
				}
			} else if tt.name == "InvalidCSV_ContentError" {
                if response["code"].(float64) != float64(400) {
					t.Errorf("Expected code 400, got %v", response["code"])
				}
            }
		})
	}
}
