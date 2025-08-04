package repository

import (
	"context"

	"github.com/Saku0512/CTFLab/ctflab/internal/models"
	"gorm.io/gorm"
)

// ChallengeRepositoryは問題に関するDB操作インターフェース
type ChallengeRepository interface {
	Create(ctx context.Context, challenge *models.Challenge) error
	FindCategoryByName(ctx context.Context, name string) (*models.ChallengeCategory, error)
	CollectByUserID(ctx context.Context, userID uint) ([]*models.Challenge, error)
	GetByID(ctx context.Context, id uint) (*models.Challenge, error)
	Update(ctx context.Context, challenge *models.Challenge) error
	Delete(ctx context.Context, id uint) error
}

type challengeRepo struct {
	db *gorm.DB
}

// challengeRepoのコンストラクタ
func NewChallengeRepository(db *gorm.DB) ChallengeRepository {
	return &challengeRepo{db: db}
}

func (r *challengeRepo) Create(ctx context.Context, challenge *models.Challenge) error {
	return r.db.WithContext(ctx).Create(challenge).Error
}

// FindCategoryByNameは、カテゴリー名に基づいてChallengeCategoryを取得します。
func (r *challengeRepo) FindCategoryByName(ctx context.Context, name string) (*models.ChallengeCategory, error) {
	var category models.ChallengeCategory
	// GORMのWhereメソッドでカテゴリー名を検索
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // カテゴリーが見つからない場合はnilを返します
		}
		return nil, err
	}
	return &category, nil
}

// CollectByUserIDは、指定されたユーザーIDに関連付けられたチャレンジをすべて取得します。
func (r *challengeRepo) CollectByUserID(ctx context.Context, userID uint) ([]*models.Challenge, error) {
	var challenges []*models.Challenge
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&challenges).Error
	if err != nil {
		return nil, err
	}
	return challenges, nil
}

func (r *challengeRepo) GetByID(ctx context.Context, id uint) (*models.Challenge, error) {
	var challenge models.Challenge
	if err := r.db.WithContext(ctx).First(&challenge, id).Error; err != nil {
		return nil, err
	}
	return &challenge, nil
}

func (r *challengeRepo) Update(ctx context.Context, challenge *models.Challenge) error {
	return r.db.WithContext(ctx).Save(challenge).Error
}

func (r *challengeRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Challenge{}, id).Error
}
