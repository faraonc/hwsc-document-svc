package main

import (
	"google.golang.org/grpc"
	"log"
	"net"

	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-document-svc/proto"
	svc "github.com/faraonc/hwsc-document-svc/service"
)

func main() {
	log.Println("[INFO] hwsc-document-svc initiating...")

	// Make TCP listener
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize TCP listener %v\n", err)
	}

	// Make gRPC server
	s := grpc.NewServer()

	// Implement services in /service/service.go
	// Register service with gRPC server
	pb.RegisterDocumentServiceServer(s, &svc.Service{})
	log.Println("[INFO] hwsc-document-svc at localhost:50051...")

	// Start gRPC server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("[FATAL] Failed to serve %v\n", err)
	}
}
