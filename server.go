package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func main() {

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	database = db
	defer db.Close()

	http.HandleFunc("/", PageHandler)
	http.HandleFunc("/create", CreateHandler)
	http.HandleFunc("/delete", DeleteHandler)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8082", nil)
}

type User struct {
	Id    int
	Name1 string
	Name2 string
}

func PageHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("select  * from users1;")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	user_current := []User{}
	for rows.Next() {
		p := User{}
		err := rows.Scan(&p.Id, &p.Name1, &p.Name2)
		if err != nil {
			fmt.Println(err)
			continue
		}
		user_current = append(user_current, p)
	}
	tmpl, _ := template.ParseFiles("templates/index.html")
	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, user_current)
	fmt.Println(user_current)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		checkError(err)
		Name1 := r.FormValue("Name1")
		Name2 := r.FormValue("Name2")

		_, err = database.Exec("insert into users1 (Name1, Name2) values (?, ?)",
			Name1, Name2)

		checkError(err)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		http.ServeFile(w, r, "templates/create.html")
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	_, err := database.Exec("delete from users1 where Id = ?", id)
	checkError(err)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
