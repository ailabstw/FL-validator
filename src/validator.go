package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"gitlab.com/fl_validator/src/edge"
	protos "gitlab.com/fl_validator/src/go_protos"
	"gitlab.com/fl_validator/src/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type baseModel struct {
	repoName string
	metadata map[string]string
	metrics  map[string]float64
}

func checkOnlyInterface(clientURI string) {

	allImplement := true
	allImplement = allImplement && sendDataValidate(true, clientURI)
	allImplement = allImplement && sendInitMessage(true, clientURI)
	allImplement = allImplement && sendLocalTrainMessage(true, clientURI, 1)
	allImplement = allImplement && sendTrainFinishMessage(true, clientURI)
	allImplement = allImplement && sendTrainInteruptMessage(true, clientURI)
	if allImplement {
		util.WriteReport("check", "all interface implemented", "")
	} else {
		util.WriteReport("check", "not all interface implemented", "not all interface implemented")
	}
}

func main() {
	clientURI := os.Getenv("APP_URI")
	serverURI := "0.0.0.0:8787"

	log.Println("clientURI: " + clientURI)

	isInterfaceOnly, _ := strconv.ParseBool(os.Getenv("DRY_RUN"))

	log.Println("Starting grpc server.... ")

	startGrpcServer(serverURI)

	time.Sleep(30 * time.Second)

	if isInterfaceOnly {
		log.Println("Check only interface .... ")
		checkOnlyInterface(clientURI)
		time.Sleep(20 * time.Second)
		return
	}

	log.Println("Sending validating msg.... ")

	sendDataValidate(isInterfaceOnly, clientURI)

	time.Sleep(10 * time.Second)

	log.Println("Sending  initialization msg.... ")

	sendInitMessage(isInterfaceOnly, clientURI)

	time.Sleep(30 * time.Second)

	log.Println("Sending local train message .... ")
	// test localtrain
	sendLocalTrainMessage(isInterfaceOnly, clientURI, 1)

	time.Sleep(10 * time.Second)

	// test TrainInterupt
	sendTrainInteruptMessage(isInterfaceOnly, clientURI)

	time.Sleep(10 * time.Second)

	log.Println("Sending training finished message .... ")
	// test train finish
	sendTrainFinishMessage(isInterfaceOnly, clientURI)

	time.Sleep(5 * time.Second)

	// all test successfully
	if util.GetResult() {
		log.Println("All FL validation completed . Congrats. ")
	} else {
		log.Println("FL validation failed. Please checkout error logs for detals. ")
	}

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

func sendDataValidate(isInterfaceOnly bool, appGrpcServerURI string) bool {
	return EmitEvent(
		isInterfaceOnly,
		"DataValidate",
		appGrpcServerURI,
		func(conn *grpc.ClientConn) interface{} {
			return protos.NewEdgeAppClient(conn)
		},
		func(ctx context.Context, client interface{}) (interface{}, error) {
			return client.(protos.EdgeAppClient).DataValidate(ctx, &protos.Empty{})
		},
		func(response interface{}) interface{} {
			if isInterfaceOnly {
				util.WriteReport("DataValidate", "implemented", "")
			} else {
				util.WriteReport("DataValidate", "DataValidate successfully.", "")
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
				util.WriteReport("TrainInit", "implemented", "")
			} else {
				util.WriteReport("TrainInit", "TrainInit successfully.", "")
			}
			return nil
		},
	)
}

func sendLocalTrainMessage(isInterfaceOnly bool, appGrpcServerURI string, epochPerRound int) bool {
	return EmitEvent(
		isInterfaceOnly,
		"LocalTrain",
		appGrpcServerURI,
		func(conn *grpc.ClientConn) interface{} {
			return protos.NewEdgeAppClient(conn)
		},
		func(ctx context.Context, client interface{}) (interface{}, error) {
			return client.(protos.EdgeAppClient).LocalTrain(ctx, &protos.LocalTrainParams{
				EpR: int32(epochPerRound),
			})
		},
		func(response interface{}) interface{} {
			if isInterfaceOnly {
				util.WriteReport("LocalTrain", "implemented", "")
			} else {
				util.WriteReport("LocalTrain", "LocalTrain successfully.", "")
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
				util.WriteReport("TrainFinish", "implemented", "")
			} else {
				util.WriteReport("TrainFinish", "TrainFinish successfully.", "")
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
				util.WriteReport("TrainInterrupt", "implemented", "")
			} else {
				util.WriteReport("TrainInterrupt", "TrainInterrupt successfully.", "")
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

	ctx := metadata.AppendToOutgoingContext(context.Background(), "draftRun", strconv.FormatBool(isInterfaceOnly))
	ctx, cancel := context.WithTimeout(ctx, 60*time.Minute)
	defer cancel()

	response, err := emitEvent(ctx, client)
	stat, ok := status.FromError(err)
	//log.Printf("Code: %d, Message: %s\n", stat.Code(), stat.Message())
	log.Println("received response", fmt.Sprintf("%v", response))

	if ok {
		responseHandler(response)
		return true
	} else {
		log.Println("errors happen ... handling")
		errorCode := stat.Code()
		if errorCode == codes.Unimplemented {
			util.MakeResultFalse()
			util.WriteReport(state, errorCode.String(), errorCode.String())
			return true
		} else if errorCode == codes.DeadlineExceeded {
			util.MakeResultFalse()
			util.WriteReport(state, errorCode.String(), errorCode.String())
			return true
		} else {
			if !isInterfaceOnly {
				util.MakeResultFalse()
				util.WriteReport(state, "Unknow Error", errorCode.String())
			}
		}
	}
	return false
}
