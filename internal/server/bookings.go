package server

import (
	"booker-be/internal/database"
	"booker-be/internal/protocol"

	"github.com/gin-gonic/gin"
)

func GetBookingsByPropertyID(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		propertyID := c.Param("propertyID")

		// Check if user belongs to the group that owns this property
		if !db.UserBelongsToPropertyGroup(userID.(string), propertyID) {
			c.JSON(403, gin.H{"error": "Forbidden: You don't have access to this property"})
			return
		}

		bookings, err := db.GetBookingsByPropertyID(propertyID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve bookings"})
			return
		}
		c.JSON(200, bookings)
	}
}

func GetBookingsByGroupID(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		groupID := c.Param("groupID")

		// Check if user belongs to the group
		if !db.UserBelongsToGroup(userID.(string), groupID) {
			c.JSON(403, gin.H{"error": "Forbidden: You don't have access to this group"})
			return
		}

		// Get all properties for the group
		properties, err := db.GetPropertiesByGroupID(groupID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve properties"})
			return
		}

		// Extract property IDs
		propertyIDs := make([]string, len(properties))
		for i, p := range properties {
			propertyIDs[i] = p.ID
		}

		// Get bookings for all properties
		if len(propertyIDs) == 0 {
			c.JSON(200, []database.Booking{})
			return
		}

		bookings, err := db.GetBookingsByPropertyIds(propertyIDs)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve bookings"})
			return
		}
		c.JSON(200, bookings)
	}
}

func CreateBooking(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		propertyID := c.Param("propertyID")
		var booking protocol.CreateBookingMessage
		if err := c.ShouldBindJSON(&booking); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if propertyID == "" {
			c.JSON(400, gin.H{"error": "Property ID is required"})
			return
		}

		// Check if the property exists
		_, err := db.GetPropertyByID(propertyID)
		if err != nil {
			c.JSON(404, gin.H{"error": "Property not found"})
			return
		}

		// Check if user belongs to the group that owns this property
		if !db.UserBelongsToPropertyGroup(userID.(string), propertyID) {
			c.JSON(403, gin.H{"error": "Forbidden: You don't have access to this property"})
			return
		}

		// Check if the booking dates are valid (they are strings)
		if booking.StartDate == "" || booking.EndDate == "" {
			c.JSON(400, gin.H{"error": "Start date and end date are required"})
			return
		}

		// Check if the booking dates are in the correct format
		if !protocol.IsValidDate(booking.StartDate) || !protocol.IsValidDate(booking.EndDate) {
			c.JSON(400, gin.H{"error": "Invalid date format"})
			return
		}

		// Check if end date is after start date
		eD, err := protocol.ParseDate(booking.EndDate)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid end date format"})
			return
		}

		sD, err := protocol.ParseDate(booking.StartDate)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid start date format"})
			return
		}

		if eD.Before(sD) {
			c.JSON(400, gin.H{"error": "End date must be after start date"})
			return
		}

		b := database.Booking{
			ID:         protocol.GenerateID(),
			PropertyID: propertyID,
			StartDate:  booking.StartDate,
			EndDate:    booking.EndDate,
			GuestName:  booking.GuestName,
			Adults:     booking.Adults,
			Children:   booking.Children,
		}

		err = db.InsertBooking(b)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create booking"})
			return
		}
		c.JSON(201, gin.H{"message": "Booking created successfully"})
	}
}

func UpdateBooking(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		var booking protocol.UpdateBookingMessage
		if err := c.ShouldBindJSON(&booking); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		bookingID := c.Param("bookingID")

		// Check if user can access this booking
		if !db.UserCanAccessBooking(userID.(string), bookingID) {
			c.JSON(403, gin.H{"error": "Forbidden: You don't have access to this booking"})
			return
		}

		b := database.Booking{
			ID:         bookingID,
			PropertyID: "", // This is not needed for update
			StartDate:  booking.StartDate,
			EndDate:    booking.EndDate,
			GuestName:  booking.GuestName,
			Adults:     booking.Adults,
			Children:   booking.Children,
		}

		err := db.UpdateBooking(b)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to update booking"})
			return
		}
		c.JSON(200, gin.H{"message": "Booking updated successfully"})
	}
}

func DeleteBooking(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		bookingID := c.Param("bookingID")

		// Check if user can access this booking
		if !db.UserCanAccessBooking(userID.(string), bookingID) {
			c.JSON(403, gin.H{"error": "Forbidden: You don't have access to this booking"})
			return
		}

		err := db.DeleteBooking(bookingID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete booking"})
			return
		}
		c.JSON(200, gin.H{"message": "Booking deleted successfully"})
	}
}
