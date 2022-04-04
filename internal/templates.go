package internal

import (
	"html/template"
	"path/filepath"

	"git.01.alem.school/rysya2000/MyForum.git/pkg/models"
)

type TemplateData struct {
	User       *models.User
	IsAuth     bool
	Post       *models.Post
	Posts      []*models.Post
	LikedPosts []*models.Post
	Err        Error
}

type Error struct {
	IsError bool
	Msg     string
}

var GTemplate *template.Template

func NewTemplateCache(dir string) (*template.Template, error) {
	var cache *template.Template

	pages, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		return nil, err
	}
	// fmt.Println(pages)
	files, err := template.ParseFiles(pages...)
	if err != nil {
		return nil, err
	}

	cache = template.Must(files, nil)

	return cache, nil
}
