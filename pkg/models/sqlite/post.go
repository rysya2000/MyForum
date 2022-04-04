package sqlite

import (
	"git.01.alem.school/rysya2000/MyForum.git/pkg/models"
)

func (s *ForumModel) InsertPost(post *models.Post) (int64, error) {
	stmt := `INSERT 
			INTO post 
			(title, author, content, created)
			VALUES(?, ?, ?, ?)`

	row, err := s.DB.Exec(stmt, post.Title, post.Author, post.Content, post.Created)
	if err != nil {
		return 0, err
	}

	return row.LastInsertId()
}

func (s *ForumModel) GetPostById(id int) (*models.Post, error) {
	stmt := `SELECT postid, title, author, content, created
			FROM post
			WHERE postid = $1`

	post := &models.Post{}

	row := s.DB.QueryRow(stmt, id)
	err := row.Scan(&post.PostID, &post.Title, &post.Author, &post.Content, &post.Created)
	if err != nil {
		return nil, err
	}

	stmt2 := `SELECT userid, symbol
	FROM like
	WHERE postid = $1`

	rows, err := s.DB.Query(stmt2, id)
	if err != nil {
		return nil, err
	}
	var (
		owner  int
		symbol string
	)
	for rows.Next() {
		err := rows.Scan(&owner, &symbol)
		if err != nil {
			return nil, err
		}

		if symbol == "1" {
			post.Like++
		} else if symbol == "0" {
			post.Dislike++
		}
	}

	stmt3 := `SELECT id, username, comment
	FROM comment
	WHERE postid = $1`

	rows, err = s.DB.Query(stmt3, id)
	if err != nil {
		return nil, err
	}

	c := models.Comment{}
	for rows.Next() {
		err = rows.Scan(&c.CommentID, &c.Author, &c.Content)
		if err != nil {
			return nil, err
		}
		post.Comments = append(post.Comments, c)
	}

	stmt4 := `SELECT tag
	FROM tag
	WHERE postid = $1`

	rows, err = s.DB.Query(stmt4, id)
	if err != nil {
		return nil, err
	}

	var tag string
	for rows.Next() {
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		post.Tags = append(post.Tags, tag)
	}

	stmt5 := `SELECT symbol
			FROM rateComment
			WHERE postid = $1 AND commentid = $2`
	for i, v := range post.Comments {
		rows, err := s.DB.Query(stmt5, post.PostID, v.CommentID)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var sym string
			err = rows.Scan(&sym)

			if sym == "1" {
				post.Comments[i].Like++
			}
			if sym == "2" {
				post.Comments[i].Dislike++
			}
		}

	}

	return post, nil
}

func (s *ForumModel) InsertRaiting(userid int, postid int, symbol string) error {
	stmt := `INSERT
			INTO like
			(userid, postid, symbol)
			VALUES(?, ?, ?)`
	_, err := s.DB.Exec(stmt, userid, postid, symbol)
	if err != nil {
		return err
	}
	return nil
}

func (s *ForumModel) GetRaiting(userid int, postid int) (string, error) {
	stmt := `SELECT symbol
			FROM like
			WHERE postid = $1`

	var symbol string

	row := s.DB.QueryRow(stmt, postid)
	err := row.Scan(&symbol)
	if err != nil {
		return "", err
	}

	return symbol, nil
}

func (s *ForumModel) DelRaiting(userid int, postid int) error {
	stmt := `DELETE 
			FROM like
			WHERE userid = $1 AND postid = $2`
	_, err := s.DB.Exec(stmt, userid, postid)
	if err != nil {
		return err
	}
	return nil
}

func (s *ForumModel) InsertComment(userid int, username string, postid int, text string) error {
	stmt := `INSERT
			INTO comment
			(userid, username, postid, comment)
			VALUES(?, ?, ?, ?)`
	_, err := s.DB.Exec(stmt, userid, username, postid, text)
	if err != nil {
		return err
	}

	return nil
}

func (s *ForumModel) InsertTag(postid int, tag string) error {
	stmt := `INSERT
			INTO tag
			(postid, tag)
			VALUES(?, ?)`
	_, err := s.DB.Exec(stmt, postid, tag)
	if err != nil {
		return err
	}

	return nil
}
