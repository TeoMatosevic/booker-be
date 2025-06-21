package main

import (
	"booker-be/internal/database"
	"booker-be/internal/server"
	"booker-be/internal/session"
)

func main() {
	// Initialize the database service
	db := database.New()

	// Initialize the session store
	store := session.NewStore()

	// Create a new Gin router
	server.StartServer(db, store)
}
