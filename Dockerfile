# Use the official Golang image as a base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Initialize Go module
RUN go mod init myapp

# Copy all the source code to the Working Directory
COPY code/ ./code/

# Install GoFiber framework
RUN go get -u github.com/gofiber/fiber/v2
RUN go get -u gorm.io/gorms
RUN go get -u gorm.io/driver/mysql
RUN go install github.com/cosmtrek/air@latest
RUN go get -u github.com/gofiber/contrib/jwt
RUN go get -u github.com/golang-jwt/jwt/v5
RUN go get golang.org/x/crypto/bcrypt
RUN go get myapp/code/credentials
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN go get -u github.com/gofiber/swagger


# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["sh", "-c", "cd code && air"]

