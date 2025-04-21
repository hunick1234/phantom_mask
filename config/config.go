package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBname   string
	Host     string
	Port     string
	User     string
	PassWord string
}

type AppConfig struct {
	DB        DBConfig
	Port      string
	JWTSecret string
}

func NewTestDBConfig() *DBConfig {
	return &DBConfig{
		DBname:   "testdb",
		Host:     "localhost",
		Port:     "5435",
		User:     "user",
		PassWord: "pass",
	}
}

func (dbConfig *DBConfig) ToDSN() string {
	return "host=" + dbConfig.Host +
		" user=" + dbConfig.User +
		" password=" + dbConfig.PassWord +
		" dbname=" + dbConfig.DBname +
		" port=" + dbConfig.Port +
		" sslmode=disable"
}

// 呼叫這個方法載入設定
func LoadConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, fallback to system environment variables")
	}

	return &AppConfig{
		DB: DBConfig{
			DBname:   getEnv("DB_NAME", "phantom_mask"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5433"),
			User:     getEnv("DB_USER", "postgres"),
			PassWord: getEnv("DB_PASSWORD", "postgres"),
		},
		Port:      getEnv("PORT", "8081"),
		JWTSecret: getEnv("JWT_SECRET", "keyForJWT"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
