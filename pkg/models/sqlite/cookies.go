package sqlite

import (
	"time"
)

func (s *ForumModel) InsertCookie(userID int, uuid string) error {
	stmt := `INSERT
			INTO cookie
			(userid, uuid, expirydate)
			VALUES(?, ?, ?)`
	Now := time.Now().Format("01-02-2006 15:04:05")
	_, err := s.DB.Exec(stmt, userID, uuid, Now)
	if err != nil {
		return err
	}

	return nil
}

func (s *ForumModel) DelCookie(userID int) {
	stmt := `DELETE
			FROM cookie
			WHERE userid = $1`
	s.DB.Exec(stmt, userID)
}

func (s *ForumModel) IsCookieInDB(userID int) error {
	stmt := `SELECT userid
			FROM cookie
			WHERE userid = $1`

	var id int

	row := s.DB.QueryRow(stmt, userID)
	err := row.Scan(&id)
	if err != nil {
		return err
	}

	return nil
}
