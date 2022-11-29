package edge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	protos "gitlab.com/fl_validator/src/go_protos"
)

type EdgeOperatorServer struct {
}

// LocalTrainFinish : event on finishing local training
func (server *EdgeOperatorServer) LocalTrainFinish(_ context.Context, localTrainResult *protos.LocalTrainResult) (*protos.Empty, error) {
	log.Println(" --- On Local Train Finish --- ", "server", fmt.Sprintf("%v", server))
	return &protos.Empty{}, nil
}

type AppLogData struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

func (server *EdgeOperatorServer) LogMessage(_ context.Context, logMsg *protos.Log) (*protos.Empty, error) {
	log.Println(" --- On LogMessage --- ", "server", fmt.Sprintf("%v", server))

	log.Println("Level: ", logMsg.Level, "Message: ", logMsg.Message)

	logMsg.Level = strings.ToLower(logMsg.Level)

	err := os.MkdirAll(filepath.Dir(os.Getenv("LOG_PATH")), os.ModePerm)
	if err != nil {
		panic(err)
	}

	fo, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	var log AppLogData
	log.Level = strings.ToLower(logMsg.Level)
	log.Timestamp = time.Now().Format(time.RFC3339)
	log.Message = logMsg.Message
	jsonStr, _ := json.Marshal(log)
	fo.Write(jsonStr)
	fo.WriteString("\n")
	fo.Close()

	return &protos.Empty{}, nil

}
