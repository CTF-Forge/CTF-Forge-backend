// Package main CTFForge API Server
//
// CTFForgeは、誰もがCTFの問題を作成し、公開できるプラットフォームです。
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
	"os"
	"time"

	"github.com/CTF-Forge/CTF-Forge-backend/config"
	_ "github.com/CTF-Forge/CTF-Forge-backend/docs" // Swagger docs
	"github.com/CTF-Forge/CTF-Forge-backend/internal/router"
	"github.com/CTF-Forge/CTF-Forge-backend/oauth"
	"gorm.io/gorm/logger"
)

// @title CTFForge API
// @version 1.0
// @description CTFForgeは、誰もがCTFの問題を作成し、公開できるプラットフォームです。
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

	// main.go の db := config.GetDB() の後に以下のコードを追加
	db.Logger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	// これで、実行されるすべてのSQLクエリがログに出力されます

	oauth.Init()

	router := router.SetupRouter(db)

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("failed to run server:", err)
	}
}
