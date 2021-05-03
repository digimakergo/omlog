package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/bmizerany/pat"
	_ "github.com/mattn/go-sqlite3"
)

// inspired from https://github.com/bopbi/simple-todo

func main() {

	db, errOpenDB := sql.Open("sqlite3", "godb.db")
	checkErr(errOpenDB)
	mainDB = db

	r := pat.New()
	r.Del("/logs/:id", http.HandlerFunc(deleteByID))
	r.Get("/logs/:id", http.HandlerFunc(getByID))
	r.Put("/logs/:id", http.HandlerFunc(updateByID))
	r.Get("/logs", http.HandlerFunc(getAll))
	r.Post("/logs", http.HandlerFunc(insert))

	http.Handle("/", r)

	log.Print(" Running on 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type Log struct {
	ID    int64  `json:"id"`
	Time  string `json:"time"`
	Level string `json:"level"`
	Msg   string `json:"msg"`

	//I created a struct with a struct to select the rows in the table and add data.
}

type Logs []Log

var mainDB *sql.DB

func getAll(w http.ResponseWriter, r *http.Request) {
	rows, err := mainDB.Query("SELECT * FROM testTable")
	checkErr(err)
	var logs Logs
	for rows.Next() {
		var log Log
		err = rows.Scan(&log.ID, &log.Time, &log.Level, &log.Msg)
		checkErr(err)
		logs = append(logs, log)
	}
	jsonB, errMarshal := json.Marshal(logs)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

func getByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	stmt, err := mainDB.Prepare(" SELECT * FROM testTable where id = ?")
	checkErr(err)
	rows, errQuery := stmt.Query(id)
	checkErr(errQuery)
	var log Log
	for rows.Next() {
		err = rows.Scan(&log.ID, &log.Time, &log.Level, &log.Msg)
		checkErr(err)
	}
	jsonB, errMarshal := json.Marshal(log)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

func insert(w http.ResponseWriter, r *http.Request) {
	time := r.FormValue("time")
	level := r.FormValue("level")
	msg := r.FormValue("msg")

	var log Log
	log.Time = time
	log.Level = level
	log.Msg = msg

	stmt, err := mainDB.Prepare("INSERT INTO testTable(time, level, msg) values (?, ?, ?)")
	checkErr(err)
	result, errExec := stmt.Exec(log.Time, log.Level, log.Msg)
	checkErr(errExec)
	newID, errLast := result.LastInsertId()
	checkErr(errLast)
	log.ID = newID
	jsonB, errMarshal := json.Marshal(log)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

func updateByID(w http.ResponseWriter, r *http.Request) {
	time := r.FormValue("time")
	level := r.FormValue("level")
	msg := r.FormValue("msg")

	id := r.URL.Query().Get(":id")
	var log Log
	ID, _ := strconv.ParseInt(id, 10, 0)
	log.ID = ID
	log.Time = time
	log.Level = level
	log.Msg = msg

	stmt, err := mainDB.Prepare("UPDATE testTable SET time = ? SET level = ? SET msg = ?  WHERE id = ?")
	checkErr(err)
	result, errExec := stmt.Exec(log.Time, log.Level, log.Msg, log.ID)
	checkErr(errExec)
	rowAffected, errLast := result.RowsAffected()
	checkErr(errLast)
	if rowAffected > 0 {
		jsonB, errMarshal := json.Marshal(log)
		checkErr(errMarshal)
		fmt.Fprintf(w, "%s", string(jsonB))
	} else {
		fmt.Fprintf(w, "{row_affected=%d}", rowAffected)
	}

}

func deleteByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	stmt, err := mainDB.Prepare("DELETE FROM testTable WHERE id = ?")
	checkErr(err)
	result, errExec := stmt.Exec(id)
	checkErr(errExec)
	rowAffected, errRow := result.RowsAffected()
	checkErr(errRow)
	fmt.Fprintf(w, "{row_affected=%d}", rowAffected)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
