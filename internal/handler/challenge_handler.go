package handler

import (
	"context"
	"net/http"

	"github.com/Saku0512/CTFLab/ctflab/internal/models"
	"github.com/Saku0512/CTFLab/ctflab/internal/service"
	"github.com/Saku0512/CTFLab/ctflab/pkg/token"
	"github.com/gin-gonic/gin"
)

// CreateChallengeRequestは問題作成APIのリクエストボディを定義します。
// データベースモデルと分離することで、柔軟に対応できます。
type CreateChallengeRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"` // カテゴリー名を文字列として受け取ります
	Score       int    `json:"score" binding:"required"`
	Flag        string `json:"flag" binding:"required"`
}

type ChallengeHandler struct {
	service service.ChallengeService
}

func NewChallengeHandler(service service.ChallengeService) *ChallengeHandler {
	return &ChallengeHandler{service: service}
}

// CreateChallengeは、認証されたユーザーが新しい問題を作成するためのハンドラです。
func (h *ChallengeHandler) CreateChallenge(c *gin.Context) {
	var req CreateChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 認証されたユーザーのIDをコンテキストから取得
	userID, exists := token.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// リクエストからモデルを構築
	challenge := &models.Challenge{
		UserID:      uint(userID),
		Title:       req.Title,
		Description: req.Description,
		Score:       req.Score,
		Flag:        req.Flag,
	}

	// サービスを呼び出して問題を作成し、カテゴリー名を渡します
	if err := h.service.CreateChallenge(context.Background(), challenge, req.Category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create challenge: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Challenge created successfully"})
}
