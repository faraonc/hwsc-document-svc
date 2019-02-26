# hwsc-document-svc

## Purpose
 Provides services to hwsc-app-gateway-svc for CRUD document, and file metadata in MongoDB

## Contract
The proto file and compiled proto buffers are located in [hwsc-api-blocks](https://github.com/hwsc-org/hwsc-api-blocks/tree/master/int/hwsc-document-svc/document) and [hwsc-api-blocks lib](https://github.com/hwsc-org/hwsc-api-blocks/tree/master/lib).
### GetStatus
- Gets the current status of the service.
### CreateDocument
- Creates a document in MongoDB.
- Returns the Document.
### ListUserDocumentCollection
- Retrieves all the MongoDB documents for a specific user with the given UUID.
- Returns a collection of Documents.
### UpdateDocument
- (completely) Updates a MongoDB document using DUID and UUID.
- Returns the updated Document.
### DeleteDocument
- Deletes a MongoDB document using DUID.
- Returns the deleted Document.
### AddFileMetadata
- Adds a new FileMetadata in a MongoDB document using a given url, DUID.
- Returns the updated Document.
### DeleteFileMetadata
- Deletes a FileMetadata in a MongoDB document using a given FUID, DUID.
- Returns the updated Document.
### ListDistinctFieldValues
- Retrieves all the unique fields values required for the front-end drop-down filter.
- Returns the QueryTransaction
### QueryDocument
- Queries the MongoDB server with the given query parameters.
- Returns a collection of Documents.

## Prerequisites
- GoLang version [go 1.11.5](https://golang.org/dl/)
- GoLang Modules [go mod](https://github.com/golang/go/wiki/Modules)
- Go Source Code Linter [golint](https://github.com/golang/lint)
- mongo-go-driver beta [1.0.0-rc1](https://github.com/mongodb/mongo-go-driver)
- Docker
- [Optional] If a new proto file and compiled proto buffer exists in [hwsc-api-blocks](https://github.com/hwsc-org/hwsc-api-blocks/tree/master/int/hwsc-document-svc/document), update dependency `$ go get -u <package name>`

## How to Run without Docker Container
1. Refer to [hwsc-dev-ops](https://github.com/hwsc-org/hwsc-dev-ops) to run DB locally
2. Grab prod/dev/test config file from Slack
3. Run main `$ go run main.go`

## How to Run with Docker Container
1. Refer to [hwsc-dev-ops](https://github.com/hwsc-org/hwsc-dev-ops) to run DB locally
2. `$ generate_container.sh`
3. Find your image `$ docker images`
4. Acquire `env.list` configuration
5. `$ docker run --env-file ./env.list -it -p 50051:50051 <imagename>`
6. Optional: run `$ docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' <container id>` to find the proper address, and update `env.list`

## How to Run Unit Test
1. Place DB migration source codes to be tested in `test_fixtures/mongodb/`
2. `$ cd service`
3. The unit test will programmatically run the test container and DB migration as required
4. For command-line summary, `$ go test -cover -v -failfast -race`
5. For comprehensive summary, `$ bash unit_test.sh`
6. If applicable, copy and push the new DB migration source codes in [hwsc-dev-ops](https://github.com/hwsc-org/hwsc-dev-ops) under `test` for integration testing

## How to Run Integration Test
- Refer to [hwsc-dev-ops](https://github.com/hwsc-org/hwsc-dev-ops) for running integration test
- Ensure the service is tested using the built container in DockerHub

