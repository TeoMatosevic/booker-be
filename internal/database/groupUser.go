package database

import (
	"database/sql"
)

func CreateGroupsUsersTable(db *sql.DB) error {
	var err error

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS group_users (
		id TEXT NOT NULL PRIMARY KEY,
		group_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetGroupUsersTableName() string {
	return groupsUsersTable
}

func (s *Service) GetAllGroupUsersByGroupID(groupID string) ([]GroupUser, error) {
	s.m.Lock()
	defer s.m.Unlock()

	rows, err := s.db.Query("SELECT * FROM "+s.groupsUsersTable+" WHERE group_id = ?", groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []GroupUser
	for rows.Next() {
		var result GroupUser
		if err := rows.Scan(&result.ID, &result.GroupID, &result.UserID); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *Service) GetAllGroupUsersByUserID(userID string) ([]GroupUser, error) {
	s.m.Lock()
	defer s.m.Unlock()

	rows, err := s.db.Query("SELECT * FROM "+s.groupsUsersTable+" WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []GroupUser
	for rows.Next() {
		var result GroupUser
		if err := rows.Scan(&result.ID, &result.GroupID, &result.UserID); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *Service) GetGroupUserByUserIDAndGroupID(userID, groupID string) (GroupUser, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var result GroupUser
	err := s.db.QueryRow("SELECT * FROM "+s.groupsUsersTable+" WHERE user_id = ? AND group_id = ?", userID, groupID).Scan(
		&result.ID,
		&result.GroupID,
		&result.UserID)
	if err != nil {
		return GroupUser{}, err
	}
	return result, nil
}

func (s *Service) InsertGroupUser(result GroupUser) error {
	s.m.Lock()
	defer s.m.Unlock()

	_, err := s.db.Exec("INSERT INTO "+s.groupsUsersTable+" (id, group_id, user_id) VALUES (?, ?, ?)",
		result.ID, result.GroupID, result.UserID)
	if err != nil {
		return err
	}

	return nil
}
