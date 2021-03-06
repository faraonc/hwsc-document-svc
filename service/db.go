package service

import (
	"github.com/hwsc-org/hwsc-document-svc/conf"
	"github.com/hwsc-org/hwsc-document-svc/consts"
	log "github.com/hwsc-org/hwsc-lib/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	mongoDBReader     *mongo.Client
	mongoDBWriter     *mongo.Client
	ctxWithTimeout, _ = context.WithTimeout(context.Background(), 5*time.Second)
)

func init() {
	// Handle Terminate Signal(Ctrl + C)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = disconnectMongoDBClient(mongoDBReader)
		_ = disconnectMongoDBClient(mongoDBWriter)
		log.Fatal(consts.MongoDBTag, "hwsc-document-svc terminated")
	}()
}

// dialMongoDB connects a client to MongoDB server.
// Returns a MongoDB Client or any dialing error.
func dialMongoDB(uri *string) (*mongo.Client, error) {
	if strings.TrimSpace(*uri) == "" {
		return nil, consts.ErrEmptyMongoDBURI
	}
	client, err := mongo.Connect(ctxWithTimeout, options.Client().ApplyURI(*uri))
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
		return consts.ErrNilMongoDBClient
	}
	return client.Disconnect(ctxWithTimeout)
}

// refreshMongoDBConnection refreshes a client's connection with MongoDB server.
// Returns if there is any connection error.
func refreshMongoDBConnection(client *mongo.Client, uri *string) error {
	if client == nil {
		newClient, err := dialMongoDB(uri)
		if err != nil {
			return consts.ErrMongoDBUnavailable
		}
		assignMongoDBClient(newClient, uri)
		return nil
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		newClient, err := dialMongoDB(uri)
		if err != nil {
			assignMongoDBClient(nil, uri)
			return consts.ErrMongoDBUnavailable
		}
		assignMongoDBClient(newClient, uri)
		return nil
	}

	return nil
}

// assignMongoDBClient assigns a new MongoDB client.
func assignMongoDBClient(newClient *mongo.Client, uri *string) {
	if strings.EqualFold(*uri, conf.DocumentDB.Reader) {
		mongoDBReader = newClient
	} else {
		mongoDBWriter = newClient
	}
}
