package config

import (
	"log"
	"math/rand"
	"mnc/msg"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(config Config) *gorm.DB {
	username, _ := os.LookupEnv("MNC_PAY_DB_USERNAME")
	password, _ := os.LookupEnv("MNC_PAY_DB_PASSWORD")
	host, _ := os.LookupEnv("MNC_PAY_DB_HOST")
	port, _ := os.LookupEnv("MNC_PAY_DB_PORT")
	dbName, _ := os.LookupEnv("MNC_PAY_DB_NAME")
	maxcon, _ := os.LookupEnv("MNC_PAY_POOL_MAX_CONN")
	idlecon, _ := os.LookupEnv("MNC_PAY_POOL_IDLE_CONN")
	lifetime, _ := os.LookupEnv("MNC_PAY_POOL_LIFE_TIME")

	maxPoolOpen, _ := strconv.Atoi(maxcon)
	maxPoolIdle, _ := strconv.Atoi(idlecon)
	maxPollLifeTime, err := strconv.Atoi(lifetime)
	msg.PanicLogging(err)

	loggerDb := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(username+":"+password+"@tcp("+host+":"+port+")/"+dbName+"?parseTime=true"), &gorm.Config{
		Logger: loggerDb,
	})
	msg.PanicLogging(err)

	sqlDB, err := db.DB()
	msg.PanicLogging(err)

	sqlDB.SetMaxOpenConns(maxPoolOpen)
	sqlDB.SetMaxIdleConns(maxPoolIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(rand.Int31n(int32(maxPollLifeTime))) * time.Millisecond)

	return db
}
