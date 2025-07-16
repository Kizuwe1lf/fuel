package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	externalsignaling "github.com/kizuwe1lf/fuel/external/service/signaling"
	internalsignaling "github.com/kizuwe1lf/fuel/internal/service/signaling"
	signalingv1 "github.com/kizuwe1lf/fuel/proto/gen/signaling/v1"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	internalService := internalsignaling.NewService()
	externalService := externalsignaling.NewService(internalService)

	grpcServer := grpc.NewServer()
	signalingv1.RegisterSignalingServer(grpcServer, externalService)

	log.Println("signaling server started on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
