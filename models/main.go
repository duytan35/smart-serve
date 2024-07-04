package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

// var newLogger = logger.New(
// 	log.New(os.Stdout, "\n", log.LstdFlags), // io writer
// 	logger.Config{
// 		LogLevel: logger.Info, // Log level
// 		Colorful: true,        // Enable color
// 	},
// )

func ConnectDB() {
	dsn := os.Getenv("DB_DSN")

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: newLogger,
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		fmt.Println("Connected to database")
	}
}
