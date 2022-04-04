package sqlite

import (
	"git.01.alem.school/rysya2000/MyForum.git/pkg/models"
)

func (s *ForumModel) GetUserByUuid(uuid string) (*models.User, error) {
	stmt := `SELECT userid FROM cookie WHERE uuid = ?`

	user := models.User{}

	row := s.DB.QueryRow(stmt, uuid)
	err := row.Scan(&user.UserID)
	if err != nil {
		return nil, err
	}

	stmt2 := `SELECT username, email, pass FROM user WHERE userid = ?`

	user2 := &models.User{}

	row2 := s.DB.QueryRow(stmt2, user.UserID)
	err = row2.Scan(&user2.Username, &user2.Email, &user2.Password)
	if err != nil {
		return nil, err
	}

	user2.UserID = user.UserID

	return user2, nil
}

func (s *ForumModel) GetMyPosts(username string) ([]*models.Post, error) {
	stmt := `SELECT postid, author, title, content
			FROM post
			WHERE author = $1`
	posts := []*models.Post{}

	rows, err := s.DB.Query(stmt, username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		post := &models.Post{}

		err = rows.Scan(&post.PostID, &post.Author, &post.Title, &post.Content)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *ForumModel) GetLikedPosts(userid int) ([]*models.Post, error) {
	stmt := `SELECT postid
			FROM like
			WHERE userid = $1 AND symbol = $2`
	rows, err := s.DB.Query(stmt, userid, "1")
	if err != nil {
		return nil, err
	}
	var p []int

	for rows.Next() {
		var id int

		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		p = append(p, id)
	}

	stmt2 := `SELECT postid, author, title, content
			FROM post
			WHERE postid = $1`
	posts := []*models.Post{}

	for _, v := range p {
		post := &models.Post{}

		row := s.DB.QueryRow(stmt2, v)
		err = row.Scan(&post.PostID, &post.Author, &post.Title, &post.Content)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
