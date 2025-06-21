package database

import (
	"database/sql"
	"strings"
)

func CreateBookingsTable(db *sql.DB) error {
	var err error

	sqlStmt := `
	create table if not exists bookings (
		id string not null primary key,
		created_at string,
		created_by string,
		property_id string not null,
		start_date string,
		end_date string,
		guest_name string
	);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetBookingsTableName() string {
	return s.bookingsTable
}

func (s *Service) GetAllBookings() ([]Booking, error) {
	s.m.Lock()
	defer s.m.Unlock()
	rows, err := s.db.Query("SELECT * FROM " + s.bookingsTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Booking
	for rows.Next() {
		var result Booking
		if err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.CreatedBy,
			&result.PropertyID,
			&result.StartDate,
			&result.EndDate,
			&result.GuestName); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *Service) GetBookingByID(id string) (Booking, error) {
	s.m.Lock()
	defer s.m.Unlock()
	var result Booking
	err := s.db.QueryRow("SELECT * FROM "+s.bookingsTable+" WHERE id = ?", id).Scan(
		&result.ID,
		&result.CreatedAt,
		&result.CreatedBy,
		&result.PropertyID,
		&result.StartDate,
		&result.EndDate,
		&result.GuestName)
	if err != nil {
		return Booking{}, err
	}
	return result, nil
}

func (s *Service) GetBookingsByPropertyID(propertyID string) ([]Booking, error) {
	s.m.Lock()
	defer s.m.Unlock()
	rows, err := s.db.Query("SELECT * FROM "+s.bookingsTable+" WHERE property_id = ?", propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []Booking
	for rows.Next() {
		var result Booking
		if err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.CreatedBy,
			&result.PropertyID,
			&result.StartDate,
			&result.EndDate,
			&result.GuestName); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (s *Service) GetBookingsByPropertyIds(propertyIDs []string) ([]Booking, error) {
	s.m.Lock()
	defer s.m.Unlock()

	query := "SELECT * FROM " + s.bookingsTable + " WHERE property_id IN (?" + strings.Repeat(",?", len(propertyIDs)-1) + ")"
	args := make([]interface{}, len(propertyIDs))
	for i, id := range propertyIDs {
		args[i] = id
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []Booking
	for rows.Next() {
		var result Booking
		if err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.CreatedBy,
			&result.PropertyID,
			&result.StartDate,
			&result.EndDate,
			&result.GuestName); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (s *Service) InsertBooking(result Booking) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.db.Exec("INSERT INTO "+s.bookingsTable+
		" (id, created_at, created_by, property_id, start_date, end_date, guest_name) VALUES (?, ?, ?, ?, ?, ?, ?)",
		result.ID,
		result.CreatedAt,
		result.CreatedBy,
		result.PropertyID,
		result.StartDate,
		result.EndDate,
		result.GuestName)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateBooking(result Booking) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.db.Exec("UPDATE "+s.bookingsTable+
		" SET start_date = ?, end_date = ?, guest_name = ? WHERE id = ?",
		result.StartDate,
		result.EndDate,
		result.GuestName,
		result.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteBooking(id string) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.db.Exec("DELETE FROM "+s.bookingsTable+" WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
