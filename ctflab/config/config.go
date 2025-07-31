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

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

func GetJWTIssuer() string {
	return os.Getenv("JWT_ISSUER")
}

func GetJWTExpireDuration() time.Duration {
	hStr := os.Getenv("JWT_EXPIRE_HOURS")
	h, err := strconv.Atoi(hStr)
	if err != nil || h <= 0 {
		h = 24 // デフォルト24時間
	}
	return time.Duration(h) * time.Hour
}
