package internal

import (
	"html/template"
	"log"

	"git.01.alem.school/rysya2000/MyForum.git/pkg/models/sqlite"
)

type App struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Forum         *sqlite.ForumModel
	TemplateCache *template.Template
}
