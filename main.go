package main

import (
	"database/sql"
	"fmt"
	"forum/database"
	"forum/handlers"
	"forum/handlers/pages"
	"forum/handlers/user"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Creer une nouvelle db
	database.StartDB(false)

	// Ouvrir une db existante
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Creer la table user si elle n'existe pas
	err = database.CreateTables(db)
	if err != nil {
		log.Fatal(err)
	}
	database.GetTags(db)

	// Page handlers
	http.HandleFunc("/", pages.HomeHandler(db))
	http.HandleFunc("/createPost", pages.CreatePostHandler(db))
	http.HandleFunc("/viewPost", pages.ViewPostHandler(db))
	http.HandleFunc("/reply", pages.ReplyHandler(db))

	http.HandleFunc("/register", user.RegisterHandler(db))
	http.HandleFunc("/login", user.LoginHandler(db))
	http.HandleFunc("/logout", user.LogoutHandler)
	http.HandleFunc("/oauth/", user.OauthHandler(db))

	http.HandleFunc("/activity/", pages.ActivityHandler(db))
	http.HandleFunc("/notifications", pages.NotifyHandler(db))
	http.HandleFunc("/edit/", pages.EditHandler(db))
	http.HandleFunc("/delete/", pages.DeleteHandler(db))
	http.HandleFunc("/user/", user.UserHandler(db))
	http.HandleFunc("/report/", pages.ReportHandler(db))
	http.HandleFunc("/admin/", pages.AdminHandler(db))

	// Lien des assets
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/favicon.ico", ignoreFavicon)

	fmt.Println("Go to: http://localhost:8181")
	handlers.Open("http://localhost:8181/")
	log.Fatal(http.ListenAndServe(":8181", nil))
}

func ignoreFavicon(_ http.ResponseWriter, _ *http.Request) {}
