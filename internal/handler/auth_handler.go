package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/CTF-Forge/CTF-Forge-backend/internal/models"
	"github.com/CTF-Forge/CTF-Forge-backend/internal/service"
)

// Swag用のレスポンス型定義
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type MessageResponse struct {
	Message string `json:"message" example:"success message"`
}

type TokenResponse struct {
	Message      string `json:"message" example:"login successful"`
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int64  `json:"expires_in" example:"3600"`
}

type OAuthResponse struct {
	Message      string `json:"message" example:"oauth authentication successful"`
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int64  `json:"expires_in" example:"3600"`
	User         struct {
		Username string `json:"username" example:"testuser"`
		Email    string `json:"email" example:"test@example.com"`
		Provider string `json:"provider" example:"github"`
	} `json:"user"`
}

type AuthHandler struct {
	authService *service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// リクエスト構造体
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20" example:"testuser"`
	Email    string `json:"email" validate:"required,email" example:"test@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"test@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Register godoc
// @Summary      ユーザー登録
// @Description  新規ユーザーを登録します
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  RegisterRequest  true  "登録情報"
// @Success      201   {object}  MessageResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPwd, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "password hashing failed"})
		return
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPwd,
		CreatedAt:    time.Now(),
	}

	if err := h.authService.RegisterUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// Login godoc
// @Summary      ユーザーログイン
// @Description  メールとパスワードでログインします
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  LoginRequest  true  "ログイン情報"
// @Success      200   {object}  TokenResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := h.authService.Login(context.Background(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
		return
	}

	var user *models.User
	user, err = h.authService.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "login successful",
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"expires_in":    tokenPair.ExpiresIn,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// RefreshToken godoc
// @Summary      トークン更新
// @Description  リフレッシュトークンを使用して新しいアクセストークンを取得します
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  RefreshTokenRequest  true  "リフレッシュトークン"
// @Success      200   {object}  TokenResponse
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := h.authService.RefreshToken(req.RefreshToken)
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

// Logout godoc
// @Summary      ログアウト
// @Description  ユーザーログアウトします
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200   {object}  MessageResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// クライアント側でトークンを削除することを前提とする
	// 必要に応じてブラックリストにトークンを追加する処理を追加可能
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// MeResponse swag用
// @Description 自分のユーザー情報
// @Tags user
// @Success 200 {object} MeResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/me [get]
type MeResponse struct {
	UserID    uint      `json:"user_id" example:"1"`
	Username  string    `json:"username" example:"testuser"`
	Email     string    `json:"email" example:"test@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2024-08-01T12:34:56Z"`
}

// Me godoc
// @Summary      自分のユーザー情報取得
// @Description  JWT認証ユーザーの情報を返す
// @Tags         user
// @Security     bearer
// @Produce      json
// @Success      200   {object}  MeResponse
// @Failure      401   {object}  ErrorResponse
// @Router       /api/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid user id"})
		return
	}
	user, err := h.authService.GetUserByID(c.Request.Context(), id)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "user not found"})
		return
	}
	c.JSON(http.StatusOK, MeResponse{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}
