package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// デフォルトミドルフェアのGinRouterを作成
	router := gin.Default()

	// シンプルなGETエンドポイントを宣言
	router.GET("/ping", func(c *gin.Context) {
		// JSON responseをreturn
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// ポート8080でサーバーを開始(デフォルト)
	// サーバーは0.0.0.0:8080をリッスンします
	if err := router.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
