package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
	_ "github.com/mattn/go-sqlite3"
)

// inspired from https://github.com/bopbi/simple-todo

func main() {

	db, errOpenDB := sql.Open("sqlite3", "./godb.db")

	checkErr(errOpenDB)
	mainDB = db

	r := pat.New()
	r.Del("/logs/:id", http.HandlerFunc(deleteByID))
	r.Get("/logs/:id", http.HandlerFunc(getByID))
	//r.Put("/logs/:id", http.HandlerFunc(updateByID))
	r.Get("/logs", http.HandlerFunc(getAll))
	//r.Post("/logs", http.HandlerFunc(insert))
	r.Get("/logs/level/:level", http.HandlerFunc(filterByLevel))
	r.Get("/logs/category/:category", http.HandlerFunc(filterByCategory))
	r.Get("/logs/userid/:userid", http.HandlerFunc(filterByUserId))
	http.Handle("/", r)

	err := http.ListenAndServe(":3001", nil)
	if err == nil {
		log.Print(" Running on 3001")
	} else {
		log.Fatal("ListenAndServe: ", err)
	}
}

//testTable (Time,Level,Msg,Category,DebugId,Ip,RequestId,Type,Uri) values (?,?,?,?,?,?,?,?,?)
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
	UserId    int32  `json:"UserId"`
	// a struct with a struct to select the rows in the table and add data.
}

type Logs []Log

var mainDB *sql.DB

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

//Reads all log data from DB
func getAll(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	rows, err := mainDB.Query("SELECT * FROM testTable")
	checkErr(err)
	var logs Logs
	for rows.Next() {
		var log Log
		err = rows.Scan(&log.ID, &log.Time, &log.Level, &log.Msg, &log.Category, &log.DebugId, &log.Ip, &log.RequestId, &log.Type, &log.Uri, &log.UserId)
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
	fmt.Println("GetById----------------------------- " + id)
	stmt, err := mainDB.Prepare(" SELECT * FROM testTable where id = ?")
	checkErr(err)
	rows, errQuery := stmt.Query(id)
	checkErr(errQuery)
	var log Log
	for rows.Next() {
		err = rows.Scan(&log.ID, &log.Time, &log.Level, &log.Msg, &log.Category, &log.DebugId, &log.Ip, &log.RequestId, &log.Type, &log.Uri, &log.UserId)
		checkErr(err)
	}
	jsonB, errMarshal := json.Marshal(log)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

//We dont need this right now!
/*
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
*/
//delete

//Filter: where by category, by from time & to time.
// F,e. time >= {​​from time}​​ and time <={​​to time}​​
//if there is no "to time", only time >= {​​from time}​​
/*

time >= {​​from time}​​ and time <={​​to time}​​
if there is no "to time", only time >= {​​from time}​​

Filter: by type or category
level: info, error, warning
category: system, permisssion, dtabase

*/

// Filters all data by Level f.e. by info, error etc.

func filterByLevel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Filter_By_Level ---------------------------")
	enableCors(&w)
	level := r.URL.Query().Get(":level")
	fmt.Println("Level filter : " + level)
	stmt, err := mainDB.Prepare(
		`SELECT * FROM testTable where level = ?`)
	checkErr(err)
	rows, errQuery := stmt.Query(level)
	checkErr(errQuery)
	var logs Logs

	for rows.Next() {
		var log Log
		err = rows.Scan(&log.ID, &log.Time, &log.Level, &log.Msg, &log.Category, &log.DebugId, &log.Ip, &log.RequestId, &log.Type, &log.Uri, &log.UserId)

		checkErr(err)
		logs = append(logs, log)
	}

	jsonB, errMarshal := json.Marshal(logs)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

// Filters all data by Category

func filterByCategory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Filter_By_Category ---------------------------")
	enableCors(&w)
	category := r.URL.Query().Get(":category")
	fmt.Println("Category filter : " + category)
	stmt, err := mainDB.Prepare(
		`SELECT * FROM testTable where category = ?`)
	checkErr(err)
	rows, errQuery := stmt.Query(category)
	checkErr(errQuery)
	var logs Logs

	for rows.Next() {
		var log Log
		err = rows.Scan(&log.ID, &log.Time, &log.Level, &log.Msg, &log.Category, &log.DebugId, &log.Ip, &log.RequestId, &log.Type, &log.Uri, &log.UserId)

		checkErr(err)
		logs = append(logs, log)
	}

	jsonB, errMarshal := json.Marshal(logs)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
}

// Filters all data by UserId f.e. user id = 0 , 1 ,2 .. int32

func filterByUserId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Filter_By_UserId ---------------------------")
	enableCors(&w)
	userid := r.URL.Query().Get(":userid")
	fmt.Println("Category filter : " + userid)
	stmt, err := mainDB.Prepare(
		`SELECT * FROM testTable where userid = ?`)
	checkErr(err)
	rows, errQuery := stmt.Query(userid)
	checkErr(errQuery)
	var logs Logs

	for rows.Next() {
		var log Log
		err = rows.Scan(&log.ID, &log.Time, &log.Level, &log.Msg, &log.Category, &log.DebugId, &log.Ip, &log.RequestId, &log.Type, &log.Uri, &log.UserId)

		checkErr(err)
		logs = append(logs, log)
	}

	jsonB, errMarshal := json.Marshal(logs)
	checkErr(errMarshal)
	fmt.Fprintf(w, "%s", string(jsonB))
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
