package repository

import (
	"errors"
	"time"

	"github.com/Saku0512/CTFLab/ctflab/internal/models"
	"gorm.io/gorm"
)

type OAuthAccountRepository interface {
	FindByProviderAndProviderUserID(provider string, providerUserID string) (*models.OAuthAccount, error)
	Create(account *models.OAuthAccount) error
	FindOrCreate(account *models.OAuthAccount) (*models.OAuthAccount, error)
	UpdateTokenInfo(accountID uint, accessToken, refreshToken string, tokenExpiry time.Time) error
}

type oauthRepo struct {
	db *gorm.DB
}

func NewOAuthAccountRepository(db *gorm.DB) OAuthAccountRepository {
	return &oauthRepo{db: db}
}

func (r *oauthRepo) FindByProviderAndProviderUserID(provider string, providerUserID string) (*models.OAuthAccount, error) {
	var account models.OAuthAccount
	if err := r.db.Where("provider = ? AND provider_user_id = ?", provider, providerUserID).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *oauthRepo) Create(account *models.OAuthAccount) error {
	return r.db.Create(account).Error
}

func (r *oauthRepo) FindOrCreate(account *models.OAuthAccount) (*models.OAuthAccount, error) {
	var existing models.OAuthAccount
	err := r.db.
		Where("provider = ? AND provider_user_id = ?", account.Provider, account.ProviderUserID).
		First(&existing).Error

	if err == nil {
		return &existing, nil // 見つかったのでそれを返す
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := r.db.Create(account).Error; err != nil {
			return nil, err
		}
		return account, nil
	}
	return nil, err
}

func (r *oauthRepo) UpdateTokenInfo(accountID uint, accessToken, refreshToken string, tokenExpiry time.Time) error {
	return r.db.Model(&models.OAuthAccount{}).
		Where("id = ?", accountID).
		Updates(map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_expiry":  tokenExpiry,
		}).Error
}
