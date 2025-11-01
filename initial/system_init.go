package initial

import (
	"airport-match-system/models"
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var RDB *redis.Client
var MysqlDB *gorm.DB

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file, %s", err)
		panic("fatal error config file")
	}
}

func InitRedis() *redis.Client {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("Redis connection failed: %v", err)
		panic("redis connection failed")
	}

	log.Printf("Connected to Redis")
	return RDB
}

func InitMysqlDB() *gorm.DB {
	var err error
	MysqlDB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic("failed to connect database")
	}
	err = MysqlDB.AutoMigrate(&models.DriverPublishOrder{}, &models.PassengerPublishOrder{}, &models.MatchRelation{}, &models.MatchOrderResult{}, &models.FeeConfig{})
	if err != nil {
		log.Fatal(err)
		panic("failed to migrate database")
	}

	log.Printf("Connected to Mysql")
	return MysqlDB
}
