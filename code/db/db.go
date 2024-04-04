package db

import (
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
	"fmt"
)


var DB *gorm.DB
var err error

const DNS = "root:123@tcp(dev-db:3306)/nest?charset=utf8mb4&parseTime=True&loc=Local"

func InitialMigration() {
	DB, err = gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to Database")
	}

	// // Ping the database to check connectivity
	// sqlDB, err := DB.DB()
	// if err != nil {
	// 	log.Fatalf("Error getting underlying SQL database: %v", err)
	// }

	
}
    // Open a connection to the database




		