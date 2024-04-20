package main

import (
	
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/golang-jwt/jwt/v5"
	"myapp/code/db"
    "myapp/code/credentials"

	
)



func main() {
	app := fiber.New()

	db.InitialMigration()
    db.DB.AutoMigrate(&credentials.Credentials{})

	// Login route
	app.Post("/login", credentials.Login)

	// Unauthenticated route
	app.Get("/", accessible)

	app.Post("/signup",credentials.Signup)

    //reset pass word token
    app.Post("/resetpasswordtoken",credentials.ResetPasswordConfirmationHandler)
	
    //reset password
    app.Post("/resetpassword",credentials.ResetPasswordConfirmation)
    // change password
    app.Post("/changepassword",credentials.ChangePasswordHandler)
	
	app.Post("/refreshtoken",credentials.RefreshToken)

    // JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	// Restricted Routes
	app.Get("/restricted", restricted)

	app.Listen(":8080")
}









func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible ")
}


func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
