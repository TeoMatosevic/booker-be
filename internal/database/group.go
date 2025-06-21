package database

import (
	"database/sql"
	"strings"
)

func CreateGroupsTable(db *sql.DB) error {
	sqlStmt := `
	create table if not exists groups (
		id string not null primary key,
		created_at string,
		name string not null,
		owner_id string not null
	);
	`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetGroupsTableName() string {
	return s.groupsTable
}

func (s *Service) GetAllGroups() ([]Group, error) {
	s.m.Lock()
	defer s.m.Unlock()

	rows, err := s.db.Query("SELECT * FROM " + s.groupsTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Group
	for rows.Next() {
		var result Group
		if err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.Name,
			&result.OwnerID); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *Service) GetGroupByID(id string) (Group, error) {
	s.m.Lock()
	defer s.m.Unlock()
	var result Group
	err := s.db.QueryRow("SELECT * FROM "+s.groupsTable+" WHERE id = ?", id).Scan(
		&result.ID,
		&result.CreatedAt,
		&result.Name,
		&result.OwnerID)
	if err != nil {
		return Group{}, err
	}
	return result, nil
}

func (s *Service) GetGroupsByID(ids []string) ([]Group, error) {
	s.m.Lock()
	defer s.m.Unlock()
	if len(ids) == 0 {
		return nil, nil // Return empty slice if no IDs are provided
	}
	questionMarks := strings.Repeat("?,", len(ids))
	questionMarks = strings.TrimSuffix(questionMarks, ",") // Remove trailing comma
	query := "SELECT * FROM " + s.groupsTable + " WHERE id IN (" + questionMarks + ")"
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []Group
	for rows.Next() {
		var result Group
		if err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.Name,
			&result.OwnerID); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (s *Service) GetGroupByOwnerID(ownerID string) ([]Group, error) {
	s.m.Lock()
	defer s.m.Unlock()

	rows, err := s.db.Query("SELECT * FROM "+s.groupsTable+" WHERE owner_id = ?", ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Group
	for rows.Next() {
		var result Group
		if err := rows.Scan(
			&result.ID,
			&result.CreatedAt,
			&result.Name,
			&result.OwnerID); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *Service) InsertGroup(result Group) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.db.Exec("INSERT INTO "+s.groupsTable+
		" (id, created_at, name, owner_id) VALUES (?, ?, ?, ?)",
		result.ID,
		result.CreatedAt,
		result.Name,
		result.OwnerID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteGroupByID(id string) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.db.Exec("DELETE FROM "+s.groupsTable+" WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
