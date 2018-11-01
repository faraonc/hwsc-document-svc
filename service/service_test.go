package service

import (
	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-document-svc/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

var (
	tempDUID string
	tempUUID string
)

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
		{&pb.DocumentRequest{Data: &pb.Document{
			Duid:         "",
			Uuid:         "0000XSNJG0MQJHBF4QX1EFD6Y3",
			LastName:     "Faraon",
			FirstName:    "Conard",
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			Region:       "some region",
			Ocean:        "Pacific Ocean",
			SensorType:   "some sensor type",
			SensorName:   "some sensor name",
			SampleRate:   100,
			Latitude:     89.123,
			Longitude:    -100.123,
			ImageUrl: map[string]string{
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwssappstorage.blob.core.windows.net/image/hulkgif.png",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwssappstorage.blob.core.windows.net/image/Rotating_earth_(large).gif"},
			AudioUrl: map[string]string{
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwssappstorage.blob.core.windows.net/audio/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.wav",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwssappstorage.blob.core.windows.net/audio/Milad Hosseini - Deli Asheghetam [128].mp3"},
			VideoUrl: map[string]string{
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwssappstorage.blob.core.windows.net/video/videoplayback.wmv",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwssappstorage.blob.core.windows.net/video/videoplayback.mp4"},
			FileUrl: map[string]string{
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwssappstorage.blob.core.windows.net/video/videoplayback.wmv",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwssappstorage.blob.core.windows.net/video/videoplayback.mp4"},
			RecordTimestamp: 1514764800,
			CreateTimestamp: 1539831496,
			UpdateTimestamp: 0,
		}}, available,
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
		{&pb.DocumentRequest{Data: &pb.Document{
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
		}}, available,
			"OK", false},
		{&pb.DocumentRequest{Data: &pb.Document{
			Uuid: "4ee30333-8ec8-45a4-ba94-5e22c4a686de",
		}}, available,
			"rpc error: code = InvalidArgument desc = invalid Document uuid", true},
		{&pb.DocumentRequest{Data: &pb.Document{
			Uuid: "xxx0XSNJG0MQJHBF4QX1EFD6Y3",
		}}, available,
			"rpc error: code = InvalidArgument desc = No document for uuid: xxx0XSNJG0MQJHBF4QX1EFD6Y3", true},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, err := s.ListUserDocumentCollection(context.TODO(), c.req)
		if !c.isExpErr {
			assert.Equal(t, c.expMsg, res.GetMessage())
			assert.NotEmpty(t, res.DocumentCollection)
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
		{&pb.DocumentRequest{Data: &pb.Document{
			Duid:         "1CMjsoGz1cNOkIYaarbcSzmNg1n",
			Uuid:         "0000XSNJG0MQJHBF4QX1EFD6Y3",
			LastName:     "Keem",
			FirstName:    "Leesa",
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			Region:       "Venus",
			Ocean:        "Pacific Ocean",
			SensorType:   "some sensor type",
			SensorName:   "some sensor name",
			SampleRate:   100,
			Latitude:     89.123,
			Longitude:    -100.123,
			ImageUrl: map[string]string{
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwssappstorage.blob.core.windows.net/image/hulkgif.png",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwssappstorage.blob.core.windows.net/image/Rotating_earth_(large).gif"},
			AudioUrl: map[string]string{
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwssappstorage.blob.core.windows.net/audio/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.wav",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwssappstorage.blob.core.windows.net/audio/Milad Hosseini - Deli Asheghetam [128].mp3"},
			VideoUrl: map[string]string{
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwssappstorage.blob.core.windows.net/video/videoplayback.wmv",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwssappstorage.blob.core.windows.net/video/videoplayback.mp4"},
			FileUrl: map[string]string{
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwssappstorage.blob.core.windows.net/video/videoplayback.wmv",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwssappstorage.blob.core.windows.net/video/videoplayback.mp4"},
			RecordTimestamp: 1514764800,
			CreateTimestamp: 1539831496,
			UpdateTimestamp: 0,
		}}, available,
			"OK", false},
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
			Duid: "xCMjlaqHYNJhnVvWiGxxxEiOno8",
			Uuid: "1100XSNJG0MQJHBF4QX1EFD6Y3",
		}}, available,
			"rpc error: code = InvalidArgument desc = Document not found, duid: xCMjlaqHYNJhnVvWiGxxxEiOno8 - uuid: 1100XSNJG0MQJHBF4QX1EFD6Y3",
			true},
		{&pb.DocumentRequest{Data: &pb.Document{
			Duid: tempDUID,
			Uuid: tempUUID,
		}}, available,
			"OK", false},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		s := Service{}
		res, err := s.DeleteDocument(context.TODO(), c.req)
		if !c.isExpErr {
			assert.Equal(t, c.expMsg, res.GetMessage())
			//assert.NotEmpty(t, res.FileMetadataCollection)
		} else {
			assert.Equal(t, c.expMsg, err.Error())
			assert.EqualError(t, err, c.expMsg)
		}

	}
}
