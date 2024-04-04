package main

import (
    "github.com/gofiber/fiber/v2"
    "myapp/code/user"
    "myapp/code/db"
)

func home(c *fiber.Ctx) error{
    return c.SendString("This is the home route")

}

func Routers(app *fiber.App){
    app.Get("/users",user.GetUsers)
    app.Get("/user/:id",user.GetUser)
    app.Post("/user",user.SaveUser)
    app.Delete("/user/:id",user.DeleteUser)
    app.Put("/user/:id",user.UpdateUser)




}

func main() {
   

    db.InitialMigration()
    db.DB.AutoMigrate(&user.User{})
    // You can now use the 'db' object to interact with the database
    // For example, you can define your models and perform CRUD operations
    app:=fiber.New()

    app.Get("/",home)
    Routers(app)


    app.Listen(":8080")
}
