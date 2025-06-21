package server

import (
	"booker-be/internal/database"
	"booker-be/internal/protocol"

	"github.com/gin-gonic/gin"
)

func GetPropertiesByGroupID(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		groupID := c.Param("groupID")
		properties, err := db.GetPropertiesByGroupID(groupID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve properties"})
			return
		}
		c.JSON(200, properties)
	}
}

func CreateProperty(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var property protocol.CreatePropertyMessage
		if err := c.ShouldBindJSON(&property); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// Get the ownerID from the URL parameter
		groupID := c.Param("groupID")
		if property.GroupID != groupID {
			c.JSON(403, gin.H{"error": "You are not allowed to create properties for this owner"})
			return
		}

		p := database.Property{
			ID:        protocol.GenerateID(),
			Name:      property.Name,
			GroupID:   groupID,
			CreatedAt: protocol.GetCurrentTime(),
		}

		err := db.InsertProperty(p)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create property"})
			return
		}
		c.JSON(201, gin.H{"message": "Property created successfully"})
	}
}
