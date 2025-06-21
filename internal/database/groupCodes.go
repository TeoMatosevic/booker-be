package database

import (
	"database/sql"
	"fmt"
	"time"
)

func CreateGroupCodesTable(db *sql.DB) error {
	sqlStmt := `
	create table if not exists group_codes (
		id string not null primary key,
		group_id string not null,
		code string not null,
		active_to string not null
	);
	`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetGroupCodesTableName() string {
	return s.groupCodesTable
}

func (s *Service) GetAllGroupCodes() ([]GroupCode, error) {
	s.m.Lock()
	defer s.m.Unlock()

	rows, err := s.db.Query("SELECT * FROM " + s.groupCodesTable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []GroupCode
	for rows.Next() {
		var result GroupCode
		if err := rows.Scan(
			&result.ID,
			&result.GroupID,
			&result.Code,
			&result.ActiveTo); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *Service) GetGroupCodeByID(id string) (GroupCode, error) {
	s.m.Lock()
	defer s.m.Unlock()
	var result []GroupCode
	err := s.db.QueryRow("SELECT * FROM "+s.groupCodesTable+" WHERE id = ?", id).Scan(
		&result[0].ID,
		&result[0].GroupID,
		&result[0].Code,
		&result[0].ActiveTo)
	if err != nil {
		return GroupCode{}, err
	}
	return result[0], nil
}

func (s *Service) GetGroupCodeByCode(code string) (GroupCode, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var result GroupCode
	err := s.db.QueryRow("SELECT * FROM "+s.groupCodesTable+" WHERE code = ?", code).Scan(
		&result.ID,
		&result.GroupID,
		&result.Code,
		&result.ActiveTo)
	if err != nil {
		return GroupCode{}, err
	}
	return result, nil
}

func (s *Service) InsertGroupCode(result GroupCode) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, err := s.db.Exec("INSERT INTO "+s.groupCodesTable+
		" (id, group_id, code, active_to) VALUES (?, ?, ?, ?)",
		result.ID,
		result.GroupID,
		result.Code,
		result.ActiveTo)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CleanUpExpiredGroupCodes() error {
	fmt.Println("Cleaning up expired group codes, current time:", time.Now().Format(time.RFC3339))
	codes, err := s.GetAllGroupCodes()
	fmt.Println("Retrieved group codes for cleanup:", len(codes))
	if err != nil {
		return err
	}
	for _, code := range codes {
		if code.ActiveTo < time.Now().Format(time.RFC3339) {
			s.m.Lock()
			_, err := s.db.Exec("DELETE FROM "+s.groupCodesTable+" WHERE id = ?", code.ID)
			s.m.Unlock()
			if err != nil {
				return err
			}
		}
	}

	fmt.Println("Expired group codes cleanup completed.")

	return nil
}
