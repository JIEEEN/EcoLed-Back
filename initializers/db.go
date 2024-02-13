package initializers

import (
	"fmt"
	"os"

	"github.com/Eco-Led/EcoLed-Back_test/logger"
	"github.com/Eco-Led/EcoLed-Back_test/models"

	_redis "github.com/go-redis/redis/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Redis *_redis.Client

func InitDB() {
	ConnectDB()
	SyncDB()
	InitRedis(1)
}

func InitRedis(selectDB ...int) {
	var redisHost = os.Getenv("REDIS_HOST")
	var redisPassword = os.Getenv("REDIS_PASSWORD")

	Redis = _redis.NewClient(&_redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       selectDB[0],
	})
}

func ConnectDB() {
	logger.Info.Println("Connect to Database")

	var err error

	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Error.Println("failed to connect to DB")
		os.Exit(-1)
	}

	dbName := os.Getenv("DB_NAME")
	res := DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))

	if res.Error != nil {
		logger.Error.Println("failed to create database")
		os.Exit(-1)
	}
}

func SyncDB() {
	logger.Debug.Println("Synchronization to Database")

	DB.AutoMigrate(&models.Users{})
	DB.AutoMigrate(&models.OAuth{})
	DB.AutoMigrate(&models.Profiles{})
	DB.AutoMigrate(&models.Accounts{})
	DB.AutoMigrate(&models.Paylogs{})
	DB.AutoMigrate(&models.Posts{})
}
