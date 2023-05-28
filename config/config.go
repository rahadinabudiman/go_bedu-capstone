package config

import (
	"fmt"
	"go_bedu/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Host     string
	DB_Name     string
}

func ConnectDB() (*gorm.DB, error) {
	// Load the Asia/Jakarta location
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Handle the error
	}

	config := Config{
		DB_Username: "r4ha",
		DB_Password: "kmoonkinan",
		DB_Port:     "3306",
		DB_Host:     "localhost",
		DB_Name:     "go_bedu",
	}

	ConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB_Username,
		config.DB_Password,
		config.DB_Host,
		config.DB_Port,
		config.DB_Name,
	)

	dbConn, err := gorm.Open(mysql.Open(ConnectionString), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	dbConn = dbConn.Session(&gorm.Session{
		NowFunc: func() time.Time {
			return time.Now().In(location)
		},
	})

	return dbConn, nil
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Administrator{},
		&models.Article{},
	)
}
