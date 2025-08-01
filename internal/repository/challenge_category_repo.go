package repository

import (
	"log"

	"github.com/Saku0512/CTFLab/ctflab/internal/models"
	"gorm.io/gorm"
)

// ChallengeCategoryRepositoryはカテゴリーを初期化するDB操作インターフェース
type ChallengeCategoryRepository interface {
	InitCategories(db *gorm.DB) error
}

var initialCategories = []string{
	"Crypto",
	"Reversing",
	"Web",
	"PWN",
	"PPC",
	"OSINT",
	"Misc",
}

func InitCategories(db *gorm.DB) error {
	for _, name := range initialCategories {
		var count int64
		if err := db.Model(&models.ChallengeCategory{}).Where("name = ?", name).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := db.Create(&models.ChallengeCategory{Name: name}).Error; err != nil {
				return err
			}
			log.Printf("Category '%s' added.", name)
		} else {
			log.Printf("Category '%s' already exists. Skipping.", name)
		}
	}
	return nil
}
