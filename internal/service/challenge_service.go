package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Saku0512/CTFLab/ctflab/internal/handler/dtos"
	"github.com/Saku0512/CTFLab/ctflab/internal/models"
	"github.com/Saku0512/CTFLab/ctflab/internal/repository"
)

type ChallengeService interface {
	CreateChallenge(ctx context.Context, challenge *models.Challenge, categoryName string) error
	CollectByUsername(ctx context.Context, username string) ([]*models.Challenge, error)
	UpdateChallenge(ctx context.Context, challengeID uint, userID uint, req *dtos.UpdateChallengeRequest) error
	DeleteChallenge(ctx context.Context, challengeID uint, userID uint) error
	GetChallengeByID(ctx context.Context, challengeID uint, userID uint) (*dtos.ChallengeDetailResponse, error)
}

type challengeService struct {
	challengerepo repository.ChallengeRepository
	userrepo      repository.UserRepository
}

// 以前の修正コード
func NewChallengeService(challengerepo repository.ChallengeRepository, userrepo repository.UserRepository) ChallengeService {
	return &challengeService{challengerepo: challengerepo, userrepo: userrepo}
}

// CreateChallengeは、カテゴリー名を解決して新しい問題をデータベースに保存します。
func (s *challengeService) CreateChallenge(ctx context.Context, challenge *models.Challenge, categoryName string) error {
	// カテゴリー名が提供されている場合、IDを検索します
	if categoryName != "" {
		category, err := s.challengerepo.FindCategoryByName(ctx, categoryName)
		if err != nil {
			return fmt.Errorf("failed to find category: %w", err)
		}
		if category == nil {
			return fmt.Errorf("category '%s' not found", categoryName)
		}
		challenge.CategoryID = &category.ID
	} else {
		// カテゴリー名がない場合はCategoryIDをnilに設定
		challenge.CategoryID = nil
	}

	// サービスはリポジトリのCreateメソッドを呼び出してデータベース操作を行います
	return s.challengerepo.Create(ctx, challenge)
}

func (s *challengeService) CollectByUsername(ctx context.Context, username string) ([]*models.Challenge, error) {
	userID, err := s.userrepo.GetIDByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}
	return s.challengerepo.CollectByUserID(ctx, userID)
}

func (s *challengeService) UpdateChallenge(ctx context.Context, challengeID uint, userID uint, req *dtos.UpdateChallengeRequest) error {
	challenge, err := s.challengerepo.GetByID(ctx, challengeID)
	if err != nil {
		return err
	}

	if challenge.UserID != userID {
		return errors.New("user is not the owner of the challenge")
	}

	if req.Title != nil {
		challenge.Title = *req.Title
	}
	if req.Description != nil {
		challenge.Description = *req.Description
	}
	if req.Score != nil {
		challenge.Score = *req.Score
	}
	if req.Flag != nil {
		challenge.Flag = *req.Flag
	}
	if req.IsPublic != nil {
		challenge.IsPublic = *req.IsPublic
	}

	if req.Category != nil {
		category, err := s.challengerepo.FindCategoryByName(ctx, *req.Category)
		if err != nil {
			return fmt.Errorf("failed to find category: %w", err)
		}
		if category == nil {
			return fmt.Errorf("category '%s' not found", *req.Category)
		}
		challenge.CategoryID = &category.ID
	} else {
		challenge.CategoryID = nil
	}

	return s.challengerepo.Update(ctx, challenge)
}

func (s *challengeService) DeleteChallenge(ctx context.Context, challengeID uint, userID uint) error {
	challenge, err := s.challengerepo.GetByID(ctx, challengeID)
	if err != nil {
		return err
	}

	if challenge.UserID != userID {
		return errors.New("user is not the owner of the challenge")
	}

	return s.challengerepo.Delete(ctx, challengeID)
}

func (s *challengeService) GetChallengeByID(ctx context.Context, challengeID uint, userID uint) (*dtos.ChallengeDetailResponse, error) {
	challenge, err := s.challengerepo.GetByID(ctx, challengeID)
	if err != nil {
		return nil, err
	}

	if challenge.UserID != userID {
		return nil, errors.New("user is not the owner of the challenge")
	}

	var categoryName *string
	if challenge.Category != nil {
		categoryName = &challenge.Category.Name
	}

	return &dtos.ChallengeDetailResponse{
		ID:          challenge.ID,
		Title:       challenge.Title,
		Description: challenge.Description,
		Category:    categoryName,
		Score:       challenge.Score,
		Flag:        challenge.Flag,
		IsPublic:    challenge.IsPublic,
	}, nil
}
