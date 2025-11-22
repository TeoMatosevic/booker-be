package server

import (
	"booker-be/internal/database"
	"booker-be/internal/session"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes the routes for the booking service
func SetupRoutes(router *gin.Engine, db database.Service, sessionValidator session.SessionValidator) {
	// CORS middleware with whitelisted origins
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Whitelist of allowed origins - configure based on environment
		allowedOrigins := map[string]bool{
			"http://localhost:5173":  true, // Vite dev server
			"http://localhost:3000":  true, // Alternative dev port
			"https://yourdomain.com": true, // Production domain (replace with actual domain)
		}

		// Check if origin is in whitelist
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Header("Content-Type", "application/json")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content
			return
		}
		c.Next()
	})

	users := router.Group("/users")
	{
		users.POST("/register", RegisterUser(db))
		users.POST("/login", LoginUser(db, sessionValidator.(*session.Store))) // Use type assertion to access Store methods
	}

	authMW := AuthMiddleware(sessionValidator) // Create the authentication middleware

	bookings := router.Group("/bookings")
	bookings.Use(authMW) // Apply authentication middleware
	{
		bookings.GET("/property/:propertyID", GetBookingsByPropertyID(db))
		bookings.GET("/group/:groupID", GetBookingsByGroupID(db))
		bookings.POST("/property/:propertyID", CreateBooking(db))
		bookings.PUT("/:bookingID", UpdateBooking(db))
		bookings.DELETE("/:bookingID", DeleteBooking(db))
	}

	groupCodes := router.Group("/group-codes")
	groupCodes.Use(authMW) // Apply authentication middleware
	{
		groupCodes.POST("/", CreateGroupCode(db))
	}

	groups := router.Group("/groups")
	groups.Use(authMW) // Apply authentication middleware
	{
		groups.GET("/:userID", GetGroupsByUserID(db))
		groups.POST("/", CreateGroup(db))
		groups.POST("/join/:code", JoinGroup(db))
	}

	properties := router.Group("/properties")
	properties.Use(authMW) // Apply authentication middleware
	{
		properties.GET("/group/:groupID", GetPropertiesByGroupID(db))
		properties.POST("/group/:groupID", CreateProperty(db))
		properties.PUT("/:propertyID", UpdateProperty(db))
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not Found"})
	})
}

// StartServer initializes the Gin router and starts the server
func StartServer(db database.Service, sessionStore session.SessionValidator) {
	router := gin.Default()
	SetupRoutes(router, db, sessionStore)

	if err := router.Run(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
