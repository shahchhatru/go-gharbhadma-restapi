package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
	"myapp/code/db"
	"myapp/code/credentials"
	"myapp/code/user"		
	
	_ "myapp/docs" // Import your generated docs
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample Swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

func main() {
	app := fiber.New()

	db.InitialMigration()
	// db.DB.AutoMigrate(&credentials.Credentials{})

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault) // default
	// User routes
	app.Get("/users", user.GetUsers)
	app.Get("/users/:id", user.GetUser)
	app.Post("/users", user.SaveUser)
	app.Delete("/users/:id", user.DeleteUser)
	app.Put("/users/:id", user.UpdateUser)

	// Define your routes with Swagger comments
	app.Post("/login", credentials.Login)
	app.Get("/", accessible)
	app.Post("/signup", credentials.Signup)
	app.Post("/resetpasswordtoken", credentials.ResetPasswordRequestHandler)
	app.Post("/resetpassword", credentials.ResetPasswordConfirmationHandler)
	app.Post("/changepassword", credentials.ChangePasswordHandler)
	app.Post("/refreshtoken", credentials.RefreshToken)

	// JWT Middleware for protected routes
	// app.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	// }))

	// Restricted routes
	//app.Get("/restricted", restricted)

	app.Listen(":8080")
}


// @Summary Accessible route
// @Description Accessible route
// @Success 200 {string} string "Accessible"
// @Router / [get]
func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}
// @Summary Restricted route
// @Description Restricted route for authenticated users
// @Success 200 {string} string "Welcome {name}"
// @Router /restricted [get]
func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
