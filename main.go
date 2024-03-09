package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id   int
	Name string
	City string
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "goblog"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
func Index(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	selDB, err := db.Query("SELECT  * FROM Employee  ORDER BY id desc")

	if err != nil {
		panic(err.Error())
	}

	emp := Employee{}

	res := []Employee{}

	for selDB.Next() {

		var id int
		var name, city string

		err := selDB.Scan(&id, &name, &city)

		if err != nil {
			panic(err.Error())
		}

		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}

	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()

}
func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nid := r.URL.Query().Get("id")

	showFrom, err := db.Query("SELECT * FROM Employee where id=?", nid)

	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}

	for showFrom.Next() {

		var id int
		var name, city string

		err := showFrom.Scan(&id, &name, &city)

		if err != nil {
			panic(err.Error())
		}
		emp.Id = id
		emp.City = city
		emp.Name = name

		defer db.Close()
		tmpl.ExecuteTemplate(w, "Show", emp)

	}
	fmt.Println("Show")
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nid := r.URL.Query().Get("id")
	editFrom, err := db.Query("SELECT * FROM Employee  where  id= ?", nid)
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	for editFrom.Next() {
		var id int
		var name, city string
		err := editFrom.Scan(&id, &name, &city)

		if err != nil {
			panic(err.Error())
		}

		emp.Id = id
		emp.City = city
		emp.Name = name
		tmpl.ExecuteTemplate(w, "Edit", emp)

		defer db.Close()
	}
}

func Insert(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")

		insForm, err := db.Prepare("INSERT INTO Employee(name, city) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, city)
		log.Println(name, city)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("uid")

		insForm, err := db.Prepare("UPDATE Employee SET name=?, city=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, city, id)
		defer db.Close()
		http.Redirect(w, r, "/", 301)
	}
	fmt.Println("update")
}

func Delete(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	emp := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id = ?")

	if err != nil {
		panic(err.Error())
	}

	delForm.Exec(emp)

	log.Println("Delete")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {

	log.Println("Server started on: http://localhost:8080")
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}
