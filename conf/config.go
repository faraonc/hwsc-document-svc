package conf

import (
	"github.com/hwsc-org/hwsc-document-svc/consts"
	"github.com/hwsc-org/hwsc-lib/hosts"
	"github.com/hwsc-org/hwsc-lib/logger"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
)

var (
	// GRPCHost address and port of gRPC microservice
	GRPCHost hosts.Host

	// DocumentDB represents the Document database
	DocumentDB hosts.DocumentDBHost
)

func init() {
	// Create new config
	conf := config.NewConfig()
	logger.Info(consts.DocumentServiceTag, "Reading ENV variables")
	src := env.NewSource(
		env.WithPrefix("hosts"),
	)
	if err := conf.Load(src); err != nil {
		logger.Fatal(consts.DocumentServiceTag, "Failed to initialize configuration", err.Error())

	}
	if err := conf.Get("hosts", "document").Scan(&GRPCHost); err != nil {
		logger.Fatal(consts.DocumentServiceTag, "Failed to get GRPC configuration", err.Error())
	}
	if err := conf.Get("hosts", "mongodb").Scan(&DocumentDB); err != nil {
		logger.Fatal(consts.DocumentServiceTag, "Failed to get MongoDB configuration", err.Error())
	}
}
