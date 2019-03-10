package main

import (
	"github.com/hwsc-org/hwsc-document-svc/conf"
	"github.com/hwsc-org/hwsc-document-svc/consts"
	"google.golang.org/grpc"
	"net"

	pbsvc "github.com/hwsc-org/hwsc-api-blocks/protobuf/hwsc-document-svc/document"
	svc "github.com/hwsc-org/hwsc-document-svc/service"
	log "github.com/hwsc-org/hwsc-lib/logger"
)

func main() {
	log.Info(consts.DocumentServiceTag, "hwsc-document-svc initiating...")

	// Make TCP listener
	lis, err := net.Listen(conf.GRPCHost.Network, conf.GRPCHost.String())
	if err != nil {
		log.Fatal(consts.DocumentServiceTag, "Failed to initialize TCP listener:", err.Error())
	}

	// Make gRPC server
	s := grpc.NewServer()

	// Implement services in /service/service.go
	// Register service with gRPC server
	pbsvc.RegisterDocumentServiceServer(s, &svc.Service{})
	log.Info(consts.DocumentServiceTag, "hwsc-document-svc started at:", conf.GRPCHost.String())

	// Start gRPC server
	if err := s.Serve(lis); err != nil {
		log.Fatal(consts.DocumentServiceTag, "Failed to serve:", err.Error())
	}
}
