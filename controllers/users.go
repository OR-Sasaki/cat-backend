package controllers

import (
	"log/slog"
	"net/http"

	"github.com/OR-Sasaki/cat-backend/authenticate"
	"github.com/gin-gonic/gin"

	"github.com/OR-Sasaki/cat-backend/models"
)

func UsersRouter(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.POST("/register", UserRegister)
		users.POST("/login", UserLogin)
	}
}

// **************************************************
// UserRegister
// **************************************************

type UserRegisterRequest struct {
	Name string `json:"name" binding:"required,min=4,max=20"`
}

type UserRegisterResponse struct {
	ID       uint   `json:"id"`
	Password string `json:"password"`
}

// @Summary		ユーザー登録
// @Description	新規ユーザーを登録し、IDとパスワードを返す
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			request	body		UserRegisterRequest	true	"ユーザー登録リクエスト"
// @Success		200		{object}	UserRegisterResponse
// @Router			/users/register [post]
func UserRegister(c *gin.Context) {
	var request UserRegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		slog.Error("failed to bind request parameters", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "parameter error"})
		return
	}

	user, password, err := models.RegisterUser(c.Request.Context(), request.Name)
	if err != nil {
		slog.Error("failed to register user", "error", err, "name", request.Name)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, UserRegisterResponse{
		ID:       user.ID,
		Password: password,
	})
}

// **************************************************
// UserLogin
// **************************************************

type UserLoginRequest struct {
	ID       uint   `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

// @Summary		ユーザーログイン
// @Description	IDとパスワードで認証し、JWTトークンを返す。このJWTトークンは Authorization: Bearer <token> としてヘッダーに付加してください。
// @Tag			users
// @Accep			json
// @Produce		json
// @Param			request	body		UserLoginRequest	true	"ユーザーログインリクエスト"
// @Success		200		{object}	UserLoginResponse
// @Router			/users/login [post]
func UserLogin(c *gin.Context) {
	var request UserLoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		slog.Error("failed to bind request parameters", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "parameter error"})
		return
	}

	// IDでユーザーを検索
	user, err := models.GetUser(c.Request.Context(), request.ID)
	if err != nil {
		slog.Error("user not found", "error", err, "id", request.ID)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	// パスワードを検証
	if !user.VerifyPassword(request.Password) {
		slog.Warn("password verification failed", "id", request.ID)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	// ログイントークンを生成
	token, err := authenticate.GenerateAuthenticateToken(user.ID)
	if err != nil {
		slog.Error("failed to generate token", "error", err, "id", request.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{Token: token})
}
