package credentials

import(
	"time"
	"github.com/gofiber/fiber/v2"
	"fmt"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"myapp/code/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"encoding/json"
)

type Credentials struct {
    ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Email string `gorm:"type:varchar(255);uniqueIndex" json:"email"`
    Pass  string `json:"pass"`
}



type ErrorResponse struct {
    Message string `json:"message"`
}

func hashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func Signup(c *fiber.Ctx) error {
    var user Credentials

    // Parse the JSON payload from the request body
    if err := c.BodyParser(&user); err != nil {
        return c.SendStatus(fiber.StatusBadRequest)
    }

    // Hash the user's password
    hashedPassword, err := hashPassword(user.Pass)
    if err != nil {
        return c.SendStatus(fiber.StatusInternalServerError)
    }

    // Print the received user details
    fmt.Printf("New user signed up: Username: %s\n", user.Email)

    // Save the user in the database with the hashed password
    hashedUser := Credentials{
        Email: user.Email,
        Pass: hashedPassword,
    }
    if err := db.DB.Create(&hashedUser).Error; err != nil {
        return c.SendStatus(fiber.StatusInternalServerError)
    }
	accessToken, err := generateAccessToken(user.Email)
    if err != nil {
        return c.SendStatus(fiber.StatusInternalServerError)
    }

    // Generate refresh token
    refreshToken, err := generateRefreshToken(user.Email)
    if err != nil {
        return c.SendStatus(fiber.StatusInternalServerError)
    }


    return c.JSON(fiber.Map{"message": "User signed up successfully","token":fiber.Map{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    }})
}

func Login(c *fiber.Ctx) error {
	// Parse the JSON payload from the request body
	var creds Credentials 
	if err := c.BodyParser(&creds); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// Log the values of user and pass
	fmt.Printf("User: %s, Pass: %s\n", creds.Email, creds.Pass)

	// Query the database to find the user by username
	var user Credentials
	if err := db.DB.Where("email = ?", creds.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User not found
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		// Other database errors
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Compare the hashed password from the database with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(creds.Pass)); err != nil {
		// Password does not match
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Password matches, create JWT token
	// claims := jwt.MapClaims{
	// 	"name":  user.User,
	// 	"admin": false, // Example: You can set admin status based on user roles from database
	// 	"exp":   time.Now().Add(time.Hour * 72).Unix(),
	// }

	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response
	// t, err := token.SignedString([]byte("secret"))
	// if err != nil {
	// 	return c.SendStatus(fiber.StatusInternalServerError)
	// }

	// Generate access token
    accessToken, err := generateAccessToken(user.Email)
    if err != nil {
        return c.SendStatus(fiber.StatusInternalServerError)
    }

    // Generate refresh token
    refreshToken, err := generateRefreshToken(user.Email)
    if err != nil {
        return c.SendStatus(fiber.StatusInternalServerError)
    }

    // Send tokens to client
    return c.JSON(fiber.Map{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}


func generateAccessToken(username string) (string, error) {
    // Create access token claims
    claims := jwt.MapClaims{
        "sub": username, // Use username as the subject
        "exp": time.Now().Add(time.Minute * 15).Unix(), // Access token expires in 15 minutes
        // Add other claims...
    }

    // Generate access token
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return accessToken.SignedString([]byte("access_secret"))
}

func generateRefreshToken(username string) (string, error) {
    // Create refresh token claims
    claims := jwt.MapClaims{
        "sub": username, // Use username as the subject
        "exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Refresh token expires in 7 days
        // Add other claims...
    }

    // Generate refresh token
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return refreshToken.SignedString([]byte("refresh_secret"))
}



func RefreshToken(c *fiber.Ctx) error {
    // Parse the JSON payload from the request body
    var refreshToken string
	fmt.Println("Raw request body data:", string(c.Body()))
    // if err := c.BodyParser(&refreshToken); err != nil {
    //     errMsg := ErrorResponse{Message: "Invalid JSON payload"}
    //     return c.Status(fiber.StatusBadRequest).JSON(errMsg)
    // }

	 // Get the raw request body data
	 bodyBytes := c.Request().Body()

	 // Parse the JSON data manually
	 var bodyMap map[string]interface{}
	 if err := json.Unmarshal(bodyBytes, &bodyMap); err != nil {
		 errMsg := ErrorResponse{Message: "Invalid JSON payload"}
		 return c.Status(fiber.StatusBadRequest).JSON(errMsg)
	 }
 
	 // Extract the refresh token from the parsed JSON data
	 refreshToken, ok := bodyMap["refresh_token"].(string)

	 if !ok {
		 errMsg := ErrorResponse{Message: "Refresh token not found or invalid format"}
		 return c.Status(fiber.StatusBadRequest).JSON(errMsg)
	 }

    // Verify the refresh token
    claims := jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte("refresh_secret"), nil
    })
    if err != nil {
        errMsg := ErrorResponse{Message: "Unauthorized: Invalid token"}
        return c.Status(fiber.StatusUnauthorized).JSON(errMsg)
    }
    if !token.Valid {
        errMsg := ErrorResponse{Message: "Unauthorized: Token is not valid"}
        return c.Status(fiber.StatusUnauthorized).JSON(errMsg)
    }

    // Extract the username from the refresh token claims
    username, ok := claims["sub"].(string)
    if !ok {
        errMsg := ErrorResponse{Message: "Unauthorized: Invalid token claims"}
        return c.Status(fiber.StatusUnauthorized).JSON(errMsg)
    }

    // Generate new access token
    accessToken, err := generateAccessToken(username)
    if err != nil {
        errMsg := ErrorResponse{Message: "Internal Server Error: Unable to generate access token"}
        return c.Status(fiber.StatusInternalServerError).JSON(errMsg)
    }

    // Generate new refresh token (optional: you may choose to reuse the existing refresh token)
    refreshToken, err = generateRefreshToken(username)
    if err != nil {
        errMsg := ErrorResponse{Message: "Internal Server Error: Unable to generate refresh token"}
        return c.Status(fiber.StatusInternalServerError).JSON(errMsg)
    }

    // Send new tokens to the client
    return c.JSON(fiber.Map{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}




// func accessible(c *fiber.Ctx) error {
// 	return c.SendString("Accessible ")
// }


// func restricted(c *fiber.Ctx) error {
// 	user := c.Locals("user").(*jwt.Token)
// 	claims := user.Claims.(jwt.MapClaims)
// 	name := claims["name"].(string)
// 	return c.SendString("Welcome " + name)
// }





