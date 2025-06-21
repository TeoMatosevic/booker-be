package database

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type Service struct {
	db               *sql.DB
	m                *sync.Mutex
	bookingsTable    string
	usersTable       string
	propertyTable    string
	groupsTable      string
	groupsUsersTable string
	groupCodesTable  string
}

var (
	usersTable       = "users"
	groupsTable      = "groups"
	propertyTable    = "properties"
	bookingsTable    = "bookings"
	groupsUsersTable = "group_users"
	groupCodesTable  = "group_codes"

	dbInstance *Service
)

func New() Service {
	var err error
	db, err := sql.Open("sqlite3", "./bookings.db")
	if err != nil {
		panic(err)
	}

	// Create the bookings table if it doesn't exist
	err = CreateBookingsTable(db)
	if err != nil {
		panic(err)
	}

	// Create the users table if it doesn't exist
	err = CreateUsersTable(db)
	if err != nil {
		panic(err)
	}

	// Create the properties table if it doesn't exist
	err = CreatePropertyTable(db)
	if err != nil {
		panic(err)
	}

	// Create the groups table if it doesn't exist
	err = CreateGroupsTable(db)
	if err != nil {
		panic(err)
	}

	// Create the groups_users table if it doesn't exist
	err = CreateGroupsUsersTable(db)
	if err != nil {
		panic(err)
	}

	// Create the group_codes table if it doesn't exist
	err = CreateGroupCodesTable(db)
	if err != nil {
		panic(err)
	}

	dbInstance = &Service{
		db:               db,
		bookingsTable:    bookingsTable,
		usersTable:       usersTable,
		propertyTable:    propertyTable,
		groupsTable:      groupsTable,
		groupsUsersTable: groupsUsersTable,
		groupCodesTable:  groupCodesTable,
		m:                &sync.Mutex{},
	}

	go func() {
		oneHour := 3600
		for {
			// Clean up expired group codes every hour
			err := dbInstance.CleanUpExpiredGroupCodes()
			if err != nil {
				panic(err)
			}
			// Sleep for one hour
			time.Sleep(time.Duration(oneHour) * time.Second)
		}
	}()

	return *dbInstance
}

func (s *Service) Close() error {
	return s.db.Close()
}
