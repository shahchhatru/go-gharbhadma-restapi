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
    "crypto/rand"
    "encoding/base64"
    
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
    

type ResetPasswordRequest struct {
    Email string `json:"email"`
}

// Define the expiration time for reset tokens (e.g., 1 hour)
const resetTokenExpiration = time.Hour

// GenerateResetToken generates a random reset token
func generateResetToken(email string) (string, error) {
    // Generate a random token
    tokenBytes := make([]byte, 32)
    _, err := rand.Read(tokenBytes)
    if err != nil {
        return "", err
    }
    token := base64.StdEncoding.EncodeToString(tokenBytes)

    // You can save the reset token in a database or cache with the associated email
    // Example: SaveTokenToDatabase(token, email)

    // Return the generated token
    return token, nil
}

func ResetPasswordRequestHandler(c *fiber.Ctx) error {
    // Parse the JSON payload from the request body
    var req ResetPasswordRequest
    if err := c.BodyParser(&req); err != nil {
        errMsg := ErrorResponse{Message: "Invalid JSON payload"}
        return c.Status(fiber.StatusBadRequest).JSON(errMsg)
    }

    // Check if the user exists in the database
    var user Credentials
    if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // User not found
            errMsg := ErrorResponse{Message: "User not found"}
            return c.Status(fiber.StatusNotFound).JSON(errMsg)
        }
        // Other database errors
        errMsg := ErrorResponse{Message: "Internal Server Error"}
        return c.Status(fiber.StatusInternalServerError).JSON(errMsg)
    }

    // Generate a reset token (you can implement this function)
    resetToken, err := generateResetToken(req.Email)
    if err != nil {
        errMsg := ErrorResponse{Message: "Failed to generate reset token"}
        return c.Status(fiber.StatusInternalServerError).JSON(errMsg)
    }

    // Send the reset token to the user via email (implement this function)

    // Return success response
    return c.JSON(fiber.Map{"message": "Reset token sent successfully","resettoken":resetToken})
}


type ResetPasswordConfirmation struct {
    Email         string `json:"email"`
    ResetToken    string `json:"reset_token"`
    NewPassword   string `json:"new_password"`
}

func validateResetToken(email string, resetToken string) error {
    // Parse the reset token
    claims := jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(resetToken, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte("reset_secret"), nil
    })
    if err != nil {
        return fmt.Errorf("Invalid reset token")
    }

    // Check if the token is valid
    if !token.Valid {
        return fmt.Errorf("Invalid reset token")
    }

    // Extract the email from the token claims
    tokenEmail, ok := claims["email"].(string)
    if !ok {
        return fmt.Errorf("Invalid reset token")
    }

    // Compare the token email with the provided email
    if tokenEmail != email {
        return fmt.Errorf("Reset token does not match the user's email")
    }

    // Token is valid and matches the user's email
    return nil
}


func ResetPasswordConfirmationHandler(c *fiber.Ctx) error {
    // Parse the JSON payload from the request body
    var req ResetPasswordConfirmation
    if err := c.BodyParser(&req); err != nil {
        errMsg := ErrorResponse{Message: "Invalid JSON payload"}
        return c.Status(fiber.StatusBadRequest).JSON(errMsg)
    }

     // Validate the reset token
     if err := validateResetToken(req.Email, req.ResetToken); err != nil {
        errMsg := ErrorResponse{Message: err.Error()}
        return c.Status(fiber.StatusUnauthorized).JSON(errMsg)
    }

    // Update the user's password in the database
    hashedPassword, err := hashPassword(req.NewPassword)
    if err != nil {
        errMsg := ErrorResponse{Message: "Failed to hash password"}
        return c.Status(fiber.StatusInternalServerError).JSON(errMsg)
    }

    if err := db.DB.Model(&Credentials{}).Where("email = ?", req.Email).Update("pass", hashedPassword).Error; err != nil {
        errMsg := ErrorResponse{Message: "Failed to update password"}
        return c.Status(fiber.StatusInternalServerError).JSON(errMsg)
    }

    // Return success response
    return c.JSON(fiber.Map{"message": "Password reset successfully"})
}

type ChangePasswordRequest struct {
    Email        string `json:"email"`
    OldPassword  string `json:"old_password"`
    NewPassword  string `json:"new_password"`
}

func ChangePasswordHandler(c *fiber.Ctx) error {
    // Parse the JSON payload from the request body
    var req ChangePasswordRequest
    if err := c.BodyParser(&req); err != nil {
        errMsg := ErrorResponse{Message: "Invalid JSON payload"}
        return c.Status(fiber.StatusBadRequest).JSON(errMsg)
    }

    // Check if the user exists in the database
    var user Credentials
    if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // User not found
            errMsg := ErrorResponse{Message: "User not found"}
            return c.Status(fiber.StatusNotFound).JSON(errMsg)
        }
        // Other database errors
        errMsg := ErrorResponse{Message: "Internal Server Error"}
        return c.Status(fiber.StatusInternalServerError).JSON(errMsg)
    }

    // Verify the old password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(req.OldPassword)); err != nil {
        // Old password does not match
        errMsg := ErrorResponse{Message: "Old password is incorrect"}
        return c.Status(fiber.StatusUnauthorized).JSON(errMsg)
    }

    // Hash the new password
    hashedPassword, err := hashPassword(req.NewPassword)
    if err != nil {
        errMsg := ErrorResponse{Message: "Failed to hash password"}
        return c.Status(fiber.StatusInternalServerError).JSON(errMsg)
    }

    // Update the user's password in the database
    if err := db.DB.Model(&Credentials{}).Where("email = ?", req.Email).Update("pass", hashedPassword).Error; err != nil {
        errMsg := ErrorResponse{Message: "Failed to update password"}
        return c.Status(fiber.StatusInternalServerError).JSON(errMsg)
    }

    // Return success response
    return c.JSON(fiber.Map{"message": "Password changed successfully"})
}








