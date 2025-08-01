package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func InitDB() {
	LoadEnv()

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
}

func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("database not initialized")
	}
	return DB
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}
}

func GetGitHubAuthKey() string {
	return os.Getenv("GITHUB_KEY")
}

func GetGitHubAuthSecret() string {
	return os.Getenv("GITHUB_SECRET")
}

func GetGitHubCallbackURL() string {
	return os.Getenv("GITHUB_CALLBACK")
}

func GetGoogleAuthKey() string {
	return os.Getenv("GOOGLE_KEY")
}

func GetGoogleAuthSecret() string {
	return os.Getenv("GOOGLE_SECRET")
}

func GetGoogleCallbackURL() string {
	return os.Getenv("GOOGLE_CALLBACK")
}

// JWT設定
func GetJWTAccessSecret() string {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		secret = os.Getenv("JWT_SECRET") // 後方互換性
	}
	if secret == "" {
		log.Fatal("JWT_ACCESS_SECRET or JWT_SECRET environment variable is required")
	}
	return secret
}

func GetJWTRefreshSecret() string {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = GetJWTAccessSecret() // デフォルトはアクセスシークレットと同じ
	}
	return secret
}

func GetJWTIssuer() string {
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "ctflab"
	}
	return issuer
}

func GetJWTAccessExpireDuration() time.Duration {
	hStr := os.Getenv("JWT_ACCESS_EXPIRE_HOURS")
	h, err := strconv.Atoi(hStr)
	if err != nil || h <= 0 {
		h = 1 // デフォルト1時間
	}
	return time.Duration(h) * time.Hour
}

func GetJWTRefreshExpireDuration() time.Duration {
	hStr := os.Getenv("JWT_REFRESH_EXPIRE_HOURS")
	h, err := strconv.Atoi(hStr)
	if err != nil || h <= 0 {
		h = 168 // デフォルト7日間（24 * 7）
	}
	return time.Duration(h) * time.Hour
}

// 後方互換性のための関数
func GetJWTSecret() string {
	return GetJWTAccessSecret()
}

func GetJWTExpireDuration() time.Duration {
	return GetJWTAccessExpireDuration()
}

func GetSessionSecret() string {
	return os.Getenv("SESSION_SECRET")
}
