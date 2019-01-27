package main

import (
	"github.com/hwsc-org/hwsc-document-svc/conf"
	"google.golang.org/grpc"
	"net"

	pb "github.com/hwsc-org/hwsc-api-blocks/int/hwsc-document-svc/proto"
	svc "github.com/hwsc-org/hwsc-document-svc/service"
	log "github.com/hwsc-org/hwsc-logger/logger"
)

func main() {
	log.Info("hwsc-document-svc initiating...")

	// Make TCP listener
	lis, err := net.Listen(conf.GRPCHost.Network, conf.GRPCHost.String())
	if err != nil {
		log.Fatal("Failed to initialize TCP listener:", err.Error())
	}

	// Make gRPC server
	s := grpc.NewServer()

	// Implement services in /service/service.go
	// Register service with gRPC server
	pb.RegisterDocumentServiceServer(s, &svc.Service{})
	log.Info("hwsc-document-svc started at:", conf.GRPCHost.String())

	// Start gRPC server
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err.Error())
	}
}
