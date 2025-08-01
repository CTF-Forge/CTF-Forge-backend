package service

import (
	"context"
	"fmt"

	"github.com/Saku0512/CTFLab/ctflab/internal/models"
	"github.com/Saku0512/CTFLab/ctflab/internal/repository"
)

type ChallengeService interface {
	CreateChallenge(ctx context.Context, challenge *models.Challenge, categoryName string) error
}

type challengeService struct {
	repo repository.ChallengeRepository
}

func NewChallengeService(repo repository.ChallengeRepository) ChallengeService {
	return &challengeService{repo: repo}
}

// CreateChallengeは、カテゴリー名を解決して新しい問題をデータベースに保存します。
func (s *challengeService) CreateChallenge(ctx context.Context, challenge *models.Challenge, categoryName string) error {
	// カテゴリー名が提供されている場合、IDを検索します
	if categoryName != "" {
		category, err := s.repo.FindCategoryByName(ctx, categoryName)
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
	return s.repo.Create(ctx, challenge)
}
