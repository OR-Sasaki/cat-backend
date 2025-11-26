package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OR-Sasaki/cat-backend/authenticate"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/OR-Sasaki/cat-backend/config"
	"github.com/OR-Sasaki/cat-backend/models"
)

type TestData[T any] struct {
	Before           func(*testing.T, *models.User) map[string]any
	Name             string
	RequestBody      interface{}
	ExpectedStatus   int
	ValidateResponse func(t *testing.T, response *httptest.ResponseRecorder, responseBody T)
}

func TestApi[T any](
	t *testing.T, tests []TestData[T],
	urlPath string,
	routerGroupAttacher func(routerGroup *gin.RouterGroup),
	useAuth bool,
	method string,
) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)
	router := gin.Default()
	routerGroup := router.Group("/api")
	routerGroupAttacher(routerGroup)

	// メソッドが指定されていない場合はPOSTをデフォルトとする
	if method == "" {
		method = http.MethodPost
	}

	for _, tt := range tests {
		t.Log(tt.Name)
		t.Run(tt.Name, func(t *testing.T) {
			// userレコード作成
			var user *models.User
			if useAuth {
				user, _, _ = models.RegisterUser(t.Context(), "testuser")
			}

			var beforeParams map[string]any

			if tt.Before != nil {
				beforeParams = tt.Before(t, user)
			}

			requestBody := resolveRequestBody(tt.RequestBody, beforeParams)

			req := createTestRequest(t, method, urlPath, requestBody)

			if useAuth {
				setAuthHeader(req, user)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedStatus, w.Code)

			var responseBody T
			shouldUnmarshal := w.Code >= http.StatusOK && w.Code < http.StatusMultipleChoices
			if !shouldUnmarshal && tt.ValidateResponse != nil {
				shouldUnmarshal = true
			}

			if shouldUnmarshal {
				if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
					t.Fatalf("レスポンスのアンマーシャルに失敗しました: %v", err)
				}
			}

			if tt.ValidateResponse != nil {
				tt.ValidateResponse(t, w, responseBody)
			}
		})
	}
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("テストデータベースへの接続に失敗しました: %v", err)
	}

	config.DB = db
	if err := models.Migrate(); err != nil {
		t.Fatalf("マイグレーションの実行に失敗しました: %v", err)
	}

	return db
}

func buildURLWithQueryParams(urlPath string, requestBody interface{}) string {
	if requestBody == nil {
		return urlPath
	}

	queryParams, ok := requestBody.(map[string]string)
	if !ok {
		return urlPath
	}

	if len(queryParams) == 0 {
		return urlPath
	}

	values := make([]string, 0, len(queryParams)*2)
	for key, value := range queryParams {
		values = append(values, key+"="+value)
	}
	return urlPath + "?" + strings.Join(values, "&")
}

func createTestRequest(t *testing.T, method, urlPath string, requestBody interface{}) *http.Request {
	var req *http.Request
	var err error

	if method == http.MethodGet {
		fullURL := buildURLWithQueryParams(urlPath, requestBody)
		req, err = http.NewRequest(method, fullURL, nil)
	} else {
		body, err := json.Marshal(requestBody)
		if err != nil {
			t.Fatalf("リクエストボディのマーシャルに失敗しました: %v", err)
		}
		req, err = http.NewRequest(method, urlPath, bytes.NewBuffer(body))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	if err != nil {
		t.Fatalf("リクエストの作成に失敗しました: %v", err)
	}

	return req
}

func setAuthHeader(req *http.Request, user *models.User) {
	token, _ := authenticate.GenerateAuthenticateToken(user.ID)
	req.Header.Set("Authorization", "Bearer "+token)
}

func resolveRequestBody(body interface{}, params map[string]any) interface{} {
	if body == nil {
		return nil
	}

	switch supplier := body.(type) {
	case func(map[string]any) interface{}:
		return supplier(params)
	case func() interface{}:
		return supplier()
	default:
		return body
	}
}
