package service

import (
	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-document-svc/proto"
	"github.com/google/uuid"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"github.com/mongodb/mongo-go-driver/mongo/mongoopt"
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
	//TODO New MongoDB Server, currently using DEV SERVER
	mongoServerDBWriter = "mongodb://hwsc-dev-mongodb:Xt1i0AF9xv5TGlkZFkO3fFQcmb2EFFczbaWJVdVoSp2ZUbRY6ttuyjkgkg6H3UELojGqHEEfIuMebJwrVcFoIA" +
		"==@hwsc-dev-mongodb.documents.azure.com:10255/?ssl=true&replicaSet=globaldb"
	mongoServerDBReader = "mongodb://hwsc-dev-mongodb:QJmIqTki1VGEPoGI4Tfn32g4rTnZrWUlv9jkWKbHSH4A8E0KGMRNUxiCHJfP4ecwEpUYm7yrpER4CLMfCcRMnQ" +
		"==@hwsc-dev-mongodb.documents.azure.com:10255/?ssl=true&replicaSet=globaldb"
	mongoDB           = "dev-document"
	mongoDBCollection = "dev-document"

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

// NewDUID generates a new document unique ID.
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
		log.Printf("[ERROR] Creating MongoDB client: %s\n", err.Error())
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Unavailable)},
			Message: codes.Unavailable.String(),
		}, nil
	}

	if err := client.Connect(context.TODO()); err != nil {
		log.Printf("[ERROR] Connecting MongoDB client: %s\n", err.Error())
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Unavailable)},
			Message: codes.Unavailable.String(),
		}, nil
	}

	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("[ERROR] Disconnecting MongoDB client: %s\n", err.Error())
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Unavailable)},
			Message: codes.Unavailable.String(),
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

	if req == nil {
		log.Println("[ERROR] Nil request")
		return nil, status.Error(codes.InvalidArgument, "Nil request")
	}

	doc := req.GetData()
	if doc == nil {
		log.Println("[ERROR] Nil request data")
		return nil, status.Error(codes.InvalidArgument, "Nil request data")
	}

	doc.Duid = duidGenerator.NewDUID()

	// Get the specific lock if it already exists, else make the lock
	lock, _ := serviceClientLocker.LoadOrStore(doc.GetDuid(), &sync.RWMutex{})
	// Lock
	lock.(*sync.RWMutex).Lock()
	// Unlock before the function exits
	defer lock.(*sync.RWMutex).Unlock()
	// TODO unit test nil lists
	// Extract image URLS
	if doc.GetImageUrlsMap() == nil {
		doc.ImageUrlsMap = make(map[string]string)
	}
	if req.GetImageUrls() != nil {
		for _, url := range req.GetImageUrls() {
			doc.ImageUrlsMap[uuid.New().String()] = url
		}
	}

	// Extract audio URLS
	if doc.GetAudioUrlsMap() == nil {
		doc.AudioUrlsMap = make(map[string]string)
	}
	if req.GetAudioUrls() != nil {
		for _, url := range req.GetAudioUrls() {
			doc.AudioUrlsMap[uuid.New().String()] = url
		}
	}

	// Extract video URLS
	if doc.GetVideoUrlsMap() == nil {
		doc.VideoUrlsMap = make(map[string]string)
	}
	if req.GetVideoUrls() != nil {
		for _, url := range req.GetVideoUrls() {
			doc.VideoUrlsMap[uuid.New().String()] = url
		}
	}
	// Extract file URLS
	if doc.GetFileUrlsMap() == nil {
		doc.FileUrlsMap = make(map[string]string)
	}
	if req.GetFileUrls() != nil {
		for _, url := range req.GetFileUrls() {
			doc.FileUrlsMap[uuid.New().String()] = url
		}
	}

	doc.CreateTimestamp = time.Now().UTC().Unix()

	if err := ValidateDocument(doc); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("[INFO] Document contains:\n %s\n\n", doc)

	log.Printf("[INFO] Connecting to mongodb://hwscmongodb duid: %s - uuid: %s\n", doc.GetDuid(), doc.GetUuid())
	client, err := mongo.NewClient(mongoServerDBWriter)
	if err != nil {
		log.Printf("[ERROR] Creating MongoDB client: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := client.Connect(context.TODO()); err != nil {
		log.Printf("[ERROR] Connecting MongoDB client: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	collection := client.Database(mongoDB).Collection(mongoDBCollection)

	res, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Printf("[ERROR] InsertOne: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("[ERROR] Inserted document _id: %v, with disconnection error: %s\n",
			res.InsertedID, err.Error())
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Internal)},
			Message: "Inserted document with MongoDB disconnection error",
			Data:    doc,
		}, nil
	}

	log.Printf("[INFO] Success inserting document _id: %v\n", res.InsertedID)

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

	if req == nil {
		log.Println("[ERROR] Nil request")
		return nil, status.Error(codes.InvalidArgument, "Nil request")
	}

	doc := req.GetData()
	if doc == nil {
		log.Println("[ERROR] Nil request data")
		return nil, status.Error(codes.InvalidArgument, "Nil request data")
	}

	if err := ValidateUUID(doc.GetUuid()); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("[INFO] Connecting to mongodb://hwscmongodb uuid: %s\n", doc.GetUuid())
	client, err := mongo.NewClient(mongoServerDBReader)
	if err != nil {
		log.Printf("[ERROR] Creating MongoDB client: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := client.Connect(context.TODO()); err != nil {
		log.Printf("[ERROR] Connecting MongoDB client: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	collection := client.Database(mongoDB).Collection(mongoDBCollection)

	// Find all MongoDB documents for the specific uuid
	filter := bson.NewDocument(bson.EC.String("uuid", doc.GetUuid()))
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("[ERROR] Find: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Close the MongoDB cursor before the function exits
	defer cur.Close(context.Background())

	// Extract the documents
	documentCollection := make([]*pb.Document, 0)
	for cur.Next(context.Background()) {
		if err := cur.Err(); err != nil {
			log.Printf("[ERROR] Cursor Err: %s\n", err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Mutate and retrieve Document
		document := &pb.Document{}
		if err := cur.Decode(document); err != nil {
			log.Printf("[ERROR] Cursor Decode: %s\n", err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Validate the retrieved Document from MongoDB
		if err := ValidateDocument(document); err != nil {
			log.Printf("[ERROR] Failed document validation, duid: %s - uuid: %s - %s\n",
				doc.GetDuid(), doc.GetUuid(), err.Error())

			return nil, status.Errorf(codes.Internal, "Failed document validation, duid: %s - uuid: %s - %s",
				doc.GetDuid(), doc.GetUuid(), err.Error())
		}
		documentCollection = append(documentCollection, document)
		log.Printf("[DEBUG] document: \n%s\n\n", document)

	}

	if err := cur.Err(); err != nil {
		log.Printf("[ERROR] Cursor Err: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(documentCollection) == 0 {
		log.Printf("[ERROR] No document for uuid: %s\n", doc.GetUuid())
		return nil, status.Errorf(codes.InvalidArgument, "No document for uuid: %s", doc.GetUuid())
	}

	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("[ERROR] Success listing documents, uuid: %s with disconnection error: %s\n",
			doc.GetUuid(), err.Error())
		return &pb.DocumentResponse{
			Status:             &pb.DocumentResponse_Code{Code: uint32(codes.Internal)},
			Message:            "Listed user documents with MongoDB disconnection error",
			DocumentCollection: documentCollection,
		}, nil
	}

	log.Printf("[INFO] Success listing documents, uuid: %s\n", doc.GetUuid())

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

	if req == nil {
		log.Println("[ERROR] Nil request")
		return nil, status.Error(codes.InvalidArgument, "Nil request")
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

	doc.UpdateTimestamp = time.Now().UTC().Unix()

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

	log.Printf("[INFO] Connecting to mongodb://hwscmongodb duid: %s - uuid: %s\n",
		doc.GetDuid(), doc.GetUuid())
	client, err := mongo.NewClient(mongoServerDBWriter)
	if err != nil {
		log.Printf("[ERROR] Creating MongoDB client: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := client.Connect(context.TODO()); err != nil {
		log.Printf("[ERROR] Connecting MongoDB client: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	collection := client.Database(mongoDB).Collection(mongoDBCollection)

	filter := bson.NewDocument(
		bson.EC.String("duid", doc.GetDuid()),
		bson.EC.String("uuid", doc.GetUuid()),
	)

	// option to return the the document after update
	option := findopt.ReplaceOneBundle{}
	result := collection.FindOneAndReplace(context.Background(), filter, doc,
		option.ReturnDocument(mongoopt.After))

	// Extract the updated MongoDB document
	if result == nil {
		log.Printf("[INFO] Document not found, duid: %s - uuid: %s\n",
			doc.GetDuid(), doc.GetUuid())

		return nil, status.Errorf(codes.InvalidArgument,
			"Document not found, duid: %s - uuid: %s",
			doc.GetDuid(), doc.GetUuid())
	}

	document := &pb.Document{}
	if err := result.Decode(document); err != nil {
		log.Printf("[ERROR] Document not found, duid: %s - uuid: %s - err: %s\n",
			doc.GetDuid(), doc.GetUuid(), err.Error())

		return nil, status.Errorf(codes.InvalidArgument,
			"Document not found, duid: %s - uuid: %s",
			doc.GetDuid(), doc.GetUuid())
	}

	if err := ValidateDocument(document); err != nil {
		log.Printf("[ERROR] Success updating document, duid: %s - uuid: %s with validation error: %s\n",
			doc.GetDuid(), doc.GetUuid(), err.Error())
		log.Printf("[ERROR] Suspected document: \n%s\n\n", doc)
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Internal)},
			Message: "Updated document with validation error",
			Data:    document,
		}, nil
	}

	log.Printf("[DEBUG] Updated document: \n%s\n\n", document)

	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("[ERROR] Success updating document, duid: %s - uuid: %s with disconnection error: %s\n",
			doc.GetDuid(), doc.GetUuid(), err.Error())
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Internal)},
			Message: "Updated document with MongoDB disconnection error",
			Data:    document,
		}, nil
	}

	log.Printf("[INFO] Success updating document, duid: %s - uuid: %s\n",
		doc.GetDuid(), doc.GetUuid())

	return &pb.DocumentResponse{
		Status:  &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
		Data:    document,
	}, nil

}

// DeleteDocument deletes a MongoDB document using UUID and DUID.
// Returns the deleted Document.
func (s Service) DeleteDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting DeleteDocument service")

	if ok := isStateAvailable(); !ok {
		return nil, status.Error(codes.Unavailable, "Service unavailable")
	}

	if req == nil {
		log.Println("[ERROR] Nil request")
		return nil, status.Error(codes.InvalidArgument, "Nil request")
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

	if err := ValidateDUID(doc.GetDuid()); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateUUID(doc.GetUuid()); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Get the specific lock if it already exists, else make the lock
	lock, _ := serviceClientLocker.LoadOrStore(doc.GetDuid(), &sync.RWMutex{})
	// Lock
	lock.(*sync.RWMutex).Lock()
	// Unlock before the function exits
	defer lock.(*sync.RWMutex).Unlock()

	log.Printf("[INFO] Connecting to mongodb://hwscmongodb duid: %s - uuid: %s\n",
		doc.GetDuid(), doc.GetUuid())
	client, err := mongo.NewClient(mongoServerDBWriter)
	if err != nil {
		log.Printf("[ERROR] Creating MongoDB client: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := client.Connect(context.TODO()); err != nil {
		log.Printf("[ERROR] Connecting MongoDB client: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	collection := client.Database(mongoDB).Collection(mongoDBCollection)

	filter := bson.NewDocument(
		bson.EC.String("duid", doc.GetDuid()),
		bson.EC.String("uuid", doc.GetUuid()),
	)
	result := collection.FindOneAndDelete(context.Background(), filter)

	// Extract the deleted MongoDB document
	if result == nil {
		log.Printf("[INFO] Document not found, duid: %s - uuid: %s\n",
			doc.GetDuid(), doc.GetUuid())

		return nil, status.Errorf(codes.InvalidArgument,
			"Document not found, duid: %s - uuid: %s",
			doc.GetDuid(), doc.GetUuid())
	}

	document := &pb.Document{}
	if err := result.Decode(document); err != nil {
		log.Printf("[ERROR] Document not found, duid: %s - uuid: %s - err: %s\n",
			doc.GetDuid(), doc.GetUuid(), err.Error())

		return nil, status.Errorf(codes.InvalidArgument,
			"Document not found, duid: %s - uuid: %s",
			doc.GetDuid(), doc.GetUuid())
	}

	if err := ValidateDocument(document); err != nil {
		log.Printf("[ERROR] Success deleting document, duid: %s - uuid: %s with validation error: %s\n",
			doc.GetDuid(), doc.GetUuid(), err.Error())
		log.Printf("[ERROR] Suspected document: \n%s\n\n", doc)
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Internal)},
			Message: "Deleted document with validation error",
			Data:    document,
		}, nil
	}

	log.Printf("[DEBUG] Deleted document: \n%s\n\n", document)

	if err := client.Disconnect(context.TODO()); err != nil {
		log.Printf("[ERROR] Success deleting document, duid: %s - uuid: %s with disconnection error: %s\n",
			doc.GetDuid(), doc.GetUuid(), err.Error())
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Internal)},
			Message: "Deleted document with MongoDB disconnection error",
			Data:    document,
		}, nil
	}

	// Log duid and uuid used for query
	log.Printf("[INFO] Success deleting document, duid: %s - uuid: %s\n",
		doc.GetDuid(), doc.GetUuid())

	return &pb.DocumentResponse{
		Status:  &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
		Data:    document,
	}, nil

}

// AddFileMetadata adds a new FileMetadata in a MongoDB document using a given url, UUID and DUID.
// Returns the updated Document.
// TODO
func (s Service) AddFileMetadata(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil

}

// DeleteFileMetadata deletes a FileMetadata in a MongoDB document using a given FUID, UUID and DUID.
// Returns the updated Document.
//TODO
func (s Service) DeleteFileMetadata(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil

}

// ListDistinctFieldValues list all the unique fields values required for the drop-down filter
// Returns the QueryTransaction.
//TODO
func (s Service) ListDistinctFieldValues(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil

}

// QueryDocument queries the MongoDB server with the given query parameters.
// Returns a collection of Documents.
//TODO
func (s Service) QueryDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {

	return nil, nil

}
