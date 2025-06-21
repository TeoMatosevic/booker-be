package database

import (
	"database/sql"
	"fmt"
)

func CreateUsersTable(db *sql.DB) error {
	var err error

	sqlStmt := `
	create table if not exists users (
		id string not null primary key,
		username string not null,
		hashed_password string not null
	);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUsersTableName() string {
	return s.usersTable
}

func (s *Service) GetUserByID(id string) (User, error) {
	s.m.Lock()
	defer s.m.Unlock()
	var result User
	err := s.db.QueryRow("SELECT * FROM "+s.usersTable+" WHERE id = ?", id).Scan(
		&result.ID,
		&result.Username,
		&result.HashedPassword)
	if err != nil {
		return User{}, err
	}
	return result, nil
}

func (s *Service) GetUserByUsername(username string) (User, error) {
	s.m.Lock()
	defer s.m.Unlock()
	var result User
	err := s.db.QueryRow("SELECT * FROM "+s.usersTable+" WHERE username = ?", username).Scan(
		&result.ID,
		&result.Username,
		&result.HashedPassword)
	if err != nil {
		fmt.Println("Error retrieving user by username:", err)
		return User{}, err
	}
	return result, nil
}

func (s *Service) InsertUser(result User) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.db.Exec("INSERT INTO "+s.usersTable+
		" (id, username, hashed_password) VALUES (?, ?, ?)",
		result.ID,
		result.Username,
		result.HashedPassword)

	if err != nil {
		return err
	}

	return nil
}
