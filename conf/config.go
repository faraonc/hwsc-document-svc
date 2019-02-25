package conf

import (
	"fmt"
	"github.com/hwsc-org/hwsc-document-svc/consts"
	log "github.com/hwsc-org/hwsc-lib/logger"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
)

var (
	// GRPCHost address and port of gRPC microservice
	GRPCHost Host

	// DocumentDB represents the Document database
	DocumentDB DocumentDBHost
)

func init() {
	// Create new config
	conf := config.NewConfig()
	log.Info(consts.DocumentServiceTag, "Reading ENV variables")
	src := env.NewSource(
		env.WithPrefix("hosts"),
	)
	if err := conf.Load(src); err != nil {
		log.Fatal(consts.DocumentServiceTag, "Failed to initialize configuration", err.Error())

	}
	if err := conf.Get("hosts", "document").Scan(&GRPCHost); err != nil {
		log.Fatal(consts.DocumentServiceTag, "Failed to get GRPC configuration", err.Error())
	}
	if err := conf.Get("hosts", "mongodb").Scan(&DocumentDB); err != nil {
		log.Fatal(consts.DocumentServiceTag, "Failed to get MongoDB configuration", err.Error())
	}
}

// Host represents a server.
type Host struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Network string `json:"network"`
}

func (h *Host) String() string {
	return fmt.Sprintf("%s:%s", h.Address, h.Port)
}

// DocumentDBHost represents the Document database
type DocumentDBHost struct {
	// Writer address for writing to MongoDB server
	Writer string `json:"writer"`

	// Reader address for reading to MongoDB server
	Reader string `json:"reader"`

	// Name database name
	Name string `json:"db"`

	// Collection database collection name
	Collection string `json:"collection"`
}
