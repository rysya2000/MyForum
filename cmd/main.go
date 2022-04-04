package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"git.01.alem.school/rysya2000/MyForum.git/internal"
	"git.01.alem.school/rysya2000/MyForum.git/pkg/models"
	"git.01.alem.school/rysya2000/MyForum.git/pkg/models/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	dsn := flag.String("dsn", "forum.db", "database")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := models.OpenDB(*dsn)
	if err != nil {
		log.Printf("%v\n", err)
		errorLog.Fatal(err)
	}

	tmpCache, err := internal.NewTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &internal.App{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		Forum:         &sqlite.ForumModel{DB: db},
		TemplateCache: tmpCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}

	infoLog.Printf("Server: http://localhost%v", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
