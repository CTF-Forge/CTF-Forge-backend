package service

import (
	"context"
	"errors"
	"time"

	"github.com/Saku0512/CTFLab/ctflab/internal/models"
	"github.com/Saku0512/CTFLab/ctflab/internal/repository"
	"github.com/Saku0512/CTFLab/ctflab/pkg/token"
	"gorm.io/gorm"
)

type OAuthService struct {
	oauthRepo  repository.OAuthAccountRepository
	userRepo   repository.UserRepository
	jwtManager *token.JWTManager
}

func NewOAuthService(oauthRepo repository.OAuthAccountRepository, userRepo repository.UserRepository, jwtManager *token.JWTManager) *OAuthService {
	return &OAuthService{
		oauthRepo:  oauthRepo,
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *OAuthService) HandleOAuthCallback(
	ctx context.Context,
	username string,
	provider string,
	providerUserID string,
	email string,
	accessToken string,
	refreshToken string,
	tokenExpiry time.Time,
) (*token.TokenPair, error) {
	// OAuthアカウントが存在するか確認
	account, err := s.oauthRepo.FindByProviderAndProviderUserID(provider, providerUserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var user *models.User
	if account != nil {
		user, err = s.userRepo.GetByID(ctx, account.UserID)
		if err != nil {
			return nil, err
		}

		// トークン更新
		if err := s.oauthRepo.UpdateTokenInfo(account.ID, accessToken, refreshToken, tokenExpiry); err != nil {
			return nil, err
		}
	} else {
		// ユーザー作成
		// 同じ名前のユーザー名のユーザーがいるかどうかチェック
		exitingUser, err := s.userRepo.GetByUsername(ctx, username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if exitingUser != nil {
			// 既存のユーザーが見つかった場合、そのユーザーを使用
			user = exitingUser
		} else {
			// ユーザー名が重複シニア場合、新しいユーザーを作成
			user = &models.User{
				Email:    email,
				Username: username,
			}
			if err := s.userRepo.Create(ctx, user); err != nil {
				return nil, err
			}
		}

		// OAuthAccount作成
		newAccount := &models.OAuthAccount{
			Provider:       provider,
			ProviderUserID: providerUserID,
			AccessToken:    accessToken,
			RefreshToken:   refreshToken,
			TokenExpiry:    tokenExpiry,
			UserID:         user.ID,
		}

		if err := s.oauthRepo.Create(newAccount); err != nil {
			return nil, err
		}
	}

	// JWTトークンペアを生成
	return s.jwtManager.GenerateTokenPair(user.ID, user.Username)
}

// GetUserByOAuthAccount OAuthアカウントからユーザー情報を取得
func (s *OAuthService) GetUserByOAuthAccount(provider, providerUserID string) (*models.User, error) {
	account, err := s.oauthRepo.FindByProviderAndProviderUserID(provider, providerUserID)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("oauth account not found")
	}

	return s.userRepo.GetByID(context.Background(), account.UserID)
}
