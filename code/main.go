package main

import (
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
    "log"
    "github.com/gofiber/fiber/v2"
    "myapp/code/user"
)

func home(c *fiber.Ctx) error{
    return c.SendString("This is the home route")

}

func main() {
    app:=fiber.New()

    app.Get("/",home)
    // MySQL connection string
    dsn := "root:123@tcp(dev-db:3306)/nest?charset=utf8mb4&parseTime=True&loc=Local"

    // Open a connection to the database
    
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    

    // Ping the database to check connectivity
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatalf("Error getting underlying SQL database: %v", err)
    }

    db.AutoMigrate(&user.User{})
    defer sqlDB.Close()

    // Database connection successful
    log.Println("Connected to MySQL database")

    // You can now use the 'db' object to interact with the database
    // For example, you can define your models and perform CRUD operations

    app.Listen(":8080")
}
