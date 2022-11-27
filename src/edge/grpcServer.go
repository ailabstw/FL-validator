package edge

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	protos "gitlab.com/fl_validator/src/go_protos"
	"go.uber.org/zap"
)

type EdgeOperatorServer struct {
}

// type AbstractOperator interface {
// 	GrpcServerRegister(*grpc.Server)
// 	GetPayload() interface{}
// 	TrainFinish()
// }

// LocalTrainFinish : event on finishing local training
func (server *EdgeOperatorServer) LocalTrainFinish(_ context.Context, localTrainResult *protos.LocalTrainResult) (*protos.Empty, error) {
	zap.L().Debug(" --- On Local Train Finish --- ", zap.String("server", fmt.Sprintf("%v", server)))
	zap.L().Debug(fmt.Sprintf("Receive localTrainResult.Metadata [%v]", localTrainResult.Metadata))
	zap.L().Debug(fmt.Sprintf("Receive localTrainResult.Metrics [%v]", localTrainResult.Metrics))

	// var metadata map[string]string
	// if localTrainResult.Metadata == nil {
	// 	metadata = map[string]string{}
	// } else {
	// 	metadata = localTrainResult.Metadata
	// }

	// var metrics map[string]float64
	// if localTrainResult.Metrics == nil {
	// 	metrics = map[string]float64{}
	// } else {
	// 	metrics = localTrainResult.Metrics
	// }

	// server.operator.Dispatch(&trainFinishAction{
	// 	errCode:     int(localTrainResult.Error),
	// 	datasetSize: int(localTrainResult.DatasetSize),
	// 	metadata:    metadata,
	// 	metrics:     metrics,
	// })

	return &protos.Empty{}, nil
}

func (server *EdgeOperatorServer) LogMessage(_ context.Context, logMsg *protos.Log) (*protos.Empty, error) {
	zap.L().Debug(" --- On LogMessage --- ", zap.String("server", fmt.Sprintf("%v", server)))
	zap.L().Debug(fmt.Sprintf("Receive logMsg.Level [%v]", logMsg.Level))
	zap.L().Debug(fmt.Sprintf("Receive logMsg.Message [%v]", logMsg.Message))

	err := os.MkdirAll("/ver/ailabs/", os.ModePerm)
	if err != nil {
		panic(err)
	}

	fo, err := os.OpenFile("/ver/ailabs/harmonia.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}

	fo.Write([]byte(time.Now().Format(time.RFC3339)))
	fo.Write([]byte(" "))
	fo.Write([]byte(strings.ToUpper(logMsg.Level)))
	fo.Write([]byte(" "))
	fo.Write([]byte(logMsg.Message))
	fo.WriteString("\n")

	fo.Close()

	return &protos.Empty{}, nil
}
