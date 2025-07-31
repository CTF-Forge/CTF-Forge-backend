package models

import "time"

type DockerChallenge struct {
	ID          uint      `gorm:"primaryKey"`
	ChallengeID uint      `gorm:"not null"`
	Challenge   Challenge `gorm:"foreignKey:ChallengeID"`
	ImageTag    string    `gorm:"not null"`
	ExposedPort int       `gorm:"not null"`
	Entrypoint  string    `gorm:"not null"`
	CreatedAt   time.Time
}
