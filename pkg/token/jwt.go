package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrInvalidClaims    = errors.New("invalid token claims")
)

// JWTManager JWTトークンの生成・検証を管理
type JWTManager struct {
	accessSecretKey  string
	refreshSecretKey string
	accessDuration   time.Duration
	refreshDuration  time.Duration
	issuer           string
}

// UserClaims JWTのカスタムクレーム
type UserClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// TokenPair アクセストークンとリフレッシュトークンのペア
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // 秒単位
}

// NewJWTManager 新しいJWTマネージャーを作成
func NewJWTManager(accessSecretKey, refreshSecretKey, issuer string, accessDuration, refreshDuration time.Duration) *JWTManager {
	return &JWTManager{
		accessSecretKey:  accessSecretKey,
		refreshSecretKey: refreshSecretKey,
		accessDuration:   accessDuration,
		refreshDuration:  refreshDuration,
		issuer:           issuer,
	}
}

// GenerateTokenPair ユーザー情報からアクセストークンとリフレッシュトークンを生成
func (j *JWTManager) GenerateTokenPair(userID uint, username string) (*TokenPair, error) {
	now := time.Now()

	// アクセストークン生成
	accessClaims := UserClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessDuration)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(j.accessSecretKey))
	if err != nil {
		return nil, err
	}

	// リフレッシュトークン生成
	refreshClaims := UserClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshDuration)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(j.refreshSecretKey))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(j.accessDuration.Seconds()),
	}, nil
}

// GenerateAccessToken アクセストークンのみを生成
func (j *JWTManager) GenerateAccessToken(userID uint, username string) (string, error) {
	now := time.Now()
	claims := UserClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.accessSecretKey))
}

// VerifyAccessToken アクセストークンを検証
func (j *JWTManager) VerifyAccessToken(tokenStr string) (*UserClaims, error) {
	return j.verifyToken(tokenStr, j.accessSecretKey)
}

// VerifyRefreshToken リフレッシュトークンを検証
func (j *JWTManager) VerifyRefreshToken(tokenStr string) (*UserClaims, error) {
	return j.verifyToken(tokenStr, j.refreshSecretKey)
}

// verifyToken トークンの検証（内部メソッド）
func (j *JWTManager) verifyToken(tokenStr, secretKey string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// HS256署名か確認
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidSignature
			}
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		// 詳細なエラーハンドリング
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// RefreshTokenPair リフレッシュトークンを使用して新しいトークンペアを生成
func (j *JWTManager) RefreshTokenPair(refreshToken string) (*TokenPair, error) {
	claims, err := j.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return j.GenerateTokenPair(claims.UserID, claims.Username)
}

// IsTokenExpired トークンが期限切れかどうかを判定
func (j *JWTManager) IsTokenExpired(tokenStr string) bool {
	_, err := j.VerifyAccessToken(tokenStr)
	return errors.Is(err, ErrExpiredToken)
}
