package service

import (
	"fmt"
	"github.com/hwsc-org/hwsc-document-svc/conf"
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	mongoDBReader *mongo.Client
	mongoDBWriter *mongo.Client
)

func init() {
	var err error
	mongoDBReader, err = dialMongoDB(&conf.DocumentDB.Reader)
	if err != nil {
		log.Fatalf("[FATAL] %s\n", err.Error())
	}
	mongoDBWriter, err = dialMongoDB(&conf.DocumentDB.Writer)
	if err != nil {
		log.Fatalf("[FATAL] %s\n", err.Error())
	}
	// Handle Terminate Signal(Ctrl + C)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = disconnectMongoDBClient(mongoDBReader)
		_ = disconnectMongoDBClient(mongoDBWriter)
		fmt.Println()
		log.Fatalln("[FATAL] hwsc-document-svc terminated")
	}()
}

// dialMongoDB connects a client to MongoDB server.
// Returns a MongoDB Client or any dialing error.
func dialMongoDB(uri *string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), *uri)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return client, nil
}

// disconnectMongoDBClient disconnects a client from MongoDB server.
// Returns if there is any disconnection error.
func disconnectMongoDBClient(client *mongo.Client) error {
	if client == nil {
		return errNilMongoDBClient
	}
	return client.Disconnect(context.TODO())
}

// refreshMongoDBConnection refreshes a client's connection with MongoDB server.
// Returns if there is any connection error.
func refreshMongoDBConnection(client *mongo.Client, uri *string) error {
	if client == nil {
		client, err := dialMongoDB(uri)
		if err != nil {
			return errMongoDBUnavailable
		}
		if strings.EqualFold(*uri, conf.DocumentDB.Reader) {
			mongoDBReader = client
		} else {
			mongoDBWriter = client
		}
		return nil
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		client, err := dialMongoDB(uri)
		if err != nil {
			if strings.EqualFold(*uri, conf.DocumentDB.Reader) {
				mongoDBReader = nil
			} else {
				mongoDBWriter = nil
			}
			return errMongoDBUnavailable
		}
		if strings.EqualFold(*uri, conf.DocumentDB.Reader) {
			mongoDBReader = client
		} else {
			mongoDBWriter = client
		}
		return nil
	}

	return nil
}
