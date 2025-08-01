package models

import "time"

type Submission struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID"`
	ChallengeID uint      `gorm:"not null"`
	Challenge   Challenge `gorm:"foreignKey:ChallengeID"`
	SubmittedAt time.Time
	Flag        string `gorm:"not null"`
	IsCorrect   bool   `gorm:"default:false"`
}
