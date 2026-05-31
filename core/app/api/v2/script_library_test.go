package v2

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1Panel-dev/1Panel/core/app/dto"
	"github.com/1Panel-dev/1Panel/core/global"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	i18nv2 "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// MockScriptService implements service.IScriptService
type MockScriptService struct {
	CreateFunc func(req dto.ScriptOperate) error
}

func (m *MockScriptService) Run() {}

func (m *MockScriptService) Search(ctx *gin.Context, req dto.SearchPageWithGroup) (int64, interface{}, error) {
	return 0, nil, nil
}

func (m *MockScriptService) Create(req dto.ScriptOperate) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(req)
	}
	return nil
}

func (m *MockScriptService) Update(req dto.ScriptOperate) error {
	return nil
}

func (m *MockScriptService) Delete(ids dto.OperateByIDs) error {
	return nil
}

func (m *MockScriptService) Sync(req dto.OperateByTaskID) error {
	return nil
}

func TestBaseApi_CreateScript(t *testing.T) {
	gin.SetMode(gin.TestMode)
	global.VALID = validator.New()

	bundle := i18nv2.NewBundle(language.English)
	global.I18n = i18nv2.NewLocalizer(bundle, "en")

	api := &BaseApi{}

	t.Run("success", func(t *testing.T) {
		mockService := &MockScriptService{
			CreateFunc: func(req dto.ScriptOperate) error {
				if req.Name != "test-script" {
					t.Errorf("expected name 'test-script', got %s", req.Name)
				}
				return nil
			},
		}
		scriptService = mockService

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := dto.ScriptOperate{
			Name:   "test-script",
			Script: "echo 'hello'",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/core/script", bytes.NewBuffer(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		api.CreateScript(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected HTTP status %d, got %d", http.StatusOK, w.Code)
		}

		var resp dto.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != http.StatusOK {
			t.Errorf("expected response code %d, got %d", http.StatusOK, resp.Code)
		}
	})

	t.Run("bind_error", func(t *testing.T) {
		scriptService = &MockScriptService{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Invalid JSON
		c.Request, _ = http.NewRequest(http.MethodPost, "/core/script", bytes.NewBufferString("{invalid_json}"))
		c.Request.Header.Set("Content-Type", "application/json")

		api.CreateScript(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected HTTP status %d, got %d", http.StatusOK, w.Code)
		}

		var resp dto.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != http.StatusBadRequest {
			t.Errorf("expected response code %d, got %d", http.StatusBadRequest, resp.Code)
		}
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := &MockScriptService{
			CreateFunc: func(req dto.ScriptOperate) error {
				return errors.New("service error")
			},
		}
		scriptService = mockService

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := dto.ScriptOperate{
			Name: "test-script",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/core/script", bytes.NewBuffer(bodyBytes))
		c.Request.Header.Set("Content-Type", "application/json")

		api.CreateScript(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected HTTP status %d, got %d", http.StatusOK, w.Code)
		}

		var resp dto.Response
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != http.StatusInternalServerError {
			t.Errorf("expected response code %d, got %d", http.StatusInternalServerError, resp.Code)
		}
	})
}
