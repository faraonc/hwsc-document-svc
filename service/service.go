package service

import (
	"fmt"
	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-document-svc/proto"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/segmentio/ksuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sync"
	"time"
)

// state of the service
type state uint32

// stateLocker synchronizes the state of the service
type stateLocker struct {
	lock                sync.RWMutex
	currentServiceState state
}

// duidLocker synchronizes the generating of duid
type duidLocker struct {
	lock sync.Mutex
}

// mongoQuery is the bson map for querying MongoDB
type mongoQuery map[string]interface{}

// Service implements services for managing document
type Service struct{}

const (
	//TODO New MongoDB Server
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
	duidGenerator      duidLocker

	// Converts State of the service to a string
	serviceStateMap map[state]string

	// Stores the lock for each uuid
	serviceClientLocker sync.Map
)

func init() {
	serviceStateLocker = stateLocker{currentServiceState: available}
	duidGenerator = duidLocker{}

	serviceStateMap = map[state]string{
		available:   "Available",
		unavailable: "Unavailable",
	}
}

func (s state) String() string {
	return serviceStateMap[s]
}

func (d duidLocker) NewDUID() string {
	d.lock.Lock()
	defer d.lock.Unlock()
	duid := ksuid.New().String()
	return duid
}

// GetStatus gets the current status of the service.
func (s Service) GetStatus(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting GetStatus service")

	// Lock the state for reading
	serviceStateLocker.lock.RLock()
	// Unlock the state before function exits
	defer serviceStateLocker.lock.RUnlock()

	log.Printf("[INFO] Service State: %s\n", serviceStateLocker.currentServiceState)
	if serviceStateLocker.currentServiceState == unavailable {
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Unavailable)},
			Message: codes.Unavailable.String(),
		}, nil
	}

	// Check MongoDB Server
	client, err := mongo.NewClient(mongoServerDBReader)
	if err != nil {
		log.Printf("[ERROR] mongo.NewClient(mongoServerDBReader): %s\n", err.Error())
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Unavailable)},
			Message: codes.Unavailable.String(),
		}, nil
	}

	if err := client.Connect(context.TODO()); err != nil {
		log.Printf("[ERROR] client.Connect(context.TODO()): %s\n", err.Error())
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Unavailable)},
			Message: codes.Unavailable.String(),
		}, nil
	}

	if err := client.Disconnect(context.TODO()); err != nil {
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Internal)},
			Message: codes.Internal.String(),
		}, nil
	}

	return &pb.DocumentResponse{
		Status:  &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil

}

// CreateDocument creates a Document in MongoDB.
func (s Service) CreateDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting CreateDocument service")

	if ok := isStateAvailable(); !ok {
		return nil, status.Error(codes.Unavailable, "Service unavailable")
	}

	doc := req.GetData()
	if doc == nil {
		log.Println("[ERROR] Nil request data")
		return nil, status.Error(codes.InvalidArgument, "Nil request data")
	}

	// Validate document
	if err := ValidateDocument(doc); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	doc.Duid = duidGenerator.NewDUID()

	// Get the specific lock if it already exists, else make the lock
	lock, _ := serviceClientLocker.LoadOrStore(doc.GetDuid(), &sync.RWMutex{})
	// Lock
	lock.(*sync.RWMutex).Lock()
	// Unlock before the function exits
	defer lock.(*sync.RWMutex).Unlock()

	// Generate CreateTimestamp in UTC
	doc.CreateTimestamp = time.Now().UTC().Unix()
	log.Printf("[INFO] Document contains:\n %s\n\n", doc)

	// Connect to MongoDB
	log.Printf("[INFO] Connecting to mongodb://hwscmongodb as %s with %s\n", doc.GetUuid(), doc.GetDuid())
	client, err := mongo.NewClient(mongoServerDBWriter)
	if err != nil {
		log.Printf("[ERROR] mongo.NewClient(mongoServerDBWriter): %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := client.Connect(context.TODO()); err != nil {
		log.Printf("[ERROR] client.Connect(context.TODO()): %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	collection := client.Database(mongoDB).Collection(mongoDBCollection)

	// Insert MongoDB document
	res, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Printf("[ERROR] collection.InsertOne(context.Background(), doc): %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("[ERROR] Inserted documment _id: %v, with disconnect error\n", res.InsertedID)
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Internal)},
			Message: fmt.Sprintf("Inserted with error: %s\n", codes.Internal.String()),
		}, nil
	}

	// Log ID of the inserted MongoDB document
	log.Printf("[INFO] Success inserting documment _id: %v\n\n", res.InsertedID)
	return &pb.DocumentResponse{
		Status:  &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
	}, nil

}

// ListUserDocumentCollection gets all the MongoDB documents for a specific user with the given UUID.
// Returns a collection of MongoDB documents.
func (s Service) ListUserDocumentCollection(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil

}

// UpdateDocument completely updates a MongoDB document with a given DUID.
// Returns the updated MongoDB document.
//TODO implementation
//TODO unit test
func (s Service) UpdateDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil

}

// DeleteDocument deletes a MongoDB document using UUID and DUID.
// Returns the deleted MongoDB document.
func (s Service) DeleteDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil

}
