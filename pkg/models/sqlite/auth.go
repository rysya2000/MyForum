package sqlite

import (
	"git.01.alem.school/rysya2000/MyForum.git/pkg/models"
)

func (s *ForumModel) InsertUser(u models.User) (int64, error) {
	stmt := `INSERT
			INTO user
			(username, email, pass)
			VALUES (?, ?, ?)`
	row, err := s.DB.Exec(stmt, u.Username, u.Email, u.Password)
	if err != nil {
		return 0, err
	}
	return row.LastInsertId()
}

func (s *ForumModel) GetUserByName(username string) (*models.User, error) {
	stmt := `SELECT userid, username, email, pass FROM user WHERE username = ?`

	user := &models.User{}

	row := s.DB.QueryRow(stmt, username)
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *ForumModel) GetUserByEmail(email string) (*models.User, error) {
	stmt := `SELECT userid, username, email, pass FROM user WHERE email = ?`

	user := &models.User{}

	row := s.DB.QueryRow(stmt, email)
	err := row.Scan(&user.UserID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
