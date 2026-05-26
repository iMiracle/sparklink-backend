package config

import (
<<<<<<< HEAD
=======
	"database/sql"
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	"fmt"
	"log"
	"os"
	"time"

<<<<<<< HEAD
=======
	_ "github.com/go-sql-driver/mysql"
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
<<<<<<< HEAD
	Port         string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	RedisHost    string
	RedisPort    string
	JWTSecret    string
	JWTExpire    time.Duration
	SMSAPIKey    string
	SMSAppSecret string
	SMSEndpoint  string
=======
	// Server
	Port string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis
	RedisHost string
	RedisPort string

	// JWT
	JWTSecret string
	JWTExpire time.Duration

	// SMS (for verification code)
	SMSAPIKey    string
	SMSAppSecret string
	SMSEndpoint  string

	// App
	AppDownloadURL string
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "3306"),
		DBUser:       getEnv("DB_USER", "root"),
		DBPassword:   getEnv("DB_PASSWORD", ""),
		DBName:       getEnv("DB_NAME", "sparklink"),
		RedisHost:    getEnv("REDIS_HOST", "localhost"),
		RedisPort:    getEnv("REDIS_PORT", "6379"),
		JWTSecret:    getEnv("JWT_SECRET", "sparklink-secret-key"),
<<<<<<< HEAD
		JWTExpire:    30 * time.Minute,
		SMSAPIKey:    getEnv("SMS_API_KEY", ""),
		SMSAppSecret: getEnv("SMS_APP_SECRET", ""),
		SMSEndpoint:  getEnv("SMS_ENDPOINT", ""),
=======
		JWTExpire:    7 * 24 * time.Hour,
		SMSAPIKey:    getEnv("SMS_API_KEY", ""),
		SMSAppSecret: getEnv("SMS_APP_SECRET", ""),
		SMSEndpoint: getEnv("SMS_ENDPOINT", ""),
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	}
}

func InitDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")
	return db, nil
}

<<<<<<< HEAD
func InitRedis(cfg *Config) (*redis.Client, error) {
=======
func InitRedis(cfg *Config) *redis.Client {
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: "",
		DB:       0,
	})

<<<<<<< HEAD
	if err := rdb.Ping(rdb.Context()).Err(); err != nil {
		return nil, err
	}

	log.Println("Redis connected successfully")
	return rdb, nil
=======
	log.Println("Redis connected successfully")
	return rdb
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
<<<<<<< HEAD
}
=======
}
>>>>>>> 4444d9abefbcf39a2473e97f16b5ac708632885f
