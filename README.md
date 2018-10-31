# hwsc-document-svc

## Purpose
 Provides services to hwsc-app-gateway-svc for CRUD document, and file metadata in Azure CosmosDB

## Contract
The proto file and compiled proto buffers are located in [hwsc-api-blocks](https://github.com/faraonc/hwsc-api-blocks/tree/master/int/hwsc-document-svc/proto).
### GetStatus
- Gets the current status of the service.
### CreateDocument
- Creates a Document in MongoDB.
### ListUserDocumentCollection
- Gets all the MongoDB documents for a specific user with the given UUID.
- Returns a collection of MongoDB documents.
### UpdateDocument
- Updates a MongoDB document using DUID.
- Returns the updated MongoDB document.
### DeleteDocument
- Deletes a MongoDB document using UUID and DUID.
- Returns the deleted MongoDB document.
## Prerequisites
- GoLang version [go 1.11.1](https://golang.org/dl/)
- GoLang Dependency Management [dep](https://github.com/golang/dep)
- Go Source Code Linter [golint](https://github.com/golang/lint)

## How to Run
1. Install dependencies and generate vendor folder ``$ dep ensure -v``
2. Run main ``go run main.go``
3. [Optional] If a new proto file and compiled proto buffer exists in [hwsc-api-blocks](https://github.com/faraonc/hwsc-api-blocks/tree/master/int/hwsc-document-svc/proto), update dependency ``$dep ensure -update``

