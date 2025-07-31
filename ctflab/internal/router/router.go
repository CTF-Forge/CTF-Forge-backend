package router

import (
	"time"

	"github.com/Saku0512/CTFLab/ctflab/config"
	"github.com/Saku0512/CTFLab/ctflab/internal/handler"
	"github.com/Saku0512/CTFLab/ctflab/internal/repository"
	"github.com/Saku0512/CTFLab/ctflab/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository(db)
	config.LoadEnv()
	authService := service.NewAuthService(userRepo, []byte(config.GetJWTSecret()), config.GetJWTIssuer(), time.Hour*config.GetJWTExpireDuration())
	authHandler := handler.NewAuthHandler(authService)

	// 認証APIグループ
	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)

		// OAuth系のルートもここに
		// authGroup.GET("/github", handler.GithubAuthHandler)
		// authGroup.GET("/github/callback", handler.GithubCallbackHandler)
		// authGroup.GET("/google", handler.GoogleAuthHandler)
		// authGroup.GET("/google/callback", handler.GoogleCallbackHandler)
	}

	// ここに他のAPIグループを追加

	return r
}
