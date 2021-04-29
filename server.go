package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	//"github.com/grpc-digimakergo/log-grpc/logpb"
	"logpb"

	//"github.com/digimakergo/log-grpc/logpb"

	"google.golang.org/grpc"

	//for DB connection
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type server struct{}

func (*server) SendLogs(stream logpb.LogService_SendLogsServer) error {
	db, _ := sql.Open("sqlite3", "./godb.db")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//Finished reading client stream
			return stream.SendAndClose(&logpb.DummyResult{
				Success: true,
				Error:   "",
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		// Convert stream data to String and read as JSON

		result, _ := json.Marshal(req)
		type LogJSON struct {
			Time  string
			Level string
			Msg   string
			//and the others! // TODO
		}

		type LogMain struct {
			Logs LogJSON
		}

		str := string(result)
		var ourLogs LogMain

		fmt.Println("Only str: ", str)

		json.Unmarshal([]byte(str), &ourLogs)

		fmt.Println("Only res: ", ourLogs)
		fmt.Println("-----------------------------------------------------------")
		fmt.Println("Res.Time: ", ourLogs.Logs.Time)
		fmt.Println("-----------------------------------------------------------")
		fmt.Println("Res.Level: ", ourLogs.Logs.Level)
		fmt.Println("-----------------------------------------------------------")
		fmt.Println("Res.Msg: ", ourLogs.Logs.Msg)

		addLogToDB(db, ourLogs.Logs.Time, ourLogs.Logs.Level, ourLogs.Logs.Msg)

		fmt.Println("-----------------------------------------------------------")
		fmt.Println("-----------------------------------------------------------")

	}
}

//DB CRUD Codes here!

func main() {

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

	//Port listening here!

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	logpb.RegisterLogServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

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
