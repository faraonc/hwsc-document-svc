# hwsc-document-svc

## Purpose
 Provides services to hwsc-app-gateway-svc for CRUD document, and file metadata in Azure CosmosDB

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
- GoLang version [go 1.11.1](https://golang.org/dl/)
- GoLang Dependency Management [dep](https://github.com/golang/dep)
- Go Source Code Linter [golint](https://github.com/golang/lint)

## How to Run
1. Install dependencies and generate vendor folder ``$ dep ensure -v``
2. Run main ``$ go run main.go``
3. [Optional] If a new proto file and compiled proto buffer exists in [hwsc-api-blocks](https://github.com/hwsc-org/hwsc-api-blocks/tree/master/int/hwsc-document-svc/proto), update dependency ``$dep ensure -update``

## How to Unit Test
1. ``$ cd service``
2. For command-line summary, ``$ go test -cover -v``
3. For comprehensive summary, ``$ bash unit_test.sh``

