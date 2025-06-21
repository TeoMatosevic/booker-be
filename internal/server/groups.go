package server

import (
	"booker-be/internal/database"
	"booker-be/internal/protocol"

	"github.com/gin-gonic/gin"
)

func GetGroupsByUserID(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		groupUsers, err := db.GetAllGroupUsersByUserID(userID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve groups"})
			return
		}

		gIds := make([]string, 0, len(groupUsers))
		for _, gu := range groupUsers {
			gIds = append(gIds, gu.GroupID)
		}
		groups, err := db.GetGroupsByID(gIds)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve groups"})
			return
		}

		c.JSON(200, groups)
	}
}

func CreateGroup(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var group protocol.GroupCreateMessage
		if err := c.ShouldBindJSON(&group); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		userID, exists := c.Get(authorizationPayloadKey)
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		uID, ok := userID.(string)
		if !ok {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		g := database.Group{
			ID:        protocol.GenerateID(),
			Name:      group.Name,
			OwnerID:   uID,
			CreatedAt: protocol.GetCurrentTime(),
		}

		err := db.InsertGroup(g)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create group"})
			return
		}

		groupUser := database.GroupUser{
			ID:      protocol.GenerateID(),
			GroupID: g.ID,
			UserID:  uID,
		}

		err = db.InsertGroupUser(groupUser)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to add user to group"})
			return
		}

		c.JSON(201, gin.H{"message": "Group created successfully"})
	}
}

func JoinGroup(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
		userID, exists := c.Get(authorizationPayloadKey)
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		uID, ok := userID.(string)
		if !ok {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		groupCode, err := db.GetGroupCodeByCode(code)
		if err != nil {
			c.JSON(404, gin.H{"error": "Group code not found"})
			return
		}

		if groupCode.ActiveTo < protocol.GetCurrentTime() {
			c.JSON(400, gin.H{"error": "Group code has expired"})
			return
		}

		// Check if the user is already a member of the group
		_, err = db.GetGroupUserByUserIDAndGroupID(uID, groupCode.GroupID)
		if err == nil {
			c.JSON(400, gin.H{"error": "User is already a member of this group"})
			return
		}

		groupUser := database.GroupUser{
			ID:      protocol.GenerateID(),
			GroupID: groupCode.GroupID,
			UserID:  uID,
		}
		err = db.InsertGroupUser(groupUser)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to join group"})
			return
		}

		c.JSON(200, gin.H{"message": "Joined group successfully"})
	}
}
