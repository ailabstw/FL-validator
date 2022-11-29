package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"gitlab.com/fl_validator/src/edge"
	protos "gitlab.com/fl_validator/src/go_protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var mutx *sync.Mutex

type baseModel struct {
	repoName string
	metadata map[string]string
	metrics  map[string]float64
}

func checkOnlyInterface(appGrpcServerURI string) {

	allImplement := true
	allImplement = allImplement && sendIsValidated(true, appGrpcServerURI)
	allImplement = allImplement && sendInitMessage(true, appGrpcServerURI)
	allImplement = allImplement && sendLocalTrainMessage(true, appGrpcServerURI, 1, baseModel{}, "")
	allImplement = allImplement && sendTrainFinishMessage(true, appGrpcServerURI)
	allImplement = allImplement && sendTrainInteruptMessage(true, appGrpcServerURI)
	if allImplement {
		WriteReport("check", "all interface implemented", "")
	}
}

func main() {
	mutx = new(sync.Mutex)
	clientURI := os.Getenv("APP_URI")
	serverURI := "0.0.0.0:8787"

	log.Println("clientURI: " + clientURI)

	isInterfaceOnly, _ := strconv.ParseBool(os.Getenv("IS_INTERFACE_ONLY"))

	log.Println("Starting grpc server.... ")

	startGrpcServer(serverURI)

	if isInterfaceOnly {
		log.Println("Check only interface .... ")
		checkOnlyInterface(serverURI)
		return
	}

	time.Sleep(10 * time.Second)

	log.Println("Sending  initialization msg.... ")

	sendIsValidated(isInterfaceOnly, serverURI)

	time.Sleep(10 * time.Second)

	log.Println("Sending  initialization msg.... ")

	sendInitMessage(isInterfaceOnly, clientURI)

	time.Sleep(10 * time.Second)

	log.Println("Sending local train message .... ")
	// test localtrain
	sendLocalTrainMessage(isInterfaceOnly, clientURI, 1, baseModel{}, "")

	time.Sleep(10 * time.Second)

	log.Println("Sending training finished message .... ")
	// test train finish
	sendTrainFinishMessage(isInterfaceOnly, clientURI)

	// all test sucessfully
	log.Println("All FL validation completed . Congrats. ")

}

func startGrpcServer(address string) *grpc.Server {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Cannot listen on the address", "address", address)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	go func() {

		log.Println("Grpc server listen the address .... ", "address", address)

		protos.RegisterEdgeOperatorServer(grpcServer, &edge.EdgeOperatorServer{})

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Cannot start server on the address", "address", address)
			log.Fatal("Cannot start server on the address")
		}
	}()

	log.Println(fmt.Sprintf("Grpc Listen [%v]", "0.0.0.0:8787"))

	return grpcServer
}

func sendIsValidated(isInterfaceOnly bool, appGrpcServerURI string) bool {
	return EmitEvent(
		isInterfaceOnly,
		"IsValidated",
		appGrpcServerURI,
		func(conn *grpc.ClientConn) interface{} {
			return protos.NewEdgeAppClient(conn)
		},
		func(ctx context.Context, client interface{}) (interface{}, error) {
			return client.(protos.EdgeAppClient).IsDataValidated(ctx, &protos.Empty{})
		},
		func(response interface{}) interface{} {
			if isInterfaceOnly {
				WriteReport("IsValidated", "implemented", "")
			} else {
				WriteReport("IsValidated", "IsValidated sucessfully.", "")
			}
			return nil
		},
	)
}

func sendInitMessage(isInterfaceOnly bool, appGrpcServerURI string) bool {
	return EmitEvent(
		isInterfaceOnly,
		"TrainInit",
		appGrpcServerURI,
		func(conn *grpc.ClientConn) interface{} {
			return protos.NewEdgeAppClient(conn)
		},
		func(ctx context.Context, client interface{}) (interface{}, error) {
			return client.(protos.EdgeAppClient).TrainInit(ctx, &protos.Empty{})
		},
		func(response interface{}) interface{} {
			if isInterfaceOnly {
				WriteReport("TrainInit", "implemented", "")
			} else {
				WriteReport("TrainInit", "TrainInit sucessfully.", "")
			}
			return nil
		},
	)
}

func sendLocalTrainMessage(isInterfaceOnly bool, appGrpcServerURI string, epochPerRound int, baseModel baseModel, edgeRepoName string) bool {
	return EmitEvent(
		isInterfaceOnly,
		"LocalTrain",
		appGrpcServerURI,
		func(conn *grpc.ClientConn) interface{} {
			return protos.NewEdgeAppClient(conn)
		},
		func(ctx context.Context, client interface{}) (interface{}, error) {
			return client.(protos.EdgeAppClient).LocalTrain(ctx, &protos.LocalTrainParams{
				BaseModel: &protos.LocalTrainParams_BaseModel{
					Path:     baseModel.repoName,
					Metadata: baseModel.metadata,
					Metrics:  baseModel.metrics,
				},
				LocalModel: &protos.LocalTrainParams_LocalModel{
					Path: edgeRepoName,
				},
				EpR: int32(epochPerRound),
			})
		},
		func(response interface{}) interface{} {
			if isInterfaceOnly {
				WriteReport("LocalTrain", "implemented", "")
			} else {
				WriteReport("LocalTrain", "LocalTrain sucessfully.", "")
			}
			return nil
		},
	)
}

func sendTrainFinishMessage(isInterfaceOnly bool, appGrpcServerURI string) bool {
	return EmitEvent(
		isInterfaceOnly,
		"TrainFinish",
		appGrpcServerURI,
		func(conn *grpc.ClientConn) interface{} {
			return protos.NewEdgeAppClient(conn)
		},
		func(ctx context.Context, client interface{}) (interface{}, error) {
			return client.(protos.EdgeAppClient).TrainFinish(ctx, &protos.Empty{})
		},
		func(response interface{}) interface{} {
			if isInterfaceOnly {
				WriteReport("TrainFinish", "implemented", "")
			} else {
				WriteReport("TrainFinish", "TrainFinish sucessfully.", "")
			}
			return nil
		},
	)
}

func sendTrainInteruptMessage(isInterfaceOnly bool, appGrpcServerURI string) bool {
	return EmitEvent(
		isInterfaceOnly,
		"TrainInterupt",
		appGrpcServerURI,
		func(conn *grpc.ClientConn) interface{} {
			return protos.NewEdgeAppClient(conn)
		},
		func(ctx context.Context, client interface{}) (interface{}, error) {
			return client.(protos.EdgeAppClient).TrainInterrupt(ctx, &protos.Empty{})
		},
		func(response interface{}) interface{} {
			if isInterfaceOnly {
				WriteReport("TrainInterrupt", "implemented", "")
			} else {
				WriteReport("TrainInterrupt", "TrainInterrupt sucessfully.", "")
			}
			return nil
		},
	)
}

func EmitEvent(
	isInterfaceOnly bool,
	state string,
	clientURI string,
	newClient func(*grpc.ClientConn) interface{},
	emitEvent func(context.Context, interface{}) (interface{}, error),
	responseHandler func(interface{}) interface{},
) bool {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}

	conn, err := grpc.Dial(clientURI, opts...)
	if err != nil {
		log.Println("fail to dail grpc")
	}
	defer conn.Close()

	client := newClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	response, err := emitEvent(ctx, client)
	stat, ok := status.FromError(err)
	log.Printf("Code: %d, Message: %s\n", stat.Code(), stat.Message())
	log.Println("received response", fmt.Sprintf("%v", response))

	if ok {
		responseHandler(response)
		return true
	} else {
		log.Println("errors happen ... handling")
		errorCode := stat.Code()
		if errorCode == codes.Unimplemented {
			WriteReport(state, "", errorCode.String())
			return true
		} else if errorCode == codes.DeadlineExceeded {
			WriteReport(state, "", errorCode.String())
			return true
		} else {
			if !isInterfaceOnly {
				WriteReport(state, "Unknow Error", errorCode.String())
			}
		}
	}
	return false
}

type ValidatingLogData struct {
	Timestamp string `json:"timestamp"`
	State     string `json:"state"`
	Error     string `json:"error"`
	Message   string `json:"message"`
}

func WriteReport(state string, msg string, er string) {
	mutx.Lock()

	log.Println("writing report message ... ")
	log.Println("path = ", os.Getenv("REPORT_PATH"))
	err := os.MkdirAll(filepath.Dir(os.Getenv("REPORT_PATH")), os.ModePerm)
	if err != nil {
		panic(err)
	}

	fo, err := os.OpenFile(os.Getenv("REPORT_PATH"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	var log ValidatingLogData
	log.Timestamp = time.Now().Format(time.RFC3339)
	log.State = state
	log.Message = msg
	log.Error = er
	jsonStr, _ := json.Marshal(log)
	fo.Write(jsonStr)
	fo.WriteString("\n")
	fo.Close()
	mutx.Unlock()

}
