package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/digimakergo/omlog/log-grpc/logpb"

	//"github.com/digimakergo/log-grpc/logpb"

	"google.golang.org/grpc"

	//for DB connection

	_ "github.com/mattn/go-sqlite3"
)

type server struct{}

type LogJSON struct {
	Time      string
	Level     string
	Msg       string
	Category  string
	DebugId   string
	Ip        string
	RequestId string
	Type      string
	Uri       string
}

type LogMain struct {
	Logs LogJSON
}

var (
	websocketConnections = make(map[*websocket.Conn]time.Time)
	upgrader             = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	con, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		log.Println(err)
		return
	}

	websocketConnections[con] = time.Now()
	con.SetReadDeadline(time.Now().Add(pongWait))
	con.SetPongHandler(func(string) error { con.SetReadDeadline(time.Now().Add(pongWait)); return nil })
}

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
		fmt.Println("-----------------------------------------------------------")
		fmt.Println("Res.Category: ", ourLogs.Logs.Category)
		fmt.Println("-----------------------------------------------------------")
		fmt.Println("Res.DebugId: ", ourLogs.Logs.DebugId)
		fmt.Println("-----------------------------------------------------------")
		fmt.Println("Res.Ip: ", ourLogs.Logs.Ip)
		fmt.Println("-----------------------------------------------------------")
		fmt.Println("Res.RequestId: ", ourLogs.Logs.RequestId)
		fmt.Println("-----------------------------------------------------------")
		fmt.Println("Res.Type: ", ourLogs.Logs.Type)
		fmt.Println("-----------------------------------------------------------")
		fmt.Println("Res.Uri: ", ourLogs.Logs.Uri)

		if ourLogs.Logs.Level != "debug" {
			AddLogToDB(db, ourLogs.Logs.Time, ourLogs.Logs.Level, ourLogs.Logs.Msg, ourLogs.Logs.Category, ourLogs.Logs.DebugId, ourLogs.Logs.Ip, ourLogs.Logs.RequestId, ourLogs.Logs.Type, ourLogs.Logs.Uri)
		} else {
			go sendLogsToWebsocketConnections(ourLogs)
		}

		fmt.Println("-----------------------------------------------------------")
		fmt.Println("-----------------------------------------------------------")
	}
}

func sendLogsToWebsocketConnections(ourLogs LogMain) {
	fmt.Println("IN SENDLOGS")
	for websocketConnection, _ := range websocketConnections {
		websocketConnection.SetWriteDeadline(time.Now().Add(writeWait))
		if err := websocketConnection.WriteJSON(ourLogs); err != nil {
			fmt.Println("Error sending to via websocket:  ", err)
			delete(websocketConnections, websocketConnection)
		}
	}
}

func setuproutes() {
	http.HandleFunc("/ws/debug-logs", websocketHandler)
}

//DB CRUD Codes here!

func main() {

	fmt.Printf("Server started at: %v", time.Now())
	setuproutes()
	go http.ListenAndServe(":6001", nil)

	//DB main func Codes
	db, _ := sql.Open("sqlite3", "./godb.db")
	db.Exec(`
	CREATE TABLE IF NOT EXISTS "testTable" (
		"id"	INTEGER UNIQUE,
		"Time"	text,
		"Level"	text,
		"Msg"	text,
		"Category" text,
		"DebugId" text,
		"Ip" text,
		"RequestId" text,
		"Type" text,
		"Uri" text,
		
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

type LogToDB struct {
	Time      string
	Level     string
	Msg       string
	Category  string
	DebugId   string
	Ip        string
	RequestId string
	Type      string
	Uri       string

	//I created a struct with a struct to select the rows in the table and add data.
}

func CheckError(err error) {
	if err != nil {
		fmt.Print("ERROR IN SERVER.GO !")
		panic(err)
	}

	// catch to error.
}

func AddLogToDB(db *sql.DB, Time string, Level string, Msg string, Category string, DebugId string, Ip string, RequestId string, Type string, Uri string) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into testTable (Time,Level,Msg,Category,DebugId,Ip,RequestId,Type,Uri) values (?,?,?,?,?,?,?,?,?)")
	_, err := stmt.Exec(Time, Level, Msg, Category, DebugId, Ip, RequestId, Type, Uri)
	CheckError(err)
	tx.Commit()
}
