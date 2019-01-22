package service

import (
	"github.com/google/uuid"
	pb "github.com/hwsc-org/hwsc-api-blocks/int/hwsc-document-svc/proto"
	"github.com/hwsc-org/hwsc-document-svc/conf"
	"github.com/kylelemons/godebug/pretty"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/segmentio/ksuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"reflect"
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

// fuidLocker synchronizes the generating of fuid
type fuidLocker struct {
	lock sync.Mutex
}

// Service implements services for managing document
type Service struct{}

const (
	// available - Service is ready and available
	available state = 0

	// unavailable - Service is unavailable. Example: Provisioning something
	unavailable state = 1
)

var (
	serviceStateLocker stateLocker
	duidGenerator      duidLocker
	fuidGenerator      fuidLocker

	// Converts State of the service to a string
	serviceStateMap map[state]string

	// Stores the lock for each duid
	duidClientLocker sync.Map
)

func init() {
	serviceStateLocker = stateLocker{currentServiceState: available}
	duidGenerator = duidLocker{}
	fuidGenerator = fuidLocker{}

	serviceStateMap = map[state]string{
		available:   "Available",
		unavailable: "Unavailable",
	}
}

func (s state) String() string {
	return serviceStateMap[s]
}

// NewDUID generates a new document unique ID.
func (d *duidLocker) NewDUID() string {
	d.lock.Lock()
	defer d.lock.Unlock()
	duid := ksuid.New().String()
	return duid
}

// NewFUID generates a new file metadata unique ID.
func (d *fuidLocker) NewFUID() string {
	d.lock.Lock()
	defer d.lock.Unlock()
	newUUID := uuid.New().String()
	return newUUID
}

// GetStatus gets the current status of the service.
func (s *Service) GetStatus(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
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

	// Check MongoDB Clients
	if err := refreshMongoDBConnection(mongoDBReader, &conf.DocumentDB.Reader); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return &pb.DocumentResponse{
			Status:  &pb.DocumentResponse_Code{Code: uint32(codes.Unavailable)},
			Message: codes.Unavailable.String(),
		}, nil
	}
	if err := refreshMongoDBConnection(mongoDBWriter, &conf.DocumentDB.Writer); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
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
func (s *Service) CreateDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting CreateDocument service")

	if ok := isStateAvailable(); !ok {
		log.Printf("[INFO] %s\n", errServiceUnavailable.Error())
		return nil, status.Error(codes.Unavailable, errServiceUnavailable.Error())
	}

	if err := refreshMongoDBConnection(mongoDBWriter, &conf.DocumentDB.Writer); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req == nil {
		log.Printf("[ERROR] %s\n", errNilRequest.Error())
		return nil, status.Error(codes.InvalidArgument, errNilRequest.Error())
	}

	doc := req.GetData()
	if doc == nil {
		log.Printf("[ERROR] %s\n", errNilRequestData.Error())
		return nil, status.Error(codes.InvalidArgument, errNilRequestData.Error())
	}

	doc.Duid = duidGenerator.NewDUID()

	// Get the specific lock if it already exists, else make the lock
	lock, _ := duidClientLocker.LoadOrStore(doc.GetDuid(), &sync.RWMutex{})
	// Lock
	lock.(*sync.RWMutex).Lock()
	// Unlock before the function exits
	defer lock.(*sync.RWMutex).Unlock()

	// Extract image URLS
	if doc.GetImageUrlsMap() == nil {
		doc.ImageUrlsMap = make(map[string]string)
	}
	if req.GetImageUrls() != nil {
		for _, url := range req.GetImageUrls() {
			doc.ImageUrlsMap[fuidGenerator.NewFUID()] = url
		}
	}

	// Extract audio URLS
	if doc.GetAudioUrlsMap() == nil {
		doc.AudioUrlsMap = make(map[string]string)
	}
	if req.GetAudioUrls() != nil {
		for _, url := range req.GetAudioUrls() {
			doc.AudioUrlsMap[fuidGenerator.NewFUID()] = url
		}
	}

	// Extract video URLS
	if doc.GetVideoUrlsMap() == nil {
		doc.VideoUrlsMap = make(map[string]string)
	}
	if req.GetVideoUrls() != nil {
		for _, url := range req.GetVideoUrls() {
			doc.VideoUrlsMap[fuidGenerator.NewFUID()] = url
		}
	}

	// Extract file URLS
	if doc.GetFileUrlsMap() == nil {
		doc.FileUrlsMap = make(map[string]string)
	}
	if req.GetFileUrls() != nil {
		for _, url := range req.GetFileUrls() {
			doc.FileUrlsMap[fuidGenerator.NewFUID()] = url
		}
	}

	doc.CreateTimestamp = time.Now().UTC().Unix()

	if err := ValidateDocument(doc); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("[INFO] Document contains:\n %s\n\n", pretty.Sprint(doc))
	collection := mongoDBWriter.Database(conf.DocumentDB.Name).Collection(conf.DocumentDB.Collection)

	res, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Printf("[ERROR] InsertOne: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
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
func (s *Service) ListUserDocumentCollection(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting ListUserDocumentCollection service")

	if ok := isStateAvailable(); !ok {
		log.Printf("[ERROR] %s\n", errServiceUnavailable.Error())
		return nil, status.Error(codes.Unavailable, errServiceUnavailable.Error())
	}

	if err := refreshMongoDBConnection(mongoDBReader, &conf.DocumentDB.Reader); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req == nil {
		log.Printf("[ERROR] %s\n", errNilRequest.Error())
		return nil, status.Error(codes.InvalidArgument, errNilRequest.Error())
	}

	doc := req.GetData()
	if doc == nil {
		log.Printf("[ERROR] %s\n", errNilRequestData.Error())
		return nil, status.Error(codes.InvalidArgument, errNilRequestData.Error())
	}

	if err := ValidateUUID(doc.GetUuid()); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	collection := mongoDBReader.Database(conf.DocumentDB.Name).Collection(conf.DocumentDB.Collection)

	// Find all MongoDB documents for the specific uuid
	filter := bson.M{"uuid": doc.GetUuid()}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("[ERROR] Find: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

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

		documentCollection = append(documentCollection, document)
		log.Printf("[INFO] document: \n%s\n\n", pretty.Sprint(document))

	}

	if err := cur.Err(); err != nil {
		log.Printf("[ERROR] Cursor Err: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := cur.Close(context.Background()); err != nil {
		log.Printf("[ERROR] Cursor Err: %s\n", err.Error())
	}

	if len(documentCollection) == 0 {
		log.Printf("[ERROR] No document for uuid: %s\n", doc.GetUuid())
		return nil, status.Errorf(codes.InvalidArgument, "No document for uuid: %s", doc.GetUuid())
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
func (s *Service) UpdateDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting UpdateDocument service")

	if ok := isStateAvailable(); !ok {
		log.Printf("[ERROR] %s\n", errServiceUnavailable.Error())
		return nil, status.Error(codes.Unavailable, errServiceUnavailable.Error())
	}

	if err := refreshMongoDBConnection(mongoDBWriter, &conf.DocumentDB.Writer); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req == nil {
		log.Printf("[ERROR] %s\n", errNilRequest.Error())
		return nil, status.Error(codes.InvalidArgument, errNilRequest.Error())
	}

	doc := req.GetData()
	if doc == nil {
		log.Printf("[ERROR] %s\n", errNilRequestData.Error())
		return nil, status.Error(codes.InvalidArgument, errNilRequestData.Error())
	}

	if doc.GetDuid() == "" {
		log.Printf("[ERROR] Missing DUID")
		return nil, status.Error(codes.InvalidArgument, errMissingDUID.Error())
	}

	// Get the specific lock if it already exists, else make the lock
	lock, _ := duidClientLocker.LoadOrStore(doc.GetDuid(), &sync.RWMutex{})
	// Lock
	lock.(*sync.RWMutex).Lock()
	// Unlock before the function exits
	defer lock.(*sync.RWMutex).Unlock()

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

	doc.UpdateTimestamp = time.Now().UTC().Unix()

	if err := ValidateDocument(doc); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("[INFO] Document contains:\n %s\n\n", pretty.Sprint(doc))
	collection := mongoDBWriter.Database(conf.DocumentDB.Name).Collection(conf.DocumentDB.Collection)

	filter := bson.M{"duid": doc.GetDuid()}
	// option to return the the document after update
	after := options.After
	option := &options.FindOneAndReplaceOptions{ReturnDocument: &after}
	result := collection.FindOneAndReplace(context.Background(), filter, doc, option)

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

	log.Printf("[INFO] Updated document: \n%s\n\n", pretty.Sprint(document))
	log.Printf("[INFO] Success updating document, duid: %s - uuid: %s\n",
		doc.GetDuid(), doc.GetUuid())

	return &pb.DocumentResponse{
		Status:  &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
		Data:    document,
	}, nil

}

// DeleteDocument deletes a MongoDB document using DUID.
// Returns the deleted Document.
func (s *Service) DeleteDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting DeleteDocument service")

	if ok := isStateAvailable(); !ok {
		log.Printf("[ERROR] %s\n", errServiceUnavailable.Error())
		return nil, status.Error(codes.Unavailable, errServiceUnavailable.Error())
	}

	if err := refreshMongoDBConnection(mongoDBWriter, &conf.DocumentDB.Writer); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req == nil {
		log.Printf("[ERROR] %s\n", errNilRequest.Error())
		return nil, status.Error(codes.InvalidArgument, errNilRequest.Error())
	}

	doc := req.GetData()
	if doc == nil {
		log.Printf("[ERROR] %s\n", errNilRequestData.Error())
		return nil, status.Error(codes.InvalidArgument, errNilRequestData.Error())
	}

	if doc.GetDuid() == "" {
		log.Printf("[ERROR] Missing DUID")
		return nil, status.Error(codes.InvalidArgument, errMissingDUID.Error())
	}

	if err := ValidateDUID(doc.GetDuid()); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Get the specific lock if it already exists, else make the lock
	lock, _ := duidClientLocker.LoadOrStore(doc.GetDuid(), &sync.RWMutex{})
	// Lock
	lock.(*sync.RWMutex).Lock()
	// Unlock before the function exits
	defer lock.(*sync.RWMutex).Unlock()

	collection := mongoDBWriter.Database(conf.DocumentDB.Name).Collection(conf.DocumentDB.Collection)

	filter := bson.M{"duid": doc.GetDuid()}
	result := collection.FindOneAndDelete(context.Background(), filter)

	// Extract the deleted MongoDB document
	if result == nil {
		log.Printf("[INFO] Document not found, duid: %s\n", doc.GetDuid())

		return nil, status.Errorf(codes.InvalidArgument,
			"Document not found, duid: %s", doc.GetDuid())
	}

	document := &pb.Document{}
	if err := result.Decode(document); err != nil {
		log.Printf("[ERROR] Document not found, duid: %s - err: %s\n", doc.GetDuid(), err.Error())

		return nil, status.Errorf(codes.InvalidArgument,
			"Document not found, duid: %s", doc.GetDuid())
	}

	log.Printf("[INFO] Deleted document: \n%s\n\n", pretty.Sprint(document))
	// Log duid and uuid used for query
	log.Printf("[INFO] Success deleting document, duid: %s - uuid: %s\n",
		document.GetDuid(), document.GetUuid())

	return &pb.DocumentResponse{
		Status:  &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
		Data:    document,
	}, nil
}

// AddFileMetadata adds a new FileMetadata in a MongoDB document using a given url, and DUID.
// Returns the updated Document.
func (s *Service) AddFileMetadata(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting AddFileMetadata service")

	if ok := isStateAvailable(); !ok {
		log.Printf("[ERROR] %s\n", errServiceUnavailable.Error())
		return nil, status.Error(codes.Unavailable, errServiceUnavailable.Error())
	}

	if err := refreshMongoDBConnection(mongoDBWriter, &conf.DocumentDB.Writer); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req == nil {
		log.Printf("[ERROR] %s\n", errNilRequest.Error())
		return nil, status.Error(codes.InvalidArgument, errNilRequest.Error())
	}

	fileMetadataParameters := req.GetFileMetadataParameters()
	if fileMetadataParameters == nil || fileMetadataParameters.GetUrl() == "" ||
		fileMetadataParameters.GetDuid() == "" {

		log.Printf("[ERROR] %s\n", errInvalidFileMetadataParameters.Error())
		return nil, status.Error(codes.InvalidArgument, errInvalidFileMetadataParameters.Error())
	}

	if err := ValidateDUID(fileMetadataParameters.GetDuid()); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	switch fileMetadataParameters.Media {
	case pb.FileType_FILE:
		break
	case pb.FileType_AUDIO:
		if !audioRegex.MatchString(fileMetadataParameters.GetUrl()) {
			log.Printf("[ERROR] %s\n", errInvalidDocumentAudioURL.Error())
			return nil, status.Error(codes.InvalidArgument, errInvalidDocumentAudioURL.Error())
		}
	case pb.FileType_IMAGE:
		if !imageRegex.MatchString(fileMetadataParameters.GetUrl()) {
			log.Printf("[ERROR] %s\n", errInvalidDocumentImageURL.Error())
			return nil, status.Error(codes.InvalidArgument, errInvalidDocumentImageURL.Error())
		}
	case pb.FileType_VIDEO:
		if !videoRegex.MatchString(fileMetadataParameters.GetUrl()) {
			log.Printf("[ERROR] %s\n", errInvalidDocumentVideoURL.Error())
			return nil, status.Error(codes.InvalidArgument, errInvalidDocumentVideoURL.Error())
		}
	default:
		return nil, status.Error(codes.InvalidArgument, errMediaType.Error())
	}

	// Test if the URI is reachable
	if err := ValidateURL(fileMetadataParameters.GetUrl()); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("[INFO] FileMetadataParameters: \n%v\n\n", pretty.Sprint(req.GetFileMetadataParameters()))

	// Get the specific lock if it already exists, else make the lock
	lock, _ := duidClientLocker.LoadOrStore(fileMetadataParameters.GetDuid(), &sync.RWMutex{})
	// Lock
	lock.(*sync.RWMutex).Lock()
	// Unlock before the function exits
	defer lock.(*sync.RWMutex).Unlock()

	collection := mongoDBWriter.Database(conf.DocumentDB.Name).Collection(conf.DocumentDB.Collection)

	filter := bson.M{"duid": fileMetadataParameters.GetDuid()}
	bsonResult := collection.FindOne(context.Background(), filter)
	if bsonResult == nil {
		log.Printf("[ERROR] FindOne: %s\n", errNoDocumentFound.Error())
		return nil, status.Error(codes.InvalidArgument, errNoDocumentFound.Error())
	}
	documentToUpdate := &pb.Document{}
	if err := bsonResult.Decode(documentToUpdate); err != nil {
		log.Printf("[ERROR] Document not found, duid: %s - err: %s\n",
			fileMetadataParameters.GetDuid(), err.Error())

		return nil, status.Errorf(codes.InvalidArgument, "Document not found, duid: %s",
			fileMetadataParameters.GetDuid())
	}

	log.Printf("[INFO] Document to update: \n%s\n\n", pretty.Sprint(documentToUpdate))
	newFuid := fuidGenerator.NewFUID()
	switch fileMetadataParameters.Media {
	case pb.FileType_FILE:
		documentToUpdate.GetFileUrlsMap()[newFuid] = fileMetadataParameters.GetUrl()
	case pb.FileType_AUDIO:
		documentToUpdate.GetAudioUrlsMap()[newFuid] = fileMetadataParameters.GetUrl()
	case pb.FileType_IMAGE:
		documentToUpdate.GetImageUrlsMap()[newFuid] = fileMetadataParameters.GetUrl()
	case pb.FileType_VIDEO:
		documentToUpdate.GetVideoUrlsMap()[newFuid] = fileMetadataParameters.GetUrl()
	default:
		return nil, status.Error(codes.InvalidArgument, errMediaType.Error())
	}
	documentToUpdate.UpdateTimestamp = time.Now().UTC().Unix()

	// option to return the the document after update
	after := options.After
	option := &options.FindOneAndReplaceOptions{ReturnDocument: &after}
	result := collection.FindOneAndReplace(context.Background(), filter, documentToUpdate, option)

	// Extract the updated MongoDB document
	if result == nil {
		log.Printf("[ERROR] Extracting updated document, duid: %s\n", documentToUpdate.GetDuid())

		return nil, status.Errorf(codes.Internal,
			"Extracting updated document duid: %s", documentToUpdate.GetDuid())
	}

	document := &pb.Document{}
	if err := result.Decode(document); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	log.Printf("[INFO] Updated document: \n%s\n\n", pretty.Sprint(document))
	log.Printf("[INFO] Success adding file metadata in document, duid: %s - fuid: %s\n",
		document.GetDuid(), newFuid)

	return &pb.DocumentResponse{
		Status:  &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
		Data:    document,
	}, nil
}

// DeleteFileMetadata deletes a FileMetadata in a MongoDB document using a given FUID, and DUID.
// Returns the updated Document.
func (s *Service) DeleteFileMetadata(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting DeleteFileMetadata service")

	if ok := isStateAvailable(); !ok {
		log.Printf("[ERROR] %s\n", errServiceUnavailable.Error())
		return nil, status.Error(codes.Unavailable, errServiceUnavailable.Error())
	}

	if err := refreshMongoDBConnection(mongoDBWriter, &conf.DocumentDB.Writer); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req == nil {
		log.Printf("[ERROR] %s\n", errNilRequest.Error())
		return nil, status.Error(codes.InvalidArgument, errNilRequest.Error())
	}

	fileMetadataParameters := req.GetFileMetadataParameters()
	if fileMetadataParameters == nil || fileMetadataParameters.GetDuid() == "" {

		log.Printf("[ERROR] %s\n", errInvalidFileMetadataParameters.Error())
		return nil, status.Error(codes.InvalidArgument, errInvalidFileMetadataParameters.Error())
	}

	if err := ValidateDUID(fileMetadataParameters.GetDuid()); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateFUID(fileMetadataParameters.GetFuid()); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if fileMetadataParameters.Media > pb.FileType_VIDEO {
		return nil, status.Error(codes.InvalidArgument, errMediaType.Error())
	}

	log.Printf("[INFO] FileMetadataParameters: \n%v\n\n", pretty.Sprint(req.GetFileMetadataParameters()))

	// Get the specific lock if it already exists, else make the lock
	lock, _ := duidClientLocker.LoadOrStore(fileMetadataParameters.GetDuid(), &sync.RWMutex{})
	// Lock
	lock.(*sync.RWMutex).Lock()
	// Unlock before the function exits
	defer lock.(*sync.RWMutex).Unlock()

	collection := mongoDBWriter.Database(conf.DocumentDB.Name).Collection(conf.DocumentDB.Collection)

	filter := bson.M{"duid": fileMetadataParameters.GetDuid()}
	bsonResult := collection.FindOne(context.Background(), filter)
	if bsonResult == nil {
		log.Printf("[ERROR] FindOne: %s\n", errNoDocumentFound.Error())
		return nil, status.Error(codes.InvalidArgument, errNoDocumentFound.Error())
	}
	documentToUpdate := &pb.Document{}
	if err := bsonResult.Decode(documentToUpdate); err != nil {
		log.Printf("[ERROR] Document not found, duid: %s - err: %s\n",
			fileMetadataParameters.GetDuid(), err.Error())

		return nil, status.Errorf(codes.InvalidArgument,
			"Document not found, duid: %s", fileMetadataParameters.GetDuid())
	}

	log.Printf("[INFO] Document to update: \n%s\n\n", pretty.Sprint(documentToUpdate))

	switch fileMetadataParameters.Media {
	case pb.FileType_FILE:
		delete(documentToUpdate.GetFileUrlsMap(), fileMetadataParameters.GetFuid())
	case pb.FileType_AUDIO:
		delete(documentToUpdate.GetAudioUrlsMap(), fileMetadataParameters.GetFuid())
	case pb.FileType_IMAGE:
		delete(documentToUpdate.GetImageUrlsMap(), fileMetadataParameters.GetFuid())
	case pb.FileType_VIDEO:
		delete(documentToUpdate.GetVideoUrlsMap(), fileMetadataParameters.GetFuid())
	default:
		return nil, status.Error(codes.InvalidArgument, errMediaType.Error())
	}
	documentToUpdate.UpdateTimestamp = time.Now().UTC().Unix()

	// option to return the the document after update
	after := options.After
	option := &options.FindOneAndReplaceOptions{ReturnDocument: &after}
	result := collection.FindOneAndReplace(context.Background(), filter, documentToUpdate, option)

	// Extract the updated MongoDB document
	if result == nil {
		log.Printf("[ERROR] Extracting updated document, duid: %s\n", documentToUpdate.GetDuid())

		return nil, status.Errorf(codes.Internal,
			"Extracting updated document duid: %s", documentToUpdate.GetDuid())
	}

	document := &pb.Document{}
	if err := result.Decode(document); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	log.Printf("[INFO] Updated document: \n%s\n\n", pretty.Sprint(document))
	log.Printf("[INFO] Success deleting file metadata in document, duid: %s - fuid: %s\n",
		document.GetDuid(), fileMetadataParameters.GetFuid())

	return &pb.DocumentResponse{
		Status:  &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message: codes.OK.String(),
		Data:    document,
	}, nil

}

// ListDistinctFieldValues list all the unique fields values required for the front-end drop-down filter
// Returns the QueryTransaction.
func (s *Service) ListDistinctFieldValues(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting ListDistinctFieldValues service")
	if ok := isStateAvailable(); !ok {
		log.Printf("[ERROR] %s\n", errServiceUnavailable.Error())
		return nil, status.Error(codes.Unavailable, errServiceUnavailable.Error())
	}

	if err := refreshMongoDBConnection(mongoDBReader, &conf.DocumentDB.Reader); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req == nil {
		log.Printf("[ERROR] %s\n", errNilRequest.Error())
		return nil, status.Error(codes.InvalidArgument, errNilRequest.Error())
	}

	collection := mongoDBReader.Database(conf.DocumentDB.Name).Collection(conf.DocumentDB.Collection)

	// Get distinct using field names in distinctSearchFieldNames
	distinctResult := make([][]interface{}, len(distinctSearchFieldNames))
	for i := 0; i < len(distinctSearchFieldNames); i++ {
		doc := &bson.D{}
		result, err := collection.Distinct(context.Background(), distinctSearchFieldNames[i], doc)
		if err != nil {
			log.Printf("[ERROR] Distinct: %s\n", err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
		distinctResult[i] = result
	}

	// Extract distinct from distinctResult, and put them in queryResult
	queryResult := &pb.QueryTransaction{}
	val := reflect.ValueOf(*queryResult)
	for i := 0; i < len(distinctResultFieldIndices); i++ {
		fieldName := val.Type().Field(i).Name
		if err := extractDistinctResults(queryResult,
			fieldName, distinctResult[distinctResultFieldIndices[fieldName]]); err != nil {

			log.Printf("[ERROR] %s\n", err.Error())
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	log.Printf("[INFO] Distinct values: \n%v\n\n", pretty.Sprint(queryResult))
	log.Println("[INFO] Success listing distinct field values")
	return &pb.DocumentResponse{
		Status:       &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message:      codes.OK.String(),
		QueryResults: queryResult,
	}, nil

}

// QueryDocument queries the MongoDB server with the given query parameters.
// Returns a collection of Documents.
func (s *Service) QueryDocument(ctx context.Context, req *pb.DocumentRequest) (*pb.DocumentResponse, error) {
	log.Println("[INFO] Requesting QueryDocument service")
	if ok := isStateAvailable(); !ok {
		log.Printf("[ERROR] %s\n", errServiceUnavailable.Error())
		return nil, status.Error(codes.Unavailable, errServiceUnavailable.Error())
	}

	if err := refreshMongoDBConnection(mongoDBReader, &conf.DocumentDB.Reader); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if req == nil {
		log.Printf("[ERROR] %s\n", errNilRequest.Error())
		return nil, status.Error(codes.InvalidArgument, errNilRequest.Error())
	}

	queryParams := req.GetQueryParameters()
	if queryParams == nil {
		log.Printf("[ERROR] %s\n", errNilQueryArgs.Error())
		return nil, status.Error(codes.InvalidArgument, errNilQueryArgs.Error())
	}

	if err := ValidateRecordTimestamp(queryParams.MinRecordTimestamp); err != nil {
		log.Printf("[ERROR] %s\n", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := ValidateRecordTimestamp(queryParams.MaxRecordTimestamp); err != nil {
		log.Printf("[ERROR] %s\n", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("[INFO] QueryParameters contains:\n %s\n", pretty.Sprint(queryParams))
	collection := mongoDBReader.Database(conf.DocumentDB.Name).Collection(conf.DocumentDB.Collection)

	pipeline, err := buildAggregatePipeline(queryParams)
	if err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	cur, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Printf("[ERROR] Aggregate: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

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

		documentCollection = append(documentCollection, document)
		log.Printf("[INFO] document: \n%s\n\n", pretty.Sprint(document))

	}

	if err := cur.Err(); err != nil {
		log.Printf("[ERROR] Cursor Err: %s\n", err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := cur.Close(context.Background()); err != nil {
		log.Printf("[ERROR] Cursor Err: %s\n", err.Error())
	}

	log.Println("[INFO] Success querying documents")
	return &pb.DocumentResponse{
		Status:             &pb.DocumentResponse_Code{Code: uint32(codes.OK)},
		Message:            codes.OK.String(),
		DocumentCollection: documentCollection,
	}, nil

}
