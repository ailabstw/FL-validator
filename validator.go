package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"gitlab.com/fl_validator/edge"
	protos "gitlab.com/fl_validator/go_protos"
	"google.golang.org/grpc"
)

type baseModel struct {
	repoName string
	metadata map[string]string
	metrics  map[string]float64
}

func main() {
	clientURI := "0.0.0.0:7878"
	serverURI := "0.0.0.0:8787"

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}

	log.Println("starting grpc server.... ")

	startGrpcServer(serverURI)

	log.Println("starting grpc client.... ")

	conn, err := grpc.Dial(clientURI, opts...)
	if err != nil {
		log.Fatal("starting grpc server.... ")
	}
	defer conn.Close()

	log.Println("sending  init msg.... ")

	sendInitMessage(clientURI)

	time.Sleep(20 * time.Second)

	log.Println("sending local train message .... ")
	// test localtrain
	sendLocalTrainMessage(clientURI, 1, baseModel{}, "")

	time.Sleep(20 * time.Second)

	log.Println("sending training finished message .... ")
	// test train finish
	sendTrainFinishMessage(clientURI)

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

func sendInitMessage(appGrpcServerURI string) {
	EmitEvent(
		appGrpcServerURI,
		func(conn *grpc.ClientConn) interface{} {
			return protos.NewEdgeAppClient(conn)
		},
		func(ctx context.Context, client interface{}) (interface{}, error) {
			return client.(protos.EdgeAppClient).TrainInit(ctx, &protos.Empty{})
		},
	)
}

func sendLocalTrainMessage(appGrpcServerURI string, epochPerRound int, baseModel baseModel, edgeRepoName string) {
	EmitEvent(
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
	)
}

func sendTrainFinishMessage(appGrpcServerURI string) {
	EmitEvent(
		appGrpcServerURI,
		func(conn *grpc.ClientConn) interface{} {
			return protos.NewEdgeAppClient(conn)
		},
		func(ctx context.Context, client interface{}) (interface{}, error) {
			return client.(protos.EdgeAppClient).TrainFinish(ctx, &protos.Empty{})
		},
	)
}

func EmitEvent(
	clientURI string,
	newClient func(*grpc.ClientConn) interface{},
	emitEvent func(context.Context, interface{}) (interface{}, error),
) {
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
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Fatal("Deadline exceeded")
		} else {
			log.Fatal("emitEvent get error", err)
		}
	}

	log.Fatal("received response", fmt.Sprintf("%v", response))
}
