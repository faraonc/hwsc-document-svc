package conf

import (
	"fmt"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"log"
)

var (
	// GRPCHost address and port of gRPC microservice
	GRPCHost Host

	// DocumentDB represents the Document database
	DocumentDB DocumentDBHost
)

func init() {
	if err := config.Load(file.NewSource(file.WithPath("conf/json/config.dev.json"))); err != nil {
		// TODO - This is a hacky solution for the unit test, because of a weird path issue with GoLang Unit Test
		if err := config.Load(file.NewSource(file.WithPath("../conf/json/config.dev.json"))); err != nil {
			log.Fatalf("[FATAL] Failed to initialize conf file %v\n", err)
		}
	}
	if err := config.Get("hosts", "grpc-server").Scan(&GRPCHost); err != nil {
		log.Fatalf("[FATAL] Failed to scan conf file %v\n", err)
	}
	if err := config.Get("hosts", "mongodb-document").Scan(&DocumentDB); err != nil {
		log.Fatalf("[FATAL] Failed to scan conf file %v\n", err)
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
	Writer string `json:"mongodb-writer"`

	// Reader address for reading to MongoDB server
	Reader string `json:"mongodb-reader"`

	// Name database name
	Name string `json:"db"`

	// Collection database collection name
	Collection string `json:"collection"`
}
