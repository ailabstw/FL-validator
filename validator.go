package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"gitlab.com/fl_validator/edge"
	protos "gitlab.com/fl_validator/go_protos"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type baseModel struct {
	repoName string
	metadata map[string]string
	metrics  map[string]float64
}

func main() {
	clientURI := "0.0.0.0:7878"

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}

	zap.L().Info("starting grpc server.... ")
	startGrpcServer(clientURI)

	conn, err := grpc.Dial(clientURI, opts...)
	if err != nil {
		zap.L().Fatal("fail to dail grpc", zap.Error(err))
	}
	defer conn.Close()

	// test init
	zap.L().Info("starting grpc server.... ")
	sendInitMessage(clientURI)

	time.Sleep(20 * time.Second)

	zap.L().Info("sending local train message .... ")

	// test localtrain
	sendLocalTrainMessage(clientURI, 1, baseModel{}, "")

	time.Sleep(20 * time.Second)
	// test train finish
	sendTrainFinishMessage(clientURI)

}

func startGrpcServer(address string) *grpc.Server {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		zap.L().Fatal("Cannot listen on the address",
			zap.String("service", "grpc"),
			zap.String("address", address),
			zap.Error(err))
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	go func() {
		zap.L().Debug("Grpc server listen the address",
			zap.String("service", "grpc"),
			zap.String("address", address))

		protos.RegisterEdgeOperatorServer(grpcServer, &edge.EdgeOperatorServer{})

		if err := grpcServer.Serve(lis); err != nil {
			zap.L().Fatal("Cannot start server on the address",
				zap.String("service", "grpc"),
				zap.String("address", address),
				zap.Error(err))
		}
	}()

	zap.L().Info(fmt.Sprintf("Grpc Listen [%v]", "0.0.0.0:8787"))
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
		zap.L().Fatal("fail to dail grpc", zap.Error(err))
	}
	defer conn.Close()

	client := newClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	response, err := emitEvent(ctx, client)
	if err != nil {
		if err == context.DeadlineExceeded {
			zap.L().Fatal("Deadline exceeded")
		} else {
			zap.L().Fatal("emitEvent get error", zap.Error(err))
		}
	}

	zap.L().Debug("received response", zap.String("response", fmt.Sprintf("%v", response)))
}
