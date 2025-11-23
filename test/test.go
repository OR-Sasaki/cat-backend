package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OR-Sasaki/cat-backend/config"
	"github.com/OR-Sasaki/cat-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestData[T any] struct {
	Before           func(t *testing.T)
	Name             string
	RequestBody      interface{}
	ExpectedStatus   int
	ValidateResponse func(t *testing.T, response *httptest.ResponseRecorder, responseBody T)
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

func TestApi[T any](t *testing.T, tests []TestData[T], urlPath string, routerGroupAttacher func(routerGroup *gin.RouterGroup)) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)
	router := gin.Default()
	routerGroup := router.Group("/api")
	routerGroupAttacher(routerGroup)

	for _, tt := range tests {
		t.Log(tt.Name)
		t.Run(tt.Name, func(t *testing.T) {

			if tt.Before != nil {
				tt.Before(t)
			}

			body, err := json.Marshal(tt.RequestBody)
			if err != nil {
				t.Fatalf("リクエストボディのマーシャルに失敗しました: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("リクエストの作成に失敗しました: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.ExpectedStatus, w.Code)

			var responseBody T
			if err := json.Unmarshal(w.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("レスポンスのアンマーシャルに失敗しました: %v", err)
			}

			if tt.ValidateResponse != nil {
				tt.ValidateResponse(t, w, responseBody)
			}
		})
	}
}
