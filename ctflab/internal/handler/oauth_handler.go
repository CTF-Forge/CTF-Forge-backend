package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Saku0512/CTFLab/ctflab/internal/service"
	"github.com/Saku0512/CTFLab/ctflab/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type OAuthHandler struct {
	oauthService *service.OAuthService
	jwtManager   *token.JWTManager
}

func NewOAuthHandler(oauthService *service.OAuthService, jwtManager *token.JWTManager) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		jwtManager:   jwtManager,
	}
}

// BeginAuthHandler OAuth認証を開始
// GET /auth/{provider}
func (h *OAuthHandler) BeginAuthHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()

	provider := c.Param("provider")
	if !isValidProvider(provider) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
		return
	}

	// デバッグログ
	fmt.Printf("Starting OAuth for provider: %s\n", provider)

	// プロバイダーを設定
	req := c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
	c.Request = req

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// CallbackAuthHandler OAuth認証コールバックを処理
// GET /auth/{provider}/callback
func (h *OAuthHandler) CallbackAuthHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}()

	provider := c.Param("provider")
	if !isValidProvider(provider) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
		return
	}

	// プロバイダーを設定
	req := c.Request.WithContext(context.WithValue(c.Request.Context(), "provider", provider))
	c.Request = req

	// OAuth認証を完了
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "oauth authentication failed"})
		return
	}

	// ユーザー情報を取得
	username := user.Name
	if username == "" {
		username = user.NickName
	}
	if username == "" {
		username = strings.Split(user.Email, "@")[0] // メールアドレスの@前をユーザー名として使用
	}

	// OAuthサービスでユーザー処理
	tokenPair, err := h.oauthService.HandleOAuthCallback(
		context.Background(),
		username,
		provider,
		user.UserID,
		user.Email,
		user.AccessToken,
		user.RefreshToken,
		user.ExpiresAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process oauth user"})
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, gin.H{
		"message":       "oauth authentication successful",
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"expires_in":    tokenPair.ExpiresIn,
		"user": gin.H{
			"username": username,
			"email":    user.Email,
			"provider": provider,
		},
	})
}

// RefreshTokenHandler リフレッシュトークンを使用して新しいトークンペアを生成
// POST /auth/refresh
func (h *OAuthHandler) RefreshTokenHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token is required"})
		return
	}

	tokenPair, err := h.jwtManager.RefreshTokenPair(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"expires_in":    tokenPair.ExpiresIn,
	})
}

// LogoutHandler ログアウト処理
// POST /auth/logout
func (h *OAuthHandler) LogoutHandler(c *gin.Context) {
	// クライアント側でトークンを削除することを前提とする
	// 必要に応じてブラックリストにトークンを追加する処理を追加可能
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// isValidProvider 有効なプロバイダーかどうかをチェック
func isValidProvider(provider string) bool {
	validProviders := []string{"github", "google"}
	for _, p := range validProviders {
		if p == provider {
			return true
		}
	}
	return false
}
