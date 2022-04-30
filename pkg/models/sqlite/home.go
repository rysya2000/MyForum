package sqlite

import (
	"git.01.alem.school/rysya2000/MyForum.git/pkg/models"
)

func (s *ForumModel) GetAllPosts() ([]*models.Post, error) {
	stmt := `SELECT postid, author, title, content FROM post ORDER BY "created" DESC`

	rows, err := s.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		post := &models.Post{}

		err := rows.Scan(&post.PostID, &post.Author, &post.Title, &post.Content)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *ForumModel) GetPostsWithTag(tag string) ([]*models.Post, error) {
	stmt := `SELECT postid
			FROM tag
			WHERE tag = $1`
	rows, err := s.DB.Query(stmt, tag)
	if err != nil {
		return nil, err
	}
	// `SELECT p.*
	// FROM tags t
	// JOIN posts p ON t.postid = p.postid
	// WHERE t.tag = $1`

	var res []int

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		res = append(res, id)
	}

	stmt = `SELECT postid, author, title, content
			FROM post
			WHERE postid = $1
			ORDER BY "created" DESC`

	posts := []*models.Post{}

	for _, v := range res {
		post := &models.Post{}

		row := s.DB.QueryRow(stmt, v)
		err = row.Scan(&post.PostID, &post.Author, &post.Title, &post.Content)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
