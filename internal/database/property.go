package database

import (
	"database/sql"
)

func CreatePropertyTable(db *sql.DB) error {
	var err error

	sqlStmt := `
	create table if not exists properties (
		id string not null primary key,
		created_at string,
		group_id string not null,
		name string not null,
		color string
	);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	// Add color column if it doesn't exist (for existing databases)
	// This will fail silently if the column already exists
	_, _ = db.Exec(`ALTER TABLE properties ADD COLUMN color string DEFAULT '';`)

	return nil
}

func (s *Service) GetPropertyTableName() string {
	return propertyTable
}

func (s *Service) GetAllProperties() ([]Property, error) {
	s.m.Lock()
	defer s.m.Unlock()
	rows, err := s.db.Query("SELECT * FROM " + propertyTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Property
	for rows.Next() {
		var result Property
		if err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.GroupID,
			&result.Name,
			&result.Color); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *Service) GetPropertyByID(id string) (Property, error) {
	s.m.Lock()
	defer s.m.Unlock()
	var result Property
	err := s.db.QueryRow("SELECT * FROM "+propertyTable+" WHERE id = ?", id).Scan(
		&result.ID,
		&result.CreatedAt,
		&result.GroupID,
		&result.Name,
		&result.Color)
	if err != nil {
		return Property{}, err
	}
	return result, nil
}

func (s *Service) InsertProperty(result Property) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.db.Exec("INSERT INTO "+propertyTable+
		" (id, created_at, group_id, name, color) VALUES (?, ?, ?, ?, ?)",
		result.ID,
		result.CreatedAt,
		result.GroupID,
		result.Name,
		result.Color)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetPropertiesByGroupID(groupID string) ([]Property, error) {
	s.m.Lock()
	defer s.m.Unlock()
	rows, err := s.db.Query("SELECT * FROM "+propertyTable+" WHERE group_id = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Property
	for rows.Next() {
		var result Property
		if err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.GroupID,
			&result.Name,
			&result.Color); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *Service) DeletePropertyByID(id string) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.db.Exec("DELETE FROM "+propertyTable+" WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdatePropertyColor(id string, color string) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.db.Exec("UPDATE "+propertyTable+" SET color = ? WHERE id = ?", color, id)
	return err
}
