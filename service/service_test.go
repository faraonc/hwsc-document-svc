package service

import (
	"encoding/json"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pbsvc "github.com/hwsc-org/hwsc-api-blocks/int/hwsc-document-svc/document"
	pbdoc "github.com/hwsc-org/hwsc-api-blocks/lib"
	"github.com/hwsc-org/hwsc-document-svc/conf"
	"github.com/hwsc-org/hwsc-document-svc/consts"
	"github.com/hwsc-org/hwsc-lib/logger"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	tempDUID            string
	tempUUID            string
	tempFileFUID        string
	tempAudioFUID       string
	tempImageFUID       string
	tempVideoFUID       string
	imaginaryDUID       = randStringBytes(27)
	imaginaryUUID       = randStringBytes(26)
	randFirstName       = randStringBytes(10)
	randLastName        = randStringBytes(12)
	randCity            = randStringBytes(13)
	randProvince        = randStringBytes(13)
	numDocumentFixtures = 32
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func TestMain(t *testing.M) {
	logger.Info(consts.TestTag, "Initializing Test, this should ONLY print during unit tests")
	pool, err := dockertest.NewPool("")
	if err != nil {
		logger.Fatal(consts.TestTag, err.Error())
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("hwsc/test-hwsc-document-svc-mongodb", "latest",
		[]string{
			"MONGO_INITDB_DATABASE=admin",
			"MONGO_INITDB_ROOT_USERNAME=mongoadmin",
			"MONGO_INITDB_ROOT_PASSWORD=secret",
		},
	)
	if err != nil {
		logger.Fatal(consts.TestTag, err.Error())
	}
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		conf.DocumentDB.Reader = fmt.Sprintf("mongodb://testDocumentReader:testDocumentPwd@localhost:%s/test-document",
			resource.GetPort("27017/tcp"))
		conf.DocumentDB.Writer = fmt.Sprintf("mongodb://testDocumentWriter:testDocumentPwd@localhost:%s/test-document",
			resource.GetPort("27017/tcp"))
		var err error
		err = refreshMongoDBConnection(mongoDBWriter, &conf.DocumentDB.Writer)
		return err
	}); err != nil {
		logger.Fatal(consts.TestTag, err.Error())
	}
	m, err := migrate.New("file://test_fixtures/mongodb", conf.DocumentDB.Writer)
	if err != nil {
		logger.Fatal(consts.TestTag, err.Error())
	}
	if err := m.Up(); err != nil {
		logger.Fatal(consts.TestTag, err.Error())
	}

	// Open our jsonFile
	jsonFile, err := ioutil.ReadFile("test_fixtures/test-document.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		logger.Fatal(consts.TestTag, err.Error())
	}
	logger.Info(consts.TestTag, "test_fixtures/test-document.json")
	// defer the closing of our jsonFile so that we can parse it later on

	var docFixtures []*pbdoc.Document
	err = json.Unmarshal(jsonFile, &docFixtures)
	if err != nil {
		logger.Fatal(consts.TestTag, err.Error())
	}
	logger.Info(consts.TestTag, "Documents unmarshaled #", strconv.Itoa(len(docFixtures)))

	var count int
	for _, doc := range docFixtures {
		collection := mongoDBWriter.Database(conf.DocumentDB.Name).Collection(conf.DocumentDB.Collection)
		_, err := collection.InsertOne(context.Background(), doc)
		if err != nil {
			logger.Fatal(consts.TestTag, err.Error(), doc.String())
		}
		count++
	}
	logger.Info(consts.TestTag, "Documents inserted #", strconv.Itoa(count))
	if count != numDocumentFixtures {
		logger.Fatal(consts.TestTag, "failed setting up document fixture")
	}
	code := t.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		logger.Fatal(consts.TestTag, err.Error())
	}

	os.Exit(code)
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestGetStatus(t *testing.T) {
	cases := []struct {
		req         *pbsvc.DocumentRequest
		serverState state
		expMsg      string
	}{
		{&pbsvc.DocumentRequest{}, available, "OK"},
		{&pbsvc.DocumentRequest{}, unavailable, "Unavailable"},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, _ := s.GetStatus(context.TODO(), c.req)
		assert.Equal(t, c.expMsg, res.GetMessage())
	}
}

func TestCreateDocument(t *testing.T) {
	cases := []struct {
		req         *pbsvc.DocumentRequest
		serverState state
		expMsg      string
		isExpErr    bool
	}{
		{
			&pbsvc.DocumentRequest{}, unavailable, "rpc error: code = Unavailable desc = service unavailable", true,
		},
		{
			nil, available, "rpc error: code = InvalidArgument desc = nil request", true,
		},
		{
			&pbsvc.DocumentRequest{}, available, "rpc error: code = InvalidArgument desc = nil request data", true,
		},
		{
			&pbsvc.DocumentRequest{Data: &pbdoc.Document{Uuid: "garbage"}}, available,
			fmt.Sprintf("rpc error: code = InvalidArgument desc = %s", consts.ErrInvalidDocumentUUID.Error()), true,
		},
		{
			&pbsvc.DocumentRequest{
				Data: &pbdoc.Document{
					Duid: "",
					Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
					PublisherName: &pbdoc.Publisher{
						LastName:  "Test LastName",
						FirstName: "Test FirstName",
					},
					CallTypeName: "some call type name",
					GroundType:   "some ground type",
					StudySite: &pbdoc.StudySite{
						City:    "Seattle",
						State:   "Washington",
						Country: "USA",
					},
					Ocean:           "Pacific Ocean",
					SensorType:      "some sensor type",
					SensorName:      "some sensor name",
					SamplingRate:    100,
					Latitude:        89.123,
					Longitude:       -100.123,
					ImageUrlsMap:    nil,
					AudioUrlsMap:    nil,
					VideoUrlsMap:    nil,
					FileUrlsMap:     nil,
					RecordTimestamp: 1514764800,
					CreateTimestamp: 1539831496,
					UpdateTimestamp: 0,
					IsPublic:        true,
				},
				ImageUrls: []string{
					"https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png",
					"https://hwscdevstorage.blob.core.windows.net/images/Rotating_earth_(large).gif"},
				AudioUrls: []string{
					"https://hwscdevstorage.blob.core.windows.net/audios/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.wav",
					"https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3"},
				VideoUrls: []string{
					"https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
					"https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4",
				},
				FileUrls: []string{
					"https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
					"https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4",
				},
			}, available, "OK", false,
		},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, err := s.CreateDocument(context.TODO(), c.req)
		if !c.isExpErr {
			assert.Equal(t, c.expMsg, res.GetMessage())
			tempDUID = res.Data.GetDuid()
			tempUUID = res.Data.GetUuid()
		} else {
			assert.Equal(t, c.expMsg, err.Error())
			assert.EqualError(t, err, c.expMsg)
		}

	}
}

func TestListUserDocumentCollection(t *testing.T) {
	cases := []struct {
		req         *pbsvc.DocumentRequest
		serverState state
		expLength   int
		expMsg      string
		isExpErr    bool
	}{
		{&pbsvc.DocumentRequest{}, unavailable, 0, "rpc error: code = Unavailable desc = service unavailable", true},
		{nil, available, 0, "rpc error: code = InvalidArgument desc = nil request", true},
		{&pbsvc.DocumentRequest{}, available, 0, "rpc error: code = InvalidArgument desc = nil request data", true},
		{
			&pbsvc.DocumentRequest{Data: &pbdoc.Document{Uuid: "garbage"}}, available, 0,
			fmt.Sprintf("rpc error: code = InvalidArgument desc = %s", consts.ErrInvalidDocumentUUID.Error()), true,
		},
		{&pbsvc.DocumentRequest{Data: &pbdoc.Document{Uuid: "0XXXXSNJG0MQJHBF4QX1EFD6Y3"}}, available, 7, "OK", false},
		{&pbsvc.DocumentRequest{Data: &pbdoc.Document{Uuid: "123XXSNJG0MQASDF4QFFFFD6Y3"}}, available, 8, "OK", false},
		{
			&pbsvc.DocumentRequest{Data: &pbdoc.Document{Uuid: "4ee30333-8ec8-45a4-ba94-5e22c4a686de"}}, available, 0,
			fmt.Sprintf("rpc error: code = InvalidArgument desc = %s", consts.ErrInvalidDocumentUUID.Error()), true,
		},
		{
			&pbsvc.DocumentRequest{Data: &pbdoc.Document{Uuid: "xxx0XSNJG0MQJHBF4QX1EFD6Y3"}}, available, 0,
			"rpc error: code = InvalidArgument desc = No document for uuid: xxx0XSNJG0MQJHBF4QX1EFD6Y3", true,
		},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, err := s.ListUserDocumentCollection(context.TODO(), c.req)
		if !c.isExpErr {
			assert.Equal(t, c.expMsg, res.GetMessage())
			assert.NotEmpty(t, res.GetDocumentCollection())
			assert.Equal(t, c.expLength, len(res.GetDocumentCollection()))
		} else {
			assert.Equal(t, c.expMsg, err.Error())
			assert.EqualError(t, err, c.expMsg)
		}

	}
}

func TestUpdateDocument(t *testing.T) {
	cases := []struct {
		req         *pbsvc.DocumentRequest
		serverState state
		expMsg      string
		isExpErr    bool
	}{
		{&pbsvc.DocumentRequest{}, unavailable, "rpc error: code = Unavailable desc = service unavailable", true},
		{nil, available, "rpc error: code = InvalidArgument desc = nil request", true},
		{&pbsvc.DocumentRequest{}, available, "rpc error: code = InvalidArgument desc = nil request data", true},
		{
			&pbsvc.DocumentRequest{Data: &pbdoc.Document{Duid: ""}}, available,
			"rpc error: code = InvalidArgument desc = missing DUID", true,
		},
		{
			&pbsvc.DocumentRequest{
				Data: &pbdoc.Document{
					Duid: tempDUID,
					Uuid: tempUUID,
					PublisherName: &pbdoc.Publisher{
						LastName:  randFirstName,
						FirstName: randLastName,
					},
					CallTypeName: "some call type name",
					GroundType:   "some ground type",
					StudySite: &pbdoc.StudySite{
						City:     randCity,
						Province: randProvince,
						Country:  "Canada",
					},
					Ocean:        "Pacific Ocean",
					SensorType:   "some sensor type",
					SensorName:   "some sensor name",
					SamplingRate: 100,
					Latitude:     89.123,
					Longitude:    -100.123,
					ImageUrlsMap: map[string]string{
						"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png",
						"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwscdevstorage.blob.core.windows.net/images/Rotating_earth_(large).gif"},
					AudioUrlsMap: map[string]string{
						"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/audios/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.wav",
						"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3"},
					VideoUrlsMap: map[string]string{
						"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
						"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4"},
					FileUrlsMap: map[string]string{
						"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
						"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4"},
					RecordTimestamp: 1514764800,
					CreateTimestamp: 1539831496,
					UpdateTimestamp: 0,
					IsPublic:        false,
				},
				ImageUrls: []string{"https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png"},
				AudioUrls: []string{"https://hwscdevstorage.blob.core.windows.net/audios/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.wav"},
				VideoUrls: []string{"https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv"},
				FileUrls:  []string{"https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv"},
			}, available, "OK", false,
		},
		{
			&pbsvc.DocumentRequest{
				Data: &pbdoc.Document{
					Duid: imaginaryDUID,
					Uuid: imaginaryUUID,
					PublisherName: &pbdoc.Publisher{
						LastName:  randFirstName,
						FirstName: randLastName,
					},
					CallTypeName: "some call type name",
					GroundType:   "some ground type",
					StudySite: &pbdoc.StudySite{
						City:     randCity,
						Province: randProvince,
						Country:  "Canada",
					},
					Ocean:           "Pacific Ocean",
					SensorType:      "some sensor type",
					SensorName:      "some sensor name",
					SamplingRate:    100,
					Latitude:        89.123,
					Longitude:       -100.123,
					ImageUrlsMap:    nil,
					AudioUrlsMap:    nil,
					VideoUrlsMap:    nil,
					FileUrlsMap:     nil,
					RecordTimestamp: 1514764800,
					CreateTimestamp: 1539831496,
					UpdateTimestamp: 0,
					IsPublic:        false,
				},
				ImageUrls: []string{"https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png"},
				AudioUrls: []string{"https://hwscdevstorage.blob.core.windows.net/audios/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.wav"},
				VideoUrls: []string{"https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv"},
				FileUrls:  []string{"https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv"},
			}, available, fmt.Sprintf("rpc error: code = InvalidArgument desc = Document not found, duid: %s - uuid: %s",
				imaginaryDUID, imaginaryUUID), true,
		},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, err := s.UpdateDocument(context.TODO(), c.req)
		if !c.isExpErr {
			assert.Equal(t, c.expMsg, res.GetMessage())
		} else {
			assert.Equal(t, c.expMsg, err.Error())
			assert.EqualError(t, err, c.expMsg)
		}

	}
}

func TestDeleteDocument(t *testing.T) {
	cases := []struct {
		req         *pbsvc.DocumentRequest
		serverState state
		expMsg      string
		isExpErr    bool
	}{
		{&pbsvc.DocumentRequest{}, unavailable, "rpc error: code = Unavailable desc = service unavailable", true},
		{nil, available, "rpc error: code = InvalidArgument desc = nil request", true},
		{&pbsvc.DocumentRequest{}, available, "rpc error: code = InvalidArgument desc = nil request data", true},
		{&pbsvc.DocumentRequest{Data: &pbdoc.Document{}}, available, "rpc error: code = InvalidArgument desc = missing DUID", true},
		{
			&pbsvc.DocumentRequest{Data: &pbdoc.Document{Duid: imaginaryDUID}}, available,
			fmt.Sprintf("rpc error: code = InvalidArgument desc = Document not found, duid: %s", imaginaryDUID),
			true,
		},
		{&pbsvc.DocumentRequest{Data: &pbdoc.Document{Duid: tempDUID}}, available, "OK", false},
		{
			&pbsvc.DocumentRequest{Data: &pbdoc.Document{Duid: tempDUID}}, available,
			fmt.Sprintf("rpc error: code = InvalidArgument desc = Document not found, duid: %s", tempDUID), true,
		},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, err := s.DeleteDocument(context.TODO(), c.req)
		if !c.isExpErr {
			assert.Equal(t, c.expMsg, res.GetMessage())
		} else {
			assert.Equal(t, c.expMsg, err.Error())
			assert.EqualError(t, err, c.expMsg)
		}

	}
}

func TestListDistinctFieldValues(t *testing.T) {
	cases := []struct {
		req         *pbsvc.DocumentRequest
		serverState state
		expMsg      string
		isExpErr    bool
		expNumDocs  int
	}{
		{&pbsvc.DocumentRequest{}, unavailable, "rpc error: code = Unavailable desc = service unavailable", true, 0},
		{nil, available, "rpc error: code = InvalidArgument desc = nil request", true, 0},
		{&pbsvc.DocumentRequest{}, available, "", false, 0},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, err := s.ListDistinctFieldValues(context.TODO(), c.req)
		if !c.isExpErr {
			assert.Nil(t, err)
			assert.Equal(t, 6, len(res.QueryResults.Publishers))
			assert.Equal(t, 20, len(res.QueryResults.StudySites))
			assert.Equal(t, 18, len(res.QueryResults.CallTypeNames))
			assert.Equal(t, 8, len(res.QueryResults.GroundTypes))
			assert.Equal(t, 9, len(res.QueryResults.SensorNames))
			assert.Equal(t, 6, len(res.QueryResults.SensorTypes))
		} else {
			assert.Equal(t, c.expMsg, err.Error())
			assert.EqualError(t, err, c.expMsg)
		}

	}
}

func TestQueryDocument(t *testing.T) {
	cases := []struct {
		req         *pbsvc.DocumentRequest
		serverState state
		expMsg      string
		isExpErr    bool
		expNumDocs  int
	}{
		{&pbsvc.DocumentRequest{}, unavailable, "rpc error: code = Unavailable desc = service unavailable", true, 0},
		{nil, available, "rpc error: code = InvalidArgument desc = nil request", true, 0},
		{&pbsvc.DocumentRequest{}, available, "rpc error: code = InvalidArgument desc = nil query arguments", true, 0},
		{
			&pbsvc.DocumentRequest{
				QueryParameters: &pbdoc.QueryTransaction{
					MinRecordTimestamp: minTimestamp,
					MaxRecordTimestamp: time.Now().UTC().Unix() - 1,
				},
			}, available, "OK", false, 32,
		},
		{
			&pbsvc.DocumentRequest{
				QueryParameters: &pbdoc.QueryTransaction{
					Publishers: []*pbdoc.Publisher{
						{
							LastName:  "Seger",
							FirstName: "Kerri",
						},
						{
							LastName:  "Abadi",
							FirstName: "Shima",
						},
					},
					MinRecordTimestamp: minTimestamp,
					MaxRecordTimestamp: time.Now().UTC().Unix() - 1,
				},
			}, available, "OK", false, 11,
		},
		{
			&pbsvc.DocumentRequest{
				QueryParameters: &pbdoc.QueryTransaction{
					Publishers: []*pbdoc.Publisher{
						{
							LastName:  "Seger",
							FirstName: "Kerri",
						},
					},
					CallTypeNames: []string{
						"Wookie",
					},
					MinRecordTimestamp: minTimestamp,
					MaxRecordTimestamp: time.Now().UTC().Unix() - 1,
				},
			}, available, "OK", false, 1,
		},
		{
			&pbsvc.DocumentRequest{
				QueryParameters: &pbdoc.QueryTransaction{
					MinRecordTimestamp: 1446744336,
					MaxRecordTimestamp: 1510287809,
				},
			}, available, "OK", false, 12,
		},
		{
			&pbsvc.DocumentRequest{
				QueryParameters: &pbdoc.QueryTransaction{
					MinRecordTimestamp: 0,
					MaxRecordTimestamp: 1510287809,
				},
			}, available,
			"rpc error: code = InvalidArgument desc = invalid Document RecordTimestamp", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				QueryParameters: &pbdoc.QueryTransaction{
					MinRecordTimestamp: 1446744336,
					MaxRecordTimestamp: 0,
				},
			}, available,
			"rpc error: code = InvalidArgument desc = invalid Document RecordTimestamp", true, 0,
		},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, err := s.QueryDocument(context.TODO(), c.req)
		if !c.isExpErr {
			assert.Nil(t, err)
			assert.Equal(t, c.expNumDocs, len(res.GetDocumentCollection()))
		} else {
			assert.Equal(t, c.expMsg, err.Error())
			assert.EqualError(t, err, c.expMsg)
		}

	}
}

func TestAddFileMetadata(t *testing.T) {
	cases := []struct {
		req         *pbsvc.DocumentRequest
		serverState state
		expMsg      string
		isExpErr    bool
		expNumDocs  int
	}{
		{&pbsvc.DocumentRequest{}, unavailable, "rpc error: code = Unavailable desc = service unavailable", true, 0},
		{nil, available, "rpc error: code = InvalidArgument desc = nil request", true, 0},
		{
			&pbsvc.DocumentRequest{}, available, "rpc error: code = InvalidArgument desc = invalid FileMetadataParameters",
			true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url: "",
				},
			}, available, "rpc error: code = InvalidArgument desc = invalid FileMetadataParameters", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:  "some url",
					Duid: "",
				},
			}, available, "rpc error: code = InvalidArgument desc = invalid FileMetadataParameters", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:  "some url",
					Duid: "some duid",
				},
			}, available, "rpc error: code = InvalidArgument desc = invalid Document duid", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:  "some url",
					Duid: "1ChHfmKs8GX7D1XVf61lwVdisWf",
				},
			}, available,
			"rpc error: code = InvalidArgument desc = unreachable URI", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:   "some url",
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_AUDIO,
				},
			}, available, "rpc error: code = InvalidArgument desc = invalid Document AudioURL", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:   "some url",
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_IMAGE,
				},
			}, available, "rpc error: code = InvalidArgument desc = invalid Document ImageURL", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:   "some url",
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_VIDEO,
				},
			}, available, "rpc error: code = InvalidArgument desc = invalid Document VideoURL", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:   "some url",
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: 4,
				},
			}, available, "rpc error: code = InvalidArgument desc = invalid media type", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:   "https://hwscdevstorage.blob.core.windows.net/images/pusheen.mp3",
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_AUDIO,
				},
			}, available, "rpc error: code = InvalidArgument desc = unreachable URI", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:   "https://hwscdevstorage.blob.core.windows.net/images/pusheen.jpg",
					Duid:  "xxxHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_IMAGE,
				},
			}, available, "rpc error: code = InvalidArgument desc = Document not found, duid: xxxHfmKs8GX7D1XVf61lwVdisWf", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:   "https://hwscdevstorage.blob.core.windows.net/images/pusheen.jpg",
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_IMAGE,
				},
			}, available, "OK", false, 3,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:   "https://hwscdevstorage.blob.core.windows.net/videos/pusheen.mp4",
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_VIDEO,
				},
			}, available, "OK", false, 3,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:   "https://hwscdevstorage.blob.core.windows.net/audios/pusheen.mp3",
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_AUDIO,
				},
			}, available, "OK", false, 3,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Url:   "https://hwscdevstorage.blob.core.windows.net/videos/pusheen.mp4",
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_FILE,
				},
			}, available, "OK", false, 3,
		},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, err := s.AddFileMetadata(context.TODO(), c.req)
		if !c.isExpErr {
			assert.Nil(t, err)
			assert.NotNil(t, res)
			switch c.req.FileMetadataParameters.GetMedia() {
			case pbdoc.FileType_FILE:
				for k, v := range res.Data.GetFileUrlsMap() {
					if v == c.req.FileMetadataParameters.GetUrl() {
						tempFileFUID = k
					}
				}
				if !assert.Equal(t, c.expNumDocs, len(res.GetData().GetFileUrlsMap())) {
					assert.Fail(t, c.req.FileMetadataParameters.GetUrl())
				}
			case pbdoc.FileType_AUDIO:
				for k, v := range res.Data.GetAudioUrlsMap() {
					if v == c.req.FileMetadataParameters.GetUrl() {
						tempAudioFUID = k
					}
				}
				if !assert.Equal(t, c.expNumDocs, len(res.GetData().GetAudioUrlsMap())) {
					assert.Fail(t, c.req.FileMetadataParameters.GetUrl())
				}
			case pbdoc.FileType_IMAGE:
				for k, v := range res.Data.GetImageUrlsMap() {
					if v == c.req.FileMetadataParameters.GetUrl() {
						tempImageFUID = k
					}
				}
				if !assert.Equal(t, c.expNumDocs, len(res.GetData().GetImageUrlsMap())) {
					assert.Fail(t, c.req.FileMetadataParameters.GetUrl())
				}
			case pbdoc.FileType_VIDEO:
				for k, v := range res.Data.GetVideoUrlsMap() {
					if v == c.req.FileMetadataParameters.GetUrl() {
						tempVideoFUID = k
					}
				}
				if !assert.Equal(t, c.expNumDocs, len(res.GetData().GetVideoUrlsMap())) {
					assert.Fail(t, c.req.FileMetadataParameters.GetUrl())
				}
			}
		} else {
			assert.Equal(t, c.expMsg, err.Error())
			assert.EqualError(t, err, c.expMsg)
		}

	}
}

func TestDeleteFileMetadata(t *testing.T) {
	cases := []struct {
		req         *pbsvc.DocumentRequest
		serverState state
		expMsg      string
		isExpErr    bool
		expNumDocs  int
	}{
		{
			&pbsvc.DocumentRequest{}, unavailable, "rpc error: code = Unavailable desc = service unavailable",
			true, 0,
		},
		{nil, available, "rpc error: code = InvalidArgument desc = nil request", true, 0},
		{
			&pbsvc.DocumentRequest{}, available,
			"rpc error: code = InvalidArgument desc = invalid FileMetadataParameters", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Duid: "",
				},
			}, available, "rpc error: code = InvalidArgument desc = invalid FileMetadataParameters", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Duid: "some duid",
				},
			}, available, "rpc error: code = InvalidArgument desc = invalid Document duid", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Duid: "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Fuid: "some fuid",
				},
			}, available, "rpc error: code = InvalidArgument desc = invalid Document fuid", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Fuid:  tempFileFUID,
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: 4,
				},
			}, available, "rpc error: code = InvalidArgument desc = invalid media type", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Fuid:  tempFileFUID,
					Duid:  "xxxHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_FILE,
				},
			}, available, "rpc error: code = InvalidArgument desc = Document not found, duid: xxxHfmKs8GX7D1XVf61lwVdisWf", true, 0,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_IMAGE,
					Fuid:  tempImageFUID,
				},
			}, available, "OK", false, 2,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_VIDEO,
					Fuid:  tempVideoFUID,
				},
			}, available, "OK", false, 2,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_AUDIO,
					Fuid:  tempAudioFUID,
				},
			}, available, "OK", false, 2,
		},
		{
			&pbsvc.DocumentRequest{
				FileMetadataParameters: &pbdoc.FileMetadataTransaction{
					Duid:  "1ChHfmKs8GX7D1XVf61lwVdisWf",
					Media: pbdoc.FileType_FILE,
					Fuid:  tempFileFUID,
				},
			}, available, "OK", false, 2,
		},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, err := s.DeleteFileMetadata(context.TODO(), c.req)
		if !c.isExpErr {
			assert.Nil(t, err)
			assert.NotNil(t, res)
			switch c.req.FileMetadataParameters.GetMedia() {
			case pbdoc.FileType_FILE:
				if !assert.Equal(t, c.expNumDocs, len(res.GetData().GetFileUrlsMap())) {
					assert.Fail(t, c.req.FileMetadataParameters.GetFuid())
				}
			case pbdoc.FileType_AUDIO:
				if !assert.Equal(t, c.expNumDocs, len(res.GetData().GetAudioUrlsMap())) {
					assert.Fail(t, c.req.FileMetadataParameters.GetFuid())
				}
			case pbdoc.FileType_IMAGE:
				if !assert.Equal(t, c.expNumDocs, len(res.GetData().GetImageUrlsMap())) {
					assert.Fail(t, c.req.FileMetadataParameters.GetFuid())
				}
			case pbdoc.FileType_VIDEO:
				if !assert.Equal(t, c.expNumDocs, len(res.GetData().GetVideoUrlsMap())) {
					assert.Fail(t, c.req.FileMetadataParameters.GetFuid())
				}
			}
		} else {
			assert.Equal(t, c.expMsg, err.Error())
			assert.EqualError(t, err, c.expMsg)
		}

	}
}
