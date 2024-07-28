package user

import (
	"gorm.io/gorm"
	"myapp/code/db"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Signup route
// @Summary Signup
// @Description Signup user
// @Accept  json
// @Produce  json
// @Param   user  body  User  true  "User Info"
// @Success 201 {object} User
// @Failure 400 {object} fiber.Map
// @Router /signup [post]
func Signup(c *fiber.Ctx) error {
	user := new(User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// Additional logic for signup can go here
	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUsers route
// @Summary Get Users
// @Description Get all users
// @Produce  json
// @Success 200 {array} User
// @Router /users [get]
func GetUsers(c *fiber.Ctx) error {
	var users []User
	db.DB.Find(&users)
	return c.JSON(&users)
}

// GetUser route
// @Summary Get User
// @Description Get a single user by ID
// @Param   id  path  int  true  "User ID"
// @Produce  json
// @Success 200 {object} User
// @Failure 404 {object} fiber.Map
// @Router /users/{id} [get]
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	result := db.DB.Find(&user, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(&user)
}

// SaveUser route
// @Summary Save User
// @Description Save a new user
// @Accept  json
// @Produce  json
// @Param   user  body  User  true  "User Info"
// @Success 201 {object} User
// @Failure 500 {object} fiber.Map
// @Router /users [post]
func SaveUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	db.DB.Create(&user)
	return c.Status(fiber.StatusCreated).JSON(&user)
}

// DeleteUser route
// @Summary Delete User
// @Description Delete a user by ID
// @Param   id  path  int  true  "User ID"
// @Success 200 {string} string "User is deleted!!!"
// @Failure 500 {object} fiber.Map
// @Router /users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	db.DB.First(&user, id)
	if user.Email == "" {
		return c.Status(500).SendString("User not available")
	}

	db.DB.Delete(&user)
	return c.SendString("User is deleted!!!")
}

// UpdateUser route
// @Summary Update User
// @Description Update a user's information
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "User ID"
// @Param   user  body  User  true  "User Info"
// @Success 200 {object} User
// @Failure 500 {object} fiber.Map
// @Router /users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(User)
	db.DB.First(&user, id)
	if user.Email == "" {
		return c.Status(500).SendString("User not available")
	}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	db.DB.Save(&user)
	return c.JSON(&user)
}
