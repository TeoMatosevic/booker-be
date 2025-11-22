package server

import (
	"booker-be/internal/database"
	"booker-be/internal/protocol"

	"github.com/gin-gonic/gin"
)

func GetPropertiesByGroupID(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		groupID := c.Param("groupID")

		// Check if user belongs to this group
		if !db.UserBelongsToGroup(userID.(string), groupID) {
			c.JSON(403, gin.H{"error": "Forbidden: You don't have access to this group"})
			return
		}

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
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		var property protocol.CreatePropertyMessage
		if err := c.ShouldBindJSON(&property); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// Get the groupID from the URL parameter
		groupID := c.Param("groupID")
		if property.GroupID != groupID {
			c.JSON(403, gin.H{"error": "You are not allowed to create properties for this owner"})
			return
		}

		// Check if user belongs to this group
		if !db.UserBelongsToGroup(userID.(string), groupID) {
			c.JSON(403, gin.H{"error": "Forbidden: You don't have access to this group"})
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

func UpdateProperty(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		propertyID := c.Param("propertyID")

		// Get the property to check group ownership
		property, err := db.GetPropertyByID(propertyID)
		if err != nil {
			c.JSON(404, gin.H{"error": "Property not found"})
			return
		}

		// Check if user belongs to this property's group
		if !db.UserBelongsToGroup(userID.(string), property.GroupID) {
			c.JSON(403, gin.H{"error": "Forbidden: You don't have access to this property"})
			return
		}

		var updateMsg protocol.UpdatePropertyMessage
		if err := c.ShouldBindJSON(&updateMsg); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// Validate color format if provided (must be hex color or empty)
		if updateMsg.Color != "" && !isValidHexColor(updateMsg.Color) {
			c.JSON(400, gin.H{"error": "Invalid color format. Must be hex color (e.g., #FF5733)"})
			return
		}

		// Update color
		if err := db.UpdatePropertyColor(propertyID, updateMsg.Color); err != nil {
			c.JSON(500, gin.H{"error": "Failed to update property"})
			return
		}

		c.JSON(200, gin.H{"message": "Property updated successfully"})
	}
}

func isValidHexColor(color string) bool {
	// Check if it matches hex color format: #RRGGBB
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	for i := 1; i < 7; i++ {
		c := color[i]
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
