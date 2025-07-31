package models

import "time"

type Challenge struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null"`
	User        User   `gorm:"foreignKey:UserID"`
	Title       string `gorm:"not null"`
	Description string
	CategoryID  *uint
	Category    *ChallengeCategory `gorm:"foreignKey:CategoryID"`
	Flag        string
	IsPublic    bool `gorm:"default:false"`
	Score       int  `gorm:"default:0"`
	CreatedAt   time.Time
}
