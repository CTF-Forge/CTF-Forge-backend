// Package main CTFLab API Server
//
// CTFLabは、誰もがCTFの問題を作成し、公開できるプラットフォームです。
//
//	Schemes: http, https
//	Host: localhost:8080
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- bearer
//
// swagger:meta
package main

import (
	"log"

	"github.com/Saku0512/CTFLab/ctflab/config"
	_ "github.com/Saku0512/CTFLab/ctflab/docs" // Swagger docs
	"github.com/Saku0512/CTFLab/ctflab/internal/models"
	"github.com/Saku0512/CTFLab/ctflab/internal/repository"
	"github.com/Saku0512/CTFLab/ctflab/internal/router"
	"github.com/Saku0512/CTFLab/ctflab/oauth"
)

// @title CTFLab API
// @version 1.0
// @description CTFLabは、誰もがCTFの問題を作成し、公開できるプラットフォームです。
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey bearer
// @in header
// @name Authorization
// @description Bearer token for authentication

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

	if err := repository.InitCategories(db); err != nil {
		log.Fatalf("カテゴリー初期化エラー: %v", err)
	}

	oauth.Init()

	router := router.SetupRouter(db)

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
