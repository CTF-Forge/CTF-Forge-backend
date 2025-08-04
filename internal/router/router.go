package router

import (
	"github.com/Saku0512/CTFLab/ctflab/config"
	"github.com/Saku0512/CTFLab/ctflab/internal/handler"
	"github.com/Saku0512/CTFLab/ctflab/internal/repository"
	"github.com/Saku0512/CTFLab/ctflab/internal/service"
	"github.com/Saku0512/CTFLab/ctflab/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// セッションストア
	gothic.Store = sessions.NewCookieStore([]byte(config.GetSessionSecret()))

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)
	oauthRepo := repository.NewOAuthAccountRepository(db)
	challengeRepo := repository.NewChallengeRepository(db)

	// JWTマネージャーの初期化
	jwtManager := token.NewJWTManager(
		config.GetJWTAccessSecret(),
		config.GetJWTRefreshSecret(),
		config.GetJWTIssuer(),
		config.GetJWTAccessExpireDuration(),
		config.GetJWTRefreshExpireDuration(),
	)

	// サービスの初期化
	authService := service.NewAuthService(userRepo, jwtManager)
	oauthService := service.NewOAuthService(oauthRepo, userRepo, jwtManager)
	challengeService := service.NewChallengeService(challengeRepo, userRepo)

	// ハンドラーの初期化
	authHandler := handler.NewAuthHandler(authService)
	oauthHandler := handler.NewOAuthHandler(oauthService, jwtManager)
	challengeHandler := handler.NewChallengeHandler(challengeService)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 認証APIグループ
	authGroup := r.Group("/auth")
	{
		// Email認証
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.POST("/logout", authHandler.Logout)

		// OAuth認証
		authGroup.GET("/:provider", oauthHandler.BeginAuthHandler)
		authGroup.GET("/:provider/callback", oauthHandler.CallbackAuthHandler)
		authGroup.POST("/oauth/refresh", oauthHandler.RefreshTokenHandler)
		authGroup.POST("/oauth/logout", oauthHandler.LogoutHandler)
	}

	// 保護されたAPIグループ（認証が必要）
	protectedGroup := r.Group("/api")
	protectedGroup.Use(token.AuthMiddleware(jwtManager))
	{
		// ユーザー関連
		protectedGroup.GET("/me", authHandler.Me)
		protectedGroup.POST("/challenges", challengeHandler.CreateChallenge)
		protectedGroup.GET("/challenges", challengeHandler.CollectChallengesByUsername)
		protectedGroup.GET("/challenges/:challengeId", challengeHandler.GetChallenge)
		protectedGroup.PUT("/challenges/:challengeId", challengeHandler.UpdateChallenge)
		protectedGroup.DELETE("/challenges/:challengeId", challengeHandler.DeleteChallenge)
		// ここに他の保護されたエンドポイントを追加
		// 例: 問題作成、提出履歴など
	}

	// 公開APIグループ（認証オプショナル）
	publicGroup := r.Group("/api/public")
	publicGroup.Use(token.OptionalAuthMiddleware(jwtManager))
	{
		// 問題一覧など、認証されていないユーザーもアクセス可能なエンドポイント
		publicGroup.GET("/challenges", func(c *gin.Context) {
			// 認証されている場合はユーザー情報を含める
			if userID, exists := token.GetUserID(c); exists {
				c.JSON(200, gin.H{
					"message": "public challenges",
					"user_id": userID,
				})
			} else {
				c.JSON(200, gin.H{
					"message": "public challenges",
				})
			}
		})
		publicGroup.GET("/challenges/:challengeId", challengeHandler.GetPublicChallenge)
	}

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}