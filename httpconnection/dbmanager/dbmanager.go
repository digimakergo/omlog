package dbmanager

import (
	"database/sql"
	"fmt"
	"strconv"

	//"github.com/grpc-digimakergo/log-grpc/logpb"

	//"github.com/digimakergo/log-grpc/logpb"

	//for DB connection

	_ "github.com/mattn/go-sqlite3"
)

func db() {

	//DB main func Codes
	db, _ := sql.Open("sqlite3", "./godb.db")
	db.Exec(`

		CREATE TABLE IF NOT EXISTS "testTable" (
			"id"	INTEGER UNIQUE,
			"Time"	text,
			"Level"	text,
			"Msg"	text,
			PRIMARY KEY("id" AUTOINCREMENT)
		);

	`)

	addLogToDB(db, "T TIME ", "T LEVEL", "TEST MSG") // added data to database

	updateLogToDB(db, 2, "U TIME", "U LEVEL", "U MSG") //update data to database

	deleteLogToDB(db, 1) // delete data to database

	fmt.Println(getLogFromDB(db, 2)) // printing the Log

}

//DB CRUD Codes here!
func addLogToDB(db *sql.DB, Time string, Level string, Msg string) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into testTable (Time,Level,Msg) values (?,?,?)")
	_, err := stmt.Exec(Time, Level, Msg)
	checkError(err)
	tx.Commit()
}

func getLogFromDB(db *sql.DB, id2 int) LogToDB {
	rows, err := db.Query("select * from testTable")
	checkError(err)
	for rows.Next() {
		var tempLogToDB LogToDB
		err =
			rows.Scan(&tempLogToDB.id, &tempLogToDB.Time, &tempLogToDB.Level, &tempLogToDB.Msg)
		checkError(err)
		if tempLogToDB.id == id2 {
			return tempLogToDB
		}

	}
	return LogToDB{}
}

func getAllLogFromDB(db *sql.DB) LogToDB {
	rows, err := db.Query("select * from testTable")
	checkError(err)
	for rows.Next() {
		var tempLogToDB LogToDB
		err =
			rows.Scan(&tempLogToDB.id, &tempLogToDB.Time, &tempLogToDB.Level, &tempLogToDB.Msg)
		checkError(err)

		return tempLogToDB
	}
	return LogToDB{}
}

func updateLogToDB(db *sql.DB, id2 int, Time string, Level string, Msg string) {
	sid := strconv.Itoa(id2) // int to string
	tx, _ := db.Begin()

	stmt, _ := tx.Prepare("update testTable set Time=?,Level=?,Msg=? where id=?")
	_, err := stmt.Exec(Time, Level, Msg, sid)
	checkError(err)
	tx.Commit()
}

func deleteLogToDB(db *sql.DB, id2 int) {
	sid := strconv.Itoa(id2) // int to string
	tx, _ := db.Begin()

	stmt, _ := tx.Prepare("delete from testTable where id=?")
	_, err := stmt.Exec(sid)
	checkError(err)
	tx.Commit()
}

type LogToDB struct {
	id    int
	Time  string
	Level string
	Msg   string

	//I created a struct with a struct to select the rows in the table and add data.
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}

	// catch to error.

}
