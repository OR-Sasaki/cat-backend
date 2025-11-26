package authenticate

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/OR-Sasaki/cat-backend/config"
	"github.com/OR-Sasaki/cat-backend/models"
)

func WithAuth(routerGroup *gin.RouterGroup, method, path string, handler func(*gin.Context, *models.User)) {
	routerGroup.Handle(method, path, func(c *gin.Context) {
		user, err := Authenticate(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			return
		}
		handler(c, user)
	})
}

func GETWithAuth(routerGroup *gin.RouterGroup, path string, handler func(*gin.Context, *models.User)) {
	WithAuth(routerGroup, "GET", path, handler)
}

func POSTWithAuth(routerGroup *gin.RouterGroup, path string, handler func(*gin.Context, *models.User)) {
	WithAuth(routerGroup, "POST", path, handler)
}

type authenticateClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func Authenticate(c *gin.Context) (*models.User, error) {
	headerToken := c.GetHeader("Authorization")
	if headerToken == "" {
		return nil, errors.New("no token found")
	}

	headerToken = strings.TrimPrefix(headerToken, "Bearer ")

	// 1. JWTトークンを検証
	claims := &authenticateClaims{}
	parsedToken, err := jwt.ParseWithClaims(headerToken, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.JWTSecret), nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, errors.New("token is invalid")
	}

	// 2. トークンからユーザーIDを取得
	// 3. ユーザーIDを使用してユーザーを取得
	user, err := models.GetUser(c.Request.Context(), claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 4. ユーザーを返す
	return user, nil
}

func GenerateAuthenticateToken(userID uint) (string, error) {
	expiration := time.Now().Add(time.Hour * time.Duration(config.JWTExpirationHours))
	claims := authenticateClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JWTSecret))
}
