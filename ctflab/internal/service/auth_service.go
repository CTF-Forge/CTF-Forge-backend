package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Saku0512/CTFLab/ctflab/internal/models"
	"github.com/Saku0512/CTFLab/ctflab/internal/repository"
)

// AuthService は認証に関わるビジネスロジックを提供します。
type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret []byte
	jwtIssuer string
	jwtExpiry time.Duration
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret []byte, jwtIssuer string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtIssuer: jwtIssuer,
		jwtExpiry: jwtExpiry,
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

// Claims はJWTのカスタムクレーム構造体です。
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken はユーザー情報を元にJWTを生成します。
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtExpiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ParseToken はJWT文字列を解析してClaimsを返します。
func (s *AuthService) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// Login はメールアドレスとパスワードを検証し、成功すればJWTを返します。
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	if err := s.VerifyPassword(user.PasswordHash, password); err != nil {
		return "", errors.New("invalid password")
	}

	return s.GenerateToken(user)
}

func (s *AuthService) RegisterUser(ctx context.Context, user *models.User) error {
	return s.userRepo.Create(ctx, user)
}
