package service

import (
	"fmt"
	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-document-svc/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"math/rand"
	"testing"
)

var (
	tempDUID      string
	tempUUID      string
	imaginaryDUID string
	imaginaryUUID string
	randFirstName string
	randLastName  string
	randCity      string
	randProvince  string
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	randFirstName = randStringBytes(10)
	randLastName = randStringBytes(12)
	randCity = randStringBytes(13)
	randProvince = randStringBytes(13)
	imaginaryDUID = randStringBytes(27)
	imaginaryUUID = randStringBytes(26)
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
		req         *pb.DocumentRequest
		serverState state
		expMsg      string
	}{
		{&pb.DocumentRequest{}, available, "OK"},
		{&pb.DocumentRequest{}, unavailable, "Unavailable"},
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
		req         *pb.DocumentRequest
		serverState state
		expMsg      string
		isExpErr    bool
	}{
		{&pb.DocumentRequest{}, unavailable,
			"rpc error: code = Unavailable desc = Service unavailable", true},
		{nil, available,
			"rpc error: code = InvalidArgument desc = Nil request", true},
		{&pb.DocumentRequest{}, available,
			"rpc error: code = InvalidArgument desc = Nil request data", true},
		{&pb.DocumentRequest{Data: &pb.Document{
			Uuid: "garbage",
		}}, available,
			"rpc error: code = InvalidArgument desc = invalid Document uuid", true},
		{&pb.DocumentRequest{
			Data: &pb.Document{
				Duid: "",
				Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
				PublisherName: &pb.Publisher{
					LastName:  "Test LastName",
					FirstName: "Test FirstName",
				},
				CallTypeName: "some call type name",
				GroundType:   "some ground type",
				StudySite: &pb.StudySite{
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
		}, available,
			"OK", false},
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
		req         *pb.DocumentRequest
		serverState state
		expLength   int
		expMsg      string
		isExpErr    bool
	}{
		{&pb.DocumentRequest{}, unavailable, 0,
			"rpc error: code = Unavailable desc = Service unavailable", true},
		{nil, available, 0,
			"rpc error: code = InvalidArgument desc = Nil request", true},
		{&pb.DocumentRequest{}, available, 0,
			"rpc error: code = InvalidArgument desc = Nil request data", true},
		{&pb.DocumentRequest{Data: &pb.Document{
			Uuid: "garbage",
		}}, available, 0,
			"rpc error: code = InvalidArgument desc = invalid Document uuid", true},
		{&pb.DocumentRequest{Data: &pb.Document{
			Uuid: "0XXXXSNJG0MQJHBF4QX1EFD6Y3",
		}}, available, 7,
			"OK", false},
		{&pb.DocumentRequest{Data: &pb.Document{
			Uuid: "123XXSNJG0MQASDF4QFFFFD6Y3",
		}}, available, 8,
			"OK", false},
		{&pb.DocumentRequest{Data: &pb.Document{
			Uuid: "4ee30333-8ec8-45a4-ba94-5e22c4a686de",
		}}, available, 0,
			"rpc error: code = InvalidArgument desc = invalid Document uuid", true},
		{&pb.DocumentRequest{Data: &pb.Document{
			Uuid: "xxx0XSNJG0MQJHBF4QX1EFD6Y3",
		}}, available, 0,
			"rpc error: code = InvalidArgument desc = No document for uuid: xxx0XSNJG0MQJHBF4QX1EFD6Y3", true},
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
		req         *pb.DocumentRequest
		serverState state
		expMsg      string
		isExpErr    bool
	}{
		{&pb.DocumentRequest{}, unavailable,
			"rpc error: code = Unavailable desc = Service unavailable", true},
		{nil, available,
			"rpc error: code = InvalidArgument desc = Nil request", true},
		{&pb.DocumentRequest{}, available,
			"rpc error: code = InvalidArgument desc = Nil request data", true},
		{&pb.DocumentRequest{Data: &pb.Document{
			Duid: "",
		}}, available,
			"rpc error: code = InvalidArgument desc = Missing DUID", true},
		{&pb.DocumentRequest{
			Data: &pb.Document{
				Duid: tempDUID,
				Uuid: tempUUID,
				PublisherName: &pb.Publisher{
					LastName:  randFirstName,
					FirstName: randLastName,
				},
				CallTypeName: "some call type name",
				GroundType:   "some ground type",
				StudySite: &pb.StudySite{
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
					"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwscdevstorage.blob.core.windows.net/video/videoplayback.mp4"},
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
		}, available,
			"OK", false},
		{&pb.DocumentRequest{
			Data: &pb.Document{
				Duid: imaginaryDUID,
				Uuid: imaginaryUUID,
				PublisherName: &pb.Publisher{
					LastName:  randFirstName,
					FirstName: randLastName,
				},
				CallTypeName: "some call type name",
				GroundType:   "some ground type",
				StudySite: &pb.StudySite{
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
		}, available,
			fmt.Sprintf("rpc error: code = InvalidArgument desc = Document not found, duid: %s - uuid: %s", imaginaryDUID, imaginaryUUID),
			true},
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
		req         *pb.DocumentRequest
		serverState state
		expMsg      string
		isExpErr    bool
	}{
		{&pb.DocumentRequest{}, unavailable,
			"rpc error: code = Unavailable desc = Service unavailable", true},
		{nil, available,
			"rpc error: code = InvalidArgument desc = Nil request", true},
		{&pb.DocumentRequest{}, available,
			"rpc error: code = InvalidArgument desc = Nil request data", true},
		{&pb.DocumentRequest{Data: &pb.Document{}}, available,
			"rpc error: code = InvalidArgument desc = Missing DUID", true},
		{&pb.DocumentRequest{Data: &pb.Document{
			Duid: "1CMjlaqHYNJhnVvWiGus3EiOno8",
		}}, available,
			"rpc error: code = InvalidArgument desc = invalid Document uuid", true},
		{&pb.DocumentRequest{Data: &pb.Document{
			Duid: imaginaryDUID,
			Uuid: imaginaryUUID,
		}}, available,
			fmt.Sprintf("rpc error: code = InvalidArgument desc = Document not found, duid: %s - uuid: %s", imaginaryDUID, imaginaryUUID),
			true},
		{&pb.DocumentRequest{Data: &pb.Document{
			Duid: imaginaryUUID,
			Uuid: imaginaryUUID,
		}}, available,
			"rpc error: code = InvalidArgument desc = invalid Document duid",
			true},
		{&pb.DocumentRequest{Data: &pb.Document{
			Duid: imaginaryDUID,
			Uuid: "",
		}}, available,
			"rpc error: code = InvalidArgument desc = invalid Document uuid",
			true},
		{&pb.DocumentRequest{Data: &pb.Document{
			Duid: tempDUID,
			Uuid: tempUUID,
		}}, available,
			"OK", false},
		{&pb.DocumentRequest{Data: &pb.Document{
			Duid: tempDUID,
			Uuid: tempUUID,
		}}, available,
			fmt.Sprintf("rpc error: code = InvalidArgument desc = Document not found, duid: %s - uuid: %s", tempDUID, tempUUID),
			true},
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
