package repository

import (
	"context"
	"errors"
	"log"

	"github.com/CTF-Forge/CTF-Forge-backend/internal/models"
	"gorm.io/gorm"
)

// UserRepository はユーザーに関するDB操作インターフェースです。
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	UpdatePassword(ctx context.Context, userID uint, newHash string) error
	GetIDByUsername(ctx context.Context, username string) (uint, error)
}

// userRepo はUserRepositoryの実装です。
type userRepo struct {
	db *gorm.DB
}

// NewUserRepository はuserRepoのコンストラクタです。
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail は email からユーザーを取得します
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	// GORMのResultオブジェクトを受け取るパターンに修正
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)

	// エラーハンドリング
	if result.Error != nil {
		// レコードが見つからないエラーの場合、nilを返す
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("User with email '%s' not found.", email)
			return nil, nil
		}
		// その他のエラーの場合、エラーを返す
		log.Printf("Database error fetching user by email '%s': %v", email, result.Error)
		return nil, result.Error
	}

	// デバッグログ: ユーザーが見つかったことを確認
	log.Printf("Found user with email '%s'. Rows affected: %d", email, result.RowsAffected)

	// 成功した場合、ユーザーポインタとnilエラーを返す
	return &user, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, userID uint, newHash string) error {
	res := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password_hash", newHash)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// GetIDByUsernameはユーザー名からUserIDを取得する関数
func (r *userRepo) GetIDByUsername(ctx context.Context, username string) (uint, error) {
	var user models.User
	// ユーザー名で検索し、結果をuser変数に格納
	result := r.db.WithContext(ctx).Select("id").Where("username = ?", username).First(&user)

	if result.Error != nil {
		// レコードが見つからない場合もエラーとして扱う
		return 0, result.Error
	}

	// ユーザーIDを返す
	return user.ID, nil
}
