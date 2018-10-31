package service

import (
	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-document-svc/proto"
	"golang.org/x/net/context"
	"sync"
)

// state of the service
type state uint32

// stateLocker synchronizes the state of the service
type stateLocker struct {
	lock                sync.RWMutex
	currentServiceState state
}

// mongoQuery is the bson map for querying MongoDB
type mongoQuery map[string]interface{}

// Service implements services for managing file metadata
type Service struct{}

const (
	//TODO MongoDB Server
	mongoServerDBWriter = "mongodb://hwscmongodb:89PoXCVmIJyg8lSpQ6aF2iaoQk4dDOYav4ZVHkibV6dIZaKF0I2gft8GgKcCOAtXkxIucq9ZBpxYTO9k8QVnTw" +
		"==@hwscmongodb.documents.azure.com:10255/?ssl=true&replicaSet=globaldb"
	mongoServerDBReader = "mongodb://hwscmongodb:mV2GqGnzoOXPF82QZbEzEi0QcFSLK4fyh2EAzU3KrZfw1wSePaQbKINUrWKfblBS3diQfJCd7ugAOYHMZK2eLA" +
		"==@hwscmongodb.documents.azure.com:10255/?ssl=true&replicaSet=globaldb"
	mongoDB           = "METADATA-FILE"
	mongoDBCollection = "metadata-file"

	// available - Service is ready and available
	available state = 0

	// unavailable - Service is unavailable. Example: Provisioning something
	unavailable state = 1
)

var (
	serviceStateLocker stateLocker

	// Converts State of the service to a string
	serviceStateMap map[state]string

	// Stores the lock for each uuid
	serviceClientLocker sync.Map
)

func init() {
	serviceStateLocker = stateLocker{currentServiceState: available}
	serviceStateMap = map[state]string{
		available:   "Available",
		unavailable: "Unavailable",
	}
}

func (s state) String() string {
	return serviceStateMap[s]
}

// GetStatus gets the current status of the service.
func (s Service) GetStatus(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	//TODO check if mongodb server is running
	return nil, nil

}

// CreateDocument creates a Document in MongoDB.
func (s Service) CreateDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil
}

// ListUserDocumentCollection gets all the MongoDB documents for a specific user with the given UUID.
// Returns a collection of MongoDB documents.
func (s Service) ListUserDocumentCollection(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil

}

// UpdateDocument completely updates a MongoDB document with a given DUID.
//TODO implementation
//TODO unit test
//TODO readme
func (s Service) UpdateDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil

}

// DeleteDocument deletes a MongoDB document using UUID and DUID.
// Returns the deleted MongoDB document.
func (s Service) DeleteDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil

}
