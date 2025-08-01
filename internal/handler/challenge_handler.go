package handler

import (
	"context"
	"net/http"

	"github.com/Saku0512/CTFLab/ctflab/internal/handler/dtos"
	"github.com/Saku0512/CTFLab/ctflab/internal/models"
	"github.com/Saku0512/CTFLab/ctflab/internal/service"
	"github.com/Saku0512/CTFLab/ctflab/pkg/token"
	"github.com/gin-gonic/gin"
)

type ChallengeHandler struct {
	service service.ChallengeService
}

func NewChallengeHandler(service service.ChallengeService) *ChallengeHandler {
	return &ChallengeHandler{service: service}
}

// @Summary 新しい問題を作成
// @Description 認証されたユーザーが新しい問題を作成します
// @Tags challenges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param challenge body dtos.CreateChallengeRequest true "問題作成情報"
// @Success 201 {object} dtos.ChallengeCreateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/challenges [post]
func (h *ChallengeHandler) CreateChallenge(c *gin.Context) {
	var req dtos.CreateChallengeRequest
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
