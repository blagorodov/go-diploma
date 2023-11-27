package drivers

import (
	"diploma/internal/config"
	"diploma/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase() Database {
	db, err := gorm.Open(postgres.Open(config.Options.DatabaseURI))
	if err != nil {
		logger.Log("Failed to connect to database")
		logger.Log(err.Error())
		panic(err)
	}
	logger.Log("Database connection established")
	return Database{
		DB: db,
	}
}
