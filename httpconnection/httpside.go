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

//estTable (Time,Level,Msg,Category,DebugId,Ip,RequestId,Type,Uri) values (?,?,?,?,?,?,?,?,?)
type Log struct {
	ID        int64  `json:"id"`
	Time      string `json:"time"`
	Level     string `json:"level"`
	Msg       string `json:"msg"`
	Category  string `json:"category"`
	DebugId   string `json:"debugId"`
	Ip        string `json:"Ip"`
	RequestId string `json:"RequestId"`
	Type      string `json:"Type"`
	Uri       string `json:"Uri"`
	//I created a struct with a struct to select the rows in the table and add data.
}

type Logs []Log

var mainDB *sql.DB

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getAll(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	rows, err := mainDB.Query("SELECT * FROM testTable")
	checkErr(err)
	var logs Logs
	for rows.Next() {
		var log Log
		err = rows.Scan(&log.ID, &log.Time, &log.Level, &log.Msg, &log.Category, &log.DebugId, &log.Ip, &log.RequestId, &log.Type, &log.Uri)
		checkErr(err)
		logs = append(logs, log)
	}
	jsonB, errMarshal := json.Marshal(logs)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

func getByID(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	id := r.URL.Query().Get(":id")
	stmt, err := mainDB.Prepare(" SELECT * FROM testTable where id = ?")
	checkErr(err)
	rows, errQuery := stmt.Query(id)
	checkErr(errQuery)
	var log Log
	for rows.Next() {
		err = rows.Scan(&log.ID, &log.Time, &log.Level, &log.Msg, &log.Category, &log.DebugId, &log.Ip, &log.RequestId, &log.Type, &log.Uri)
		checkErr(err)
	}
	jsonB, errMarshal := json.Marshal(log)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

func insert(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	time := r.FormValue("time")
	level := r.FormValue("level")
	msg := r.FormValue("msg")
	category := r.FormValue("category")
	debugid := r.FormValue("debugid")
	ip := r.FormValue("ip")
	requestid := r.FormValue("requestid")
	Type := r.FormValue("type")
	uri := r.FormValue("uri")

	var log Log
	log.Time = time
	log.Level = level
	log.Msg = msg
	log.Category = category
	log.DebugId = debugid
	log.Ip = ip
	log.RequestId = requestid
	log.Type = Type
	log.Uri = uri

	stmt, err := mainDB.Prepare("INSERT INTO testTable (Time,Level,Msg,Category,DebugId,Ip,RequestId,Type,Uri) values (?,?,?,?,?,?,?,?,?)")
	checkErr(err)
	result, errExec := stmt.Exec(log.Time, log.Level, log.Msg, log.Category, log.DebugId, log.Ip, log.RequestId, log.Type, log.Uri)
	checkErr(errExec)
	newID, errLast := result.LastInsertId()
	checkErr(errLast)
	log.ID = newID
	jsonB, errMarshal := json.Marshal(log)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

func updateByID(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	time := r.FormValue("time")
	level := r.FormValue("level")
	msg := r.FormValue("msg")
	category := r.FormValue("category")
	debugid := r.FormValue("debugid")
	ip := r.FormValue("ip")
	requestid := r.FormValue("requestid")
	Type := r.FormValue("type")
	uri := r.FormValue("uri")

	id := r.URL.Query().Get(":id")
	var log Log
	ID, _ := strconv.ParseInt(id, 10, 0)
	log.ID = ID
	log.Time = time
	log.Level = level
	log.Msg = msg
	log.Category = category
	log.DebugId = debugid
	log.Ip = ip
	log.RequestId = requestid
	log.Type = Type
	log.Uri = uri

	stmt, err := mainDB.Prepare("UPDATE testTable SET time = ? SET level = ? SET msg = ? SET category = ? SET debugid = ? SET ip = ? SET requestid = ? SET Type = ? SET uri = ? WHERE id = ?")
	checkErr(err)
	result, errExec := stmt.Exec(log.Time, log.Level, log.Msg, log.Category, log.DebugId, log.Ip, log.RequestId, log.Type, log.Uri, log.ID)
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
	enableCors(&w)
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
