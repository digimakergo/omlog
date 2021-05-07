//Author xc, Created on 2020-04-09 19:00
//{COPYRIGHTS}
package log

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digimakergo/digimaker/core/log/log-grpc/logpb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type RemoteHook struct {
}

type Client struct{}

func (hook *RemoteHook) Fire(entry *log.Entry) error {
	/*
		_, err := entry.Bytes()
		if err != nil {
			return err
		}
	*/

	fmt.Println("In Fire")
	client := Client{}
	client.SendWithgRPC(entry)
	//maybe create goroutine here for sending with grpc?

	//todo: based on settings(eg. debug by ip/user), output context log information.
	/*f, err := os.OpenFile("request-debug.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	f.Write(line)
	defer f.Close()
	*/
	return nil
}

// Levels define on which log levels this hook would trigger
func (hook *RemoteHook) Levels() []log.Level {
	return log.AllLevels
}

func SendWS(entry *log.Entry) []byte {

	logFields := entry.Data
	category := ""

	if logFields["category"] != nil {
		category = logFields["category"].(string)
	}

	logEntry := &logpb.Log{
		Time:      entry.Time.String(),
		Level:     entry.Level.String(),
		Msg:       entry.Message,
		Category:  category,
		DebugId:   logFields["debug_id"].(string),
		Ip:        logFields["ip"].(string),
		RequestId: logFields["request_id"].(string),
		Type:      logFields["type"].(string),
		Uri:       logFields["uri"].(string),
		Id:        1,
	}

	logs := []*logpb.Log{logEntry}
	logarray := []*logpb.SendLogsRequest{
		&logpb.SendLogsRequest{
			Logs: logs,
		},
	}

	result, _ := json.Marshal(logarray)
	return result
}

func (*Client) SendWithgRPC(entry *log.Entry) error {

	//Connection should be established only once
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()

	c := logpb.NewLogServiceClient(cc)

	stream, err := c.SendLogs(context.Background())
	if err != nil {
		log.Fatalf("Error while calling SendLogs: %v", err)
	}

	logFields := entry.Data
	category := ""

	if logFields["category"] != nil {
		category = logFields["category"].(string)
	}

	logEntry := &logpb.Log{
		Time:      entry.Time.String(),
		Level:     entry.Level.String(),
		Msg:       entry.Message,
		Category:  category,
		DebugId:   logFields["debug_id"].(string),
		Ip:        logFields["ip"].(string),
		RequestId: logFields["request_id"].(string),
		Type:      logFields["type"].(string),
		Uri:       logFields["uri"].(string),
		Id:        1,
	}

	logs := []*logpb.Log{logEntry}
	logarray := []*logpb.SendLogsRequest{
		&logpb.SendLogsRequest{
			Logs: logs,
		},
	}

	fmt.Print(logarray)

	for _, req := range logarray {
		stream.Send(req)
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("Error with response: %v", err)
	}

	return nil
}
