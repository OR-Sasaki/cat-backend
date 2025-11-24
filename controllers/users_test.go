package controllers

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    "github.com/OR-Sasaki/cat-backend/config"
    "github.com/OR-Sasaki/cat-backend/models"
    "github.com/OR-Sasaki/cat-backend/test"
)

func TestUserRegister(t *testing.T) {
    testUserLogin(t)
}

func testUserLogin(t *testing.T) {
    testDatas := userLoginTestDatas()

    test.TestApi(t, testDatas, "/api/users/login", UsersRouter, false, "")
}

func userLoginTestDatas() []test.TestData[UserLoginResponse] {
    testDatas := []test.TestData[UserLoginResponse]{}

    {
        passwordHash, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
        testDatas = append(testDatas, test.TestData[UserLoginResponse]{
            Before: func(t *testing.T, _ *models.User) map[string]any {
                user := models.User{
                    Name:         "testuser",
                    PasswordHash: string(passwordHash),
                }
                gorm.G[models.User](config.DB).Create(t.Context(), &user)
                return nil
            },
            Name: "正常系: 有効なidとpasswordでログイン",
            RequestBody: func(_ map[string]any) map[string]any {
                return map[string]any{
                    "username": "testuser",
                    "password": "testpassword",
                }
            },
            ExpectedStatus: http.StatusOK,
            ValidateResponse: func(t *testing.T, w *httptest.ResponseRecorder, responseBody UserLoginResponse) {
                assert.NotEmpty(t, responseBody.Token, "トークンは空でない必要があります")
            },
        })
    }

    testDatas = append(testDatas, []test.TestData[UserLoginResponse]{
        {
            Name: "異常系: idが存在しない",
            RequestBody: func(_ map[string]any) map[string]any {
                return map[string]interface{}{
                    "password": "testpassword"}
            },
            ExpectedStatus: http.StatusBadRequest,
        },
        {
            Name: "異常系: passwordが存在しない",
            RequestBody: func(_ map[string]any) map[string]any {
                return map[string]interface{}{
                    "id": 1}
            },
            ExpectedStatus: http.StatusBadRequest,
        },
        {
            Name:           "異常系: idとpasswordが存在しない",
            RequestBody:    nil,
            ExpectedStatus: http.StatusBadRequest,
        },
    }...)

    return testDatas
}
