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

//DB CRUD Codes here!
func AddLogToDB(db *sql.DB, Time string, Level string, Msg string) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into testTable (Time,Level,Msg) values (?,?,?)")
	_, err := stmt.Exec(Time, Level, Msg)
	CheckError(err)
	tx.Commit()
}

func GetLogFromDB(db *sql.DB, id2 int) LogToDB {
	rows, err := db.Query("select * from testTable")
	CheckError(err)
	for rows.Next() {
		var tempLogToDB LogToDB
		err =
			rows.Scan(&tempLogToDB.id, &tempLogToDB.Time, &tempLogToDB.Level, &tempLogToDB.Msg)
		CheckError(err)
		if tempLogToDB.id == id2 {
			return tempLogToDB
		}

	}
	return LogToDB{}
}

func GetAllLogFromDB(db *sql.DB) LogToDB {
	rows, err := db.Query("select * from testTable")
	CheckError(err)
	for rows.Next() {
		var tempLogToDB LogToDB
		err =
			rows.Scan(&tempLogToDB.id, &tempLogToDB.Time, &tempLogToDB.Level, &tempLogToDB.Msg)
		CheckError(err)

		return tempLogToDB
	}
	return LogToDB{}
}

func UpdateLogToDB(db *sql.DB, id2 int, Time string, Level string, Msg string) {
	sid := strconv.Itoa(id2) // int to string
	tx, _ := db.Begin()

	stmt, _ := tx.Prepare("update testTable set Time=?,Level=?,Msg=? where id=?")
	_, err := stmt.Exec(Time, Level, Msg, sid)
	CheckError(err)
	tx.Commit()
}

func DeleteLogToDB(db *sql.DB, id2 int) {
	sid := strconv.Itoa(id2) // int to string
	tx, _ := db.Begin()

	stmt, _ := tx.Prepare("delete from testTable where id=?")
	_, err := stmt.Exec(sid)
	CheckError(err)
	tx.Commit()
}

type LogToDB struct {
	id    int
	Time  string
	Level string
	Msg   string

	//I created a struct with a struct to select the rows in the table and add data.
}

func CheckError(err error) {
	if err != nil {

		fmt.Print("ERROR FROM DBMANAGER!")
		panic(err)
	}

	// catch to error.

}
