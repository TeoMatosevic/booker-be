package server

import (
	"booker-be/internal/database"
	"booker-be/internal/protocol"
	"math/rand/v2"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	groupCodeDuration = 24 * 60 * 60 // 24 hours in seconds
)

func generateGroupCode() string {
	code := make([]rune, 6) // Generate a 6-character code
	for i := range code {
		code[i] = rune(rand.IntN(26) + 'A') // Randomly select a character from A-Z
	}

	return string(code)
}

func CreateGroupCode(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var groupCode protocol.GroupCodeMessage
		if err := c.ShouldBindJSON(&groupCode); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		gCode := database.GroupCode{
			ID:       protocol.GenerateID(),
			Code:     generateGroupCode(),
			GroupID:  groupCode.GroupID,
			ActiveTo: (time.Now().Add(time.Duration(groupCodeDuration) * time.Second)).Format(time.RFC3339),
		}

		err := db.InsertGroupCode(gCode)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create group code"})
			return
		}
		c.JSON(201, gin.H{
			"message":   "Group code created successfully",
			"groupCode": gCode.Code,
		})
	}
}
