package sqlite

func (s *ForumModel) InsertRateComment(userid int, postid int, commentid int, symbol string) error {
	stmt := `INSERT
			INTO rateComment
			(userid, commentid, postid, symbol)
			VALUES(?, ?, ?, ?)`
	_, err := s.DB.Exec(stmt, userid, commentid, postid, symbol)
	if err != nil {
		return err
	}

	return nil
}

func (s *ForumModel) GetRateComment(userid int, postid int, commentid int) (string, error) {
	stmt := `SELECT symbol
			FROM rateComment
			WHERE commentid = $1 AND postid = $2`

	var symbol string

	row := s.DB.QueryRow(stmt, commentid, postid)
	err := row.Scan(&symbol)
	if err != nil {
		return "", err
	}

	return symbol, nil
}

func (s *ForumModel) DelRateComment(userid int, commentid int) error {
	stmt := `DELETE
			FROM rateComment
			WHERE userid = $1 AND commentid = $2`
	_, err := s.DB.Exec(stmt, userid, commentid)
	if err != nil {
		return err
	}
	return nil
}
