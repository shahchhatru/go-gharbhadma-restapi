package user

import(
	"gorm.io/gorm"
	"myapp/code/db"
	"crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
	"github.com/dgrijalva/jwt-go"
	jwtware "github.com/gofiber/contrib/jwt"
    "github.com/gofiber/fiber/v2"
    "golang.org/x/crypto/bcrypt"

)

type User struct {
	gorm.Model
	FirstName string  `json:"firstname"`
	LastName string   `json:"lastname"`
	Email string 	   `json:"email"`


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

// const (
//     jwtSecret,  err := generateRandomString(length) // Change this to a strong secret key in production
//     jwtExpiration = time.Hour * 24    // Token expiration time (1 day)
// )


func Signup(c *fiber.Ctx) error{
	user := new(User)

	if err :=c.BodyParser(user); err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}
}


func GetUsers(c *fiber.Ctx) error {
	var users []User
	db.DB.Find(&users)
	return c.JSON(&users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	db.DB.Find(&user, id)
	return c.JSON(&user)
}

func SaveUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	db.DB.Create(&user)
	return c.JSON(&user)
}

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