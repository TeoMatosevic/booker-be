package server

import (
	"booker-be/internal/database"
	"booker-be/internal/protocol"
	"booker-be/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterUser(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user protocol.CreateUserMessage
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// Check if the username already exists
		existingUser, err := db.GetUserByUsername(user.Username)
		if err == nil && existingUser.Username != "" {
			c.JSON(409, gin.H{"error": "Username already exists"})
			return
		}

		hashedPassword := protocol.Sha256Hash(user.Password)
		id := protocol.GenerateID()

		err = db.InsertUser(database.User{
			ID:             id,
			Username:       user.Username,
			HashedPassword: hashedPassword,
		})

		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to register user"})
			return
		}

		c.JSON(201, gin.H{"message": "User registered successfully"})
	}
}

func LoginUser(db database.Service, sessionStore *session.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user protocol.LoginUserMessage
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		dbUser, err := db.GetUserByUsername(user.Username)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			return
		}

		if dbUser.HashedPassword != protocol.Sha256Hash(user.Password) {
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			return
		}

		token, err := sessionStore.CreateSession(dbUser.ID, time.Hour*24) // Create a session valid for 24 hours
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create session"})
			return
		}

		c.JSON(200, gin.H{
			"message":  "Login successful",
			"token":    token,
			"userID":   dbUser.ID,
			"username": dbUser.Username,
		})
	}
}
