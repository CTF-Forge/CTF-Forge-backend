package main

import (
	"github.com/gin-gonic/gin"

	// Swagger 関連
	_ "github.com/Saku0512/CTFLab/ctflab/docs" // ← ここを自分のモジュール名に合わせて変更
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"net/http"
)

// @title CTFLab API
// @version 1.0
// @description CTFの問題管理・共有アプリ用API仕様
// @host localhost:8080
// @BasePath /api/v1

func main() {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", PingHandler)
	}

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}

// PingHandler godoc
// @Summary Ping Example
// @Description ヘルスチェック用エンドポイント
// @Tags HealthCheck
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
