package handler

import (
	"context"
	"net/http"
	"strconv"

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

	c.JSON(http.StatusCreated, gin.H{
		"message": "Challenge created successfully",
		"id":      challenge.ID,
	})
}

// @Summary ユーザーが作成した問題を取得
// @Description 認証されたユーザーが作成した問題のリストを取得します
// @Tags challenges
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Challenge
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/me/challenges [get]
func (h *ChallengeHandler) CollectChallengesByUsername(c *gin.Context) {
	// 認証されたユーザー名を取得
	username, exists := token.GetUsername(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// サービス層のCollectByUsernameを呼び出してチャレンジを取得
	challenges, err := h.service.CollectByUsername(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get challenges: " + err.Error()})
		return
	}

	// 取得したチャレンジをJSONで返す
	c.JSON(http.StatusOK, challenges)
}

// @Summary 問題を更新
// @Description 既存の問題を更新します
// @Tags challenges
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param challengeId path int true "Challenge ID"
// @Param challenge body dtos.UpdateChallengeRequest true "問題更新情報"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/challenges/{challengeId} [put]
func (h *ChallengeHandler) UpdateChallenge(c *gin.Context) {
	challengeID, err := strconv.ParseUint(c.Param("challengeId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid challenge ID"})
		return
	}

	var req dtos.UpdateChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	userID, exists := token.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := h.service.UpdateChallenge(c.Request.Context(), uint(challengeID), userID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update challenge: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Challenge updated successfully"})
}
