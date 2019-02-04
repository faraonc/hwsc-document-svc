# hwsc-document-svc

## Purpose
 Provides services to hwsc-app-gateway-svc for CRUD document, and file metadata in MongoDB

## Contract
The proto file and compiled proto buffers are located in [hwsc-api-blocks](https://github.com/hwsc-org/hwsc-api-blocks/tree/master/int/hwsc-document-svc/proto).
### GetStatus
- Gets the current status of the service.
### CreateDocument
- Creates a document in MongoDB.
- Returns the Document.
### ListUserDocumentCollection
- Retrieves all the MongoDB documents for a specific user with the given UUID.
- Returns a collection of Documents.
### UpdateDocument
- (completely) Updates a MongoDB document using DUID.
- Returns the updated Document.
### DeleteDocument
- Deletes a MongoDB document using UUID and DUID.
- Returns the deleted Document.
### AddFileMetadata
- Adds a new FileMetadata in a MongoDB document using a given url, UUID and DUID.
- Returns the updated Document.
### DeleteFileMetadata
- Deletes a FileMetadata in a MongoDB document using a given FUID, UUID and DUID.
- Returns the updated Document.
### ListDistinctFieldValues
- Retrieves all the unique fields values required for the front-end drop-down filter.
- Returns the QueryTransaction
### QueryDocument
- Queries the MongoDB server with the given query parameters.
- Returns a collection of Documents.

## Prerequisites
- GoLang version [go 1.11.5](https://golang.org/dl/)
- GoLang Dependency Management [dep](https://github.com/golang/dep)
- Go Source Code Linter [golint](https://github.com/golang/lint)
- mongo-go-driver beta [0.3.0](https://github.com/mongodb/mongo-go-driver)
- Docker
- [Optional] If a new proto file and compiled proto buffer exists in [hwsc-api-blocks](https://github.com/hwsc-org/hwsc-api-blocks/tree/master/int/hwsc-document-svc/proto), update dependency ``$dep ensure -update``

## How to Run without Docker Container
1. Install dependencies and generate vendor folder `$ dep ensure -v`
2. Update ENV variables
3. Run main `$ go run main.go`

## How to Run with Docker Container
1. Install dependencies and generate vendor folder `$ dep ensure -v`
2. `$ generate_container.sh`
3. Find your image `$ docker images`
4. Acquire `env.list` configuration
5. `$ docker run --env-file ./env.list -it -p 50051:50051 <imagename>`

## How to Unit Test
1. `$ cd service`
2. For command-line summary, `$ go test -cover -v`
3. For comprehensive summary, `$ bash unit_test.sh`

