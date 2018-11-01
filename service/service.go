package service

import (
	"fmt"
	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-document-svc/proto"
	"github.com/mongodb/mongo-go-driver/bson"
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

// Service implements services for managing document
type Service struct{}

const (
	//TODO New MongoDB Server
	mongoServerDBWriter = "mongodb://hwscmongodb:89PoXCVmIJyg8lSpQ6aF2iaoQk4dDOYav4ZVHkibV6dIZaKF0I2gft8GgKcCOAtXkxIucq9ZBpxYTO9k8QVnTw" +
		"==@hwscmongodb.documents.azure.com:10255/?ssl=true&replicaSet=globaldb"
	mongoServerDBReader = "mongodb://hwscmongodb:mV2GqGnzoOXPF82QZbEzEi0QcFSLK4fyh2EAzU3KrZfw1wSePaQbKINUrWKfblBS3diQfJCd7ugAOYHMZK2eLA" +
		"==@hwscmongodb.documents.azure.com:10255/?ssl=true&replicaSet=globaldb"
	mongoDB           = "hwsc-document"
	mongoDBCollection = "hwsc-document"

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

// CreateDocument creates a document in MongoDB.
// Returns the Document.
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

	doc.Duid = duidGenerator.NewDUID()

	// Validate document
	if err := ValidateDocument(doc); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

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
	log.Printf("[INFO] Connecting to mongodb://hwscmongodb duid: %s uuid: %s\n", doc.GetDuid(), doc.GetUuid())
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
		log.Printf("[ERROR] Inserted documment _id: %v, with disconnection error\n", res.InsertedID)
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Internal)},
			Message: fmt.Sprintf("Inserted document with error: %s\n", codes.Internal.String()),
			Data:    doc,
		}, nil
	}

	// Log ID of the inserted MongoDB document
	log.Printf("[INFO] Success inserting documment _id: %v\n\n", res.InsertedID)
	return &pb.DocumentResponse{
		Status:  &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
		Data:    doc,
	}, nil

}

// ListUserDocumentCollection retrieves all the MongoDB documents for a specific user with the given UUID.
// Returns a collection of Documents.
func (s Service) ListUserDocumentCollection(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting ListUserDocumentCollection service")

	if ok := isStateAvailable(); !ok {
		return nil, status.Error(codes.Unavailable, "Service unavailable")
	}

	doc := req.GetData()
	if doc == nil {
		log.Println("[ERROR] Nil request data")
		return nil, status.Error(codes.InvalidArgument, "Nil request data")
	}

	// Validate UUID field
	if err := ValidateUUID(doc.GetUuid()); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Connect to MongoDB
	log.Printf("[INFO] Connecting to mongodb://hwscmongodb uuid: %s\n", doc.GetUuid())
	client, err := mongo.NewClient(mongoServerDBReader)
	if err != nil {
		log.Printf("[ERROR] mongo.NewClient(mongoServerDBReader): %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := client.Connect(context.TODO()); err != nil {
		log.Printf("[ERROR] client.Connect(context.TODO()): %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	collection := client.Database(mongoDB).Collection(mongoDBCollection)

	// Find all MongoDB documents for the specific uuid
	filter := bson.NewDocument(bson.EC.String("uuid", doc.GetUuid()))
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("[ERROR] collection.Find(context.Background(), mongoQuery{\"uuid\": doc.GetUuid(),}): %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Close the MongoDB cursor before the function exits
	defer cur.Close(context.Background())

	// Extract the documents
	documentCollection := make([]*pb.Document, 0)
	for cur.Next(context.Background()) {
		if err := cur.Err(); err != nil {
			log.Printf("[ERROR] cur.Err(): %s\n", err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Mutate and retrieve Document
		document := &pb.Document{}
		if err := cur.Decode(document); err != nil {
			log.Printf("[ERROR] %s\n", err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Validate the retrieved Document from MongoDB
		if err := ValidateDocument(document); err != nil {
			log.Printf("[ERROR] Failed document validation, duid: %v - uuid: %v - %s\n",
				doc.GetDuid(), doc.GetUuid(), err.Error())

			return nil, status.Errorf(codes.Internal, "Failed document validations, duid: %v - uuid: %v - %s",
				doc.GetDuid(), doc.GetUuid(), err.Error())
		}
		documentCollection = append(documentCollection, document)
		log.Printf("[DEBUG] document: \n%s\n\n", document)

	}

	if err := cur.Err(); err != nil {
		log.Printf("[ERROR] cur.Err(): %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(documentCollection) == 0 {
		log.Printf("[ERROR] No documments for uuid: %v\n\n", doc.GetUuid())
		return nil, status.Errorf(codes.Internal, "No documments for uuid: %v", doc.GetUuid())
	}

	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("[ERROR] Success listing documments, uuid: %v with disconnection error\n", doc.GetUuid())
		return &pb.DocumentResponse{
			Status:             &pb.DocumentResponse_Code{Code: uint32(codes.Internal)},
			Message:            fmt.Sprintf("Listed user documents with disconnection error: %s\n", codes.Internal.String()),
			DocumentCollection: documentCollection,
		}, nil
	}

	// Log ID of the uuid used for query
	log.Printf("[INFO] Success listing documents, uuid: %v\n\n", doc.GetUuid())

	return &pb.DocumentResponse{
		Status:             &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message:            codes.OK.String(),
		DocumentCollection: documentCollection,
	}, nil

}

// UpdateDocument (completely) updates a MongoDB document with a given DUID.
// Returns the updated Document.
func (s Service) UpdateDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting UpdateDocument service")

	if ok := isStateAvailable(); !ok {
		return nil, status.Error(codes.Unavailable, "Service unavailable")
	}

	doc := req.GetData()
	if doc == nil {
		log.Println("[ERROR] Nil request data")
		return nil, status.Error(codes.InvalidArgument, "Nil request data")
	}

	if doc.GetDuid() == "" {
		log.Printf("[ERROR] Missing DUID")
		return nil, status.Error(codes.InvalidArgument, "Missing DUID")
	}

	// Generate UpdateTimestamp
	doc.UpdateTimestamp = time.Now().UTC().Unix()

	// Validate document
	if err := ValidateDocument(doc); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("[INFO] Document contains:\n %s\n\n", doc)

	// Get the specific lock if it already exists, else make the lock
	lock, _ := serviceClientLocker.LoadOrStore(doc.GetDuid(), &sync.RWMutex{})

	// Lock
	lock.(*sync.RWMutex).Lock()
	// Unlock before the function exits
	defer lock.(*sync.RWMutex).Unlock()

	// Connect to MongoDB
	log.Printf("[INFO] Connecting to mongodb://hwscmongodb duid: %s uuid: %s\n",
		doc.GetDuid(), doc.GetUuid())
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

	filter := bson.NewDocument(
		bson.EC.String("duid", doc.GetDuid()),
		bson.EC.String("uuid", doc.GetUuid()),
	)
	result := collection.FindOneAndReplace(context.Background(), filter, doc)

	// Extract the updated MongoDB document
	if result == nil {
		log.Printf("[INFO] Document not found, duid: %v - uuid: %v\n\n",
			doc.GetDuid(), doc.GetUuid())

		return nil, status.Errorf(codes.InvalidArgument,
			"Document not found, duid: %v - uuid: %v",
			doc.GetDuid(), doc.GetUuid())
	}

	if err := result.Decode(&pb.Document{}); err != nil {
		log.Printf("[ERROR] Document not found, duid: %v - uuid: %v\n\n",
			doc.GetDuid(), doc.GetUuid())

		return nil, status.Errorf(codes.InvalidArgument,
			"Document not found, duid: %v - uuid: %v",
			doc.GetDuid(), doc.GetUuid())
	}

	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("[ERROR] Success updating document, duid: %v - uuid: %v with disconnection error\n\n", doc.GetDuid(), doc.GetUuid())
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Internal)},
			Message: fmt.Sprintf("Updated document with disconnection error: %s\n", codes.Internal.String()),
			Data:    doc,
		}, nil
	}

	// Log duid and uuid used for query
	log.Printf("[INFO] Success updating document, duid: %v - uuid: %v\n\n", doc.GetDuid(), doc.GetUuid())

	return &pb.DocumentResponse{
		Status:  &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
		Data:    doc,
	}, nil

}

// DeleteDocument deletes a MongoDB document using UUID and DUID.
// Returns the deleted Document.
func (s Service) DeleteDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting DeleteDocument service")

	if ok := isStateAvailable(); !ok {
		return nil, status.Error(codes.Unavailable, "Service unavailable")
	}

	doc := req.GetData()
	if doc == nil {
		log.Println("[ERROR] Nil request data")
		return nil, status.Error(codes.InvalidArgument, "Nil request data")
	}

	if doc.GetDuid() == "" {
		log.Printf("[ERROR] Missing DUID")
		return nil, status.Error(codes.InvalidArgument, "Missing DUID")
	}

	// Validate document
	if err := ValidateDocument(doc); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("[INFO] Document contains:\n %s\n\n", doc)

	// Get the specific lock if it already exists, else make the lock
	lock, _ := serviceClientLocker.LoadOrStore(doc.GetDuid(), &sync.RWMutex{})

	// Lock
	lock.(*sync.RWMutex).Lock()
	// Unlock before the function exits
	defer lock.(*sync.RWMutex).Unlock()

	// Connect to MongoDB
	log.Printf("[INFO] Connecting to mongodb://hwscmongodb duid: %s uuid: %s\n",
		doc.GetDuid(), doc.GetUuid())
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

	filter := bson.NewDocument(
		bson.EC.String("duid", doc.GetDuid()),
		bson.EC.String("uuid", doc.GetUuid()),
	)
	result := collection.FindOneAndDelete(context.Background(), filter)

	// Extract the updated MongoDB document
	if result == nil {
		log.Printf("[INFO] Document not found, duid: %v - uuid: %v\n\n",
			doc.GetDuid(), doc.GetUuid())

		return nil, status.Errorf(codes.InvalidArgument,
			"Document not found, duid: %v - uuid: %v",
			doc.GetDuid(), doc.GetUuid())
	}

	if err := result.Decode(&pb.Document{}); err != nil {
		log.Printf("[ERROR] Document not found, duid: %v - uuid: %v\n\n",
			doc.GetDuid(), doc.GetUuid())

		return nil, status.Errorf(codes.InvalidArgument,
			"Document not found, duid: %v - uuid: %v",
			doc.GetDuid(), doc.GetUuid())
	}


	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("[ERROR] Success deleting document, duid: %v - uuid: %v with disconnection error\n\n", doc.GetDuid(), doc.GetUuid())
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Internal)},
			Message: fmt.Sprintf("Deleted document with disconnection error: %s\n", codes.Internal.String()),
			Data:    doc,
		}, nil
	}

	// Log duid and uuid used for query
	log.Printf("[INFO] Success deleting document, duid: %v - uuid: %v\n\n", doc.GetDuid(), doc.GetUuid())

	return &pb.DocumentResponse{
		Status:  &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
		Data:    doc,
	}, nil

}

// AddFileMetadata adds a new FileMetadata in a MongoDB document using a given url, UUID and DUID.
// Returns the updated Document.
func (s Service) AddFileMetadata(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil

}

// DeleteFileMetadata deletes a FileMetadata in a MongoDB document using a given FUID, UUID and DUID.
// Returns the updated Document.
func (s Service) DeleteFileMetadata(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil

}
