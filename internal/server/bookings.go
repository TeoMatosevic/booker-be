package server

import (
	"booker-be/internal/database"
	"booker-be/internal/protocol"

	"github.com/gin-gonic/gin"
)

func GetBookingsByPropertyID(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		propertyID := c.Param("propertyID")
		bookings, err := db.GetBookingsByPropertyID(propertyID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to retrieve bookings"})
			return
		}
		c.JSON(200, bookings)
	}
}

func CreateBooking(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		var booking database.Booking
		if err := c.ShouldBindJSON(&booking); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		bookingID := c.Param("bookingID")
		if booking.ID != bookingID {
			c.JSON(403, gin.H{"error": "You are not allowed to update this booking"})
			return
		}

		err := db.UpdateBooking(booking)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to update booking"})
			return
		}
		c.JSON(200, gin.H{"message": "Booking updated successfully"})
	}
}

func DeleteBooking(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookingID := c.Param("bookingID")
		err := db.DeleteBooking(bookingID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete booking"})
			return
		}
		c.JSON(200, gin.H{"message": "Booking deleted successfully"})
	}
}
