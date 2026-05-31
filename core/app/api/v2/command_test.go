package v2

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1Panel-dev/1Panel/core/app/dto"
	"github.com/1Panel-dev/1Panel/core/i18n"
	"github.com/gin-gonic/gin"
)

type MockCommandService struct {
	ExportResult string
	ExportErr    error
}

func (m *MockCommandService) List(req dto.OperateByType) ([]dto.CommandInfo, error) {
	return nil, nil
}
func (m *MockCommandService) SearchForTree(req dto.OperateByType) ([]dto.CommandTree, error) {
	return nil, nil
}
func (m *MockCommandService) SearchWithPage(search dto.SearchCommandWithPage) (int64, interface{}, error) {
	return 0, nil, nil
}
func (m *MockCommandService) Create(req dto.CommandOperate) error {
	return nil
}
func (m *MockCommandService) Update(req dto.CommandOperate) error {
	return nil
}
func (m *MockCommandService) Delete(ids []uint) error {
	return nil
}
func (m *MockCommandService) Export() (string, error) {
	return m.ExportResult, m.ExportErr
}

func TestBaseApi_ExportCommands(t *testing.T) {
	// Initialize i18n for tests to prevent panic
	i18n.Init()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockExportFile string
		mockExportErr  error
		expectedCode   int
		expectedPath   string
		expectedMsg    string
	}{
		{
			name:           "successful export",
			mockExportFile: "/tmp/export/commands/1panel-commands-20231027120000.csv",
			mockExportErr:  nil,
			expectedCode:   http.StatusOK,
			expectedPath:   "/tmp/export/commands/1panel-commands-20231027120000.csv",
			expectedMsg:    "",
		},
		{
			name:           "export failure",
			mockExportFile: "",
			mockExportErr:  errors.New("db error"),
			expectedCode:   http.StatusInternalServerError,
			expectedPath:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the service
			originalService := commandService
			defer func() { commandService = originalService }()

			commandService = &MockCommandService{
				ExportResult: tt.mockExportFile,
				ExportErr:    tt.mockExportErr,
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Needs a request with some method and URL just to not be nil
			c.Request, _ = http.NewRequest("POST", "/core/commands/export", bytes.NewBufferString(""))

			api := &BaseApi{}
			api.ExportCommands(c)

			var response dto.Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

            // Log response body for inspection
            t.Logf("Response body: %s", w.Body.String())

			if tt.mockExportErr == nil {
				if response.Code != http.StatusOK {
					t.Errorf("Expected code %d, got %d", http.StatusOK, response.Code)
				}
				if response.Data != tt.expectedPath {
					t.Errorf("Expected data %v, got %v", tt.expectedPath, response.Data)
				}
			} else {
				if response.Code != http.StatusInternalServerError {
					t.Errorf("Expected code %d, got %d", http.StatusInternalServerError, response.Code)
				}
			}
		})
	}
}
