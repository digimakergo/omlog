package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	//	dbmanager "github.com/digimakergo/omlog/dbmanager"

	//"github.com/grpc-digimakergo/log-grpc/logpb"
	"logpb"

	//"github.com/digimakergo/log-grpc/logpb"

	"google.golang.org/grpc"

	//for DB connection

	_ "github.com/mattn/go-sqlite3"
)

type server struct{}

func (*server) SendLogs(stream logpb.LogService_SendLogsServer) error {
	db, _ := sql.Open("sqlite3", "./httpconnection/godb.db")
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

		AddLogToDB(db, ourLogs.Logs.Time, ourLogs.Logs.Level, ourLogs.Logs.Msg)

		fmt.Println("-----------------------------------------------------------")
		fmt.Println("-----------------------------------------------------------")

	}
}

//DB CRUD Codes here!

func main() {

	//DB main func Codes
	db, err := sql.Open("sqlite3", "./httpconnection/godb.db")

	fmt.Print("Creates new database!")
	db.Exec(`
		CREATE TABLE IF NOT EXISTS "testTable" (
			"id"	INTEGER UNIQUE,
			"Time"	text,
			"Level"	text,
			"Msg"	text,
			PRIMARY KEY("id" AUTOINCREMENT)
		);
		
		`)

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

func CheckError(err error) {
	if err != nil {
		fmt.Print("ERROR IN SERVER.GO !")
		panic(err)
	}

	// catch to error.
}

func AddLogToDB(db *sql.DB, Time string, Level string, Msg string) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into testTable (Time,Level,Msg) values (?,?,?)")
	_, err := stmt.Exec(Time, Level, Msg)
	CheckError(err)
	tx.Commit()
}
