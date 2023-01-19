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
	"gitlab.com/fl_validator/src/util"
)

type EdgeOperatorServer struct {
}

// LocalTrainFinish : event on finishing local training
func (server *EdgeOperatorServer) LocalTrainFinish(_ context.Context, localTrainResult *protos.LocalTrainResult) (*protos.Empty, error) {
	log.Println(" --- On Local Train Finish --- ", "server", fmt.Sprintf("%v", server))
	if localTrainResult.Metadata != nil {
		if localTrainResult.Metadata.DatasetSize != 0 {
			util.WriteReport("LocalTrain", "LocalTrainFinish:Metadata.DatasetSize correct", "")
		} else {
			util.WriteReport("LocalTrain", "LocalTrainFinish:Metadata.DatasetSize not correct", "LocalTrainFinish:Metadata.DatasetSize not correct")
			util.MakeResultFalse()
		}

		if localTrainResult.Metadata.Importance != 0 {
			util.WriteReport("LocalTrain", "LocalTrainFinish:Metadata.Importance correct", "")
		} else {
			util.WriteReport("LocalTrain", "LocalTrainFinish:Metadata.Importance not correct", "LocalTrainFinish:Metadata.Importance not correct")
			util.MakeResultFalse()
		}
	} else {
		util.WriteReport("LocalTrain", "LocalTrainFinish:Metadata not correct", "LocalTrainFinish:Metadata not correct")
		util.MakeResultFalse()
	}

	if localTrainResult.Metrics != nil {
		s := map[string]bool{"basic/confusion_tn": true, "basic/confusion_fp": true, "basic/confusion_fn": true, "basic/confusion_tp": true}
		for k, _ := range s {
			_, ok := localTrainResult.Metrics[k]
			if !ok {
				util.WriteReport("LocalTrain", "LocalTrainFinish:Metrics"+k+"not found", "LocalTrainFinish:Metrics"+k+"not found")
			} else {
				util.WriteReport("LocalTrain", "LocalTrainFinish:Metrics"+k+"found and correct", "")
				util.MakeResultFalse()
			}
		}
	} else {
		util.WriteReport("LocalTrain", "LocalTrainFinish:Metrics not correct", "LocalTrainFinish:Metrics not correct")
		util.MakeResultFalse()
	}

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
