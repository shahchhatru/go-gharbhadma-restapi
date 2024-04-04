package user

import(
	"gorm.io/gorm"
	"github.com/gofiber/fiber/v2"
	"myapp/code/db"
)

type User struct {
	gorm.Model
	FirstName string  `json:"firstname"`
	LastName string   `json:"lastname"`
	Email string 	   `json:"email"`

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