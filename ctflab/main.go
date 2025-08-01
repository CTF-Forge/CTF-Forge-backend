package main

import (
	"log"

	"github.com/Saku0512/CTFLab/ctflab/config"
	"github.com/Saku0512/CTFLab/ctflab/internal/models"
	"github.com/Saku0512/CTFLab/ctflab/internal/router"
	"github.com/Saku0512/CTFLab/ctflab/oauth"
)

func main() {
	config.InitDB()
	db := config.GetDB()

	err := db.AutoMigrate(
		&models.User{},
		&models.OAuthAccount{}, // OAuthアカウントテーブルを追加
		&models.ChallengeCategory{},
		&models.Challenge{},
		&models.ChallengeFile{},
		&models.DockerChallenge{},
		&models.Submission{},
	)

	if err != nil {
		log.Fatal("migration failed:", err)
	}

	log.Println("Migration successful")

	oauth.Init()

	router := router.SetupRouter(db)

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
