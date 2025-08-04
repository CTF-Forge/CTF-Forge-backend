package token

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT認証ミドルウェア
func AuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorizationヘッダーからトークンを取得
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// "Bearer "プレフィックスを除去
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		// トークンを検証
		claims, err := jwtManager.VerifyAccessToken(tokenString)
		if err != nil {
			switch err {
			case ErrExpiredToken:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
			case ErrInvalidToken, ErrInvalidSignature, ErrInvalidClaims:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
			}
			c.Abort()
			return
		}

		// ユーザー情報をコンテキストに設定
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware オプショナルなJWT認証ミドルウェア（認証されていない場合も続行）
func OptionalAuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}

		claims, err := jwtManager.VerifyAccessToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// ユーザー情報をコンテキストに設定
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_claims", claims)

		c.Next()
	}
}

// GetUserID コンテキストからユーザーIDを取得
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	if id, ok := userID.(uint); ok {
		return id, true
	}
	return 0, false
}

// GetUsername コンテキストからユーザー名を取得
func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}

	if name, ok := username.(string); ok {
		return name, true
	}
	return "", false
}

// GetUserClaims コンテキストからユーザークレームを取得
func GetUserClaims(c *gin.Context) (*UserClaims, bool) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	if userClaims, ok := claims.(*UserClaims); ok {
		return userClaims, true
	}
	return nil, false
}
