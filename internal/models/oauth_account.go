package models

import "time"

type OAuthAccount struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"not null" json:"user_id"`
	User           User      `gorm:"foreignKey:UserID" json:"-"`
	Provider       string    `gorm:"not null;uniqueIndex:idx_provider_user" json:"provider"`         // 例: github, google
	ProviderUserID string    `gorm:"not null;uniqueIndex:idx_provider_user" json:"provider_user_id"` // プロバイダー側のユーザーID
	AccessToken    string    `json:"access_token,omitempty"`                                         // 必要なら保存
	RefreshToken   string    `json:"refresh_token,omitempty"`
	TokenExpiry    time.Time `json:"token_expiry,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	// 一意制約: provider と provider_user_id の組み合わせ
	// UNIQUE(provider, provider_user_id) は GORM v2 で以下のように書ける
	// -> gorm:"uniqueIndex:idx_provider_user"
}

func (OAuthAccount) TableName() string {
	return "oauth_accounts"
}
