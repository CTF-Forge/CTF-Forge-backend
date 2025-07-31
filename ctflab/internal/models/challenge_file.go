package models

import "time"

type ChallengeFile struct {
	ID          uint      `gorm:"primaryKey"`
	ChallengeID uint      `gorm:"not null"`
	Challenge   Challenge `gorm:"foreignKey:ChallengeID"`
	Filename    string    `gorm:"not null"`
	Filepath    string    `gorm:"not null"`
	Mimetype    string    `gorm:"not null;check:mimetype = 'application/zip'"`
	Size        int       `gorm:"not null"`
	UploadedAt  time.Time
}
