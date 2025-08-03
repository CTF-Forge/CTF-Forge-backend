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
	"os"
	"time"

	"github.com/Saku0512/CTFLab/ctflab/config"
	_ "github.com/Saku0512/CTFLab/ctflab/docs" // Swagger docs
	"github.com/Saku0512/CTFLab/ctflab/internal/router"
	"github.com/Saku0512/CTFLab/ctflab/oauth"
	"gorm.io/gorm/logger"
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
