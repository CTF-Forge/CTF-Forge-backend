package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/Saku0512/CTFLab/ctflab/internal/models"
	"github.com/Saku0512/CTFLab/ctflab/internal/repository"
	"github.com/Saku0512/CTFLab/ctflab/pkg/token"
)

// AuthService は認証に関わるビジネスロジックを提供します。
type AuthService struct {
	userRepo   repository.UserRepository
	jwtManager *token.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *token.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// HashPassword は平文パスワードをbcryptでハッシュ化します。
func (s *AuthService) HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

// VerifyPassword は平文パスワードとハッシュが一致するか検証します。
func (s *AuthService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Login はメールアドレスとパスワードを検証し、成功すればトークンペアを返します。
func (s *AuthService) Login(ctx context.Context, email, password string) (*token.TokenPair, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if err := s.VerifyPassword(user.PasswordHash, password); err != nil {
		return nil, errors.New("invalid password")
	}

	return s.jwtManager.GenerateTokenPair(user.ID, user.Username)
}

// RegisterUser は新しいユーザーを登録します。
func (s *AuthService) RegisterUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Create(ctx, user)
}

// RefreshToken はリフレッシュトークンを使用して新しいトークンペアを生成します。
func (s *AuthService) RefreshToken(refreshToken string) (*token.TokenPair, error) {
	return s.jwtManager.RefreshTokenPair(refreshToken)
}

// ValidateToken はアクセストークンを検証します。
func (s *AuthService) ValidateToken(tokenStr string) (*token.UserClaims, error) {
	return s.jwtManager.VerifyAccessToken(tokenStr)
}

// IsTokenExpired はトークンが期限切れかどうかを判定します。
func (s *AuthService) IsTokenExpired(tokenStr string) bool {
	return s.jwtManager.IsTokenExpired(tokenStr)
}
