package service

import (
	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-document-svc/proto"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValidateDocument(t *testing.T) {

	cases := []struct {
		input    *pb.Document
		isExpErr bool
		errorStr string
	}{
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
			IsPublic:        true,
		},
			false, ""},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
			},
			Ocean:        "Pacific Ocean",
			SensorType:   "some sensor type",
			SensorName:   "some sensor name",
			SamplingRate: 100,
			Latitude:     89.123,
			Longitude:    -100.123,
			ImageUrlsMap: map[string]string{},
			AudioUrlsMap: map[string]string{},
			VideoUrlsMap: map[string]string{
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4"},
			FileUrlsMap: map[string]string{
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4"},
			RecordTimestamp: 1514764800,
			CreateTimestamp: 1539831496,
			UpdateTimestamp: 0,
			IsPublic:        true,
		},
			true, errAtLeastOneImageAudioURL.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
			},
			Ocean:        "Pacific Ocean",
			SensorType:   "some sensor type",
			SensorName:   "some sensor name",
			SamplingRate: 100,
			Latitude:     89.123,
			Longitude:    -100.123,
			ImageUrlsMap: map[string]string{
				"4ff303392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df":  "https://hwscdevstorage.blob.core.windows.net/images/Rotating_earth_(large).gif"},
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
			IsPublic:        true,
		},
			true, errInvalidDocumentFUID.Error()},
		{&pb.Document{
			Duid: "0ujssszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentDUID.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "000s0XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentUUID.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentLastName.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentFirstName.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentCallTypeName.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentGroundType.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "",
				Country: "USA",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentCity.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				State:   "123456789012345678901234567890123",
				Country: "USA",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentState.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:     "Vancouver",
				Province: "1234567890123456789012345678901234567890123456789",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentProvince.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentCountry.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
			},
			Ocean:        "",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentOcean.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
			},
			Ocean:        "Pacific Ocean",
			SensorType:   "",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentSensorType.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
			},
			Ocean:        "Pacific Ocean",
			SensorType:   "some sensor type",
			SensorName:   "",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentSensorName.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
			},
			Ocean:        "Pacific Ocean",
			SensorType:   "some sensor type",
			SensorName:   "some sensor name",
			SamplingRate: maxSamplingRate + 1,
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
			IsPublic:        true,
		},
			true, errInvalidDocumentSamplingRate.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
			},
			Ocean:        "Pacific Ocean",
			SensorType:   "some sensor type",
			SensorName:   "some sensor name",
			SamplingRate: 100,
			Latitude:     99.123,
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
			IsPublic:        true,
		},
			true, errInvalidDocumentLatitude.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
			},
			Ocean:        "Pacific Ocean",
			SensorType:   "some sensor type",
			SensorName:   "some sensor name",
			SamplingRate: 100,
			Latitude:     89.123,
			Longitude:    -180.123,
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
			IsPublic:        true,
		},
			true, errInvalidDocumentLongitude.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
			},
			Ocean:        "Pacific Ocean",
			SensorType:   "some sensor type",
			SensorName:   "some sensor name",
			SamplingRate: 100,
			Latitude:     89.123,
			Longitude:    -100.123,
			ImageUrlsMap: map[string]string{
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentImageURL.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
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
			IsPublic:        true,
		},
			true, errInvalidDocumentAudioURL.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4"},
			FileUrlsMap: map[string]string{
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4"},
			RecordTimestamp: 1514764800,
			CreateTimestamp: 1539831496,
			UpdateTimestamp: 0,
			IsPublic:        true,
		},
			true, errInvalidDocumentVideoURL.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
				"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
				"4ff30392-8ec8-45a4-ba94-5e22c4a686df": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4"},
			RecordTimestamp: 1514764800,
			CreateTimestamp: 1539831496,
			UpdateTimestamp: 0,
			IsPublic:        true,
		},
			true, errInvalidDocumentFileURL.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
			RecordTimestamp: 0,
			CreateTimestamp: 1539831496,
			UpdateTimestamp: 0,
			IsPublic:        true,
		},
			true, errInvalidDocumentRecordTimestamp.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
			RecordTimestamp: 1539831497,
			CreateTimestamp: 1539831496,
			UpdateTimestamp: 0,
			IsPublic:        true,
		},
			true, errInvalidDocumentCreateTimestamp.Error()},
		{&pb.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pb.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pb.StudySite{
				City:    "Seattle",
				Country: "USA",
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
			UpdateTimestamp: 1539831495,
			IsPublic:        true,
		},
			true, errInvalidUpdateTimestamp.Error()},
	}

	for _, c := range cases {
		err := ValidateDocument(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
			if err == nil {
				t.Fatal(err)
			}
		} else {
			assert.Nil(t, err)
			if err != nil {
				t.Fatal(err)
			}
		}

	}

}

func TestValidateDUID(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"0ujsszwN8NRY24YaXiTIE2VWDTS", false, ""},
		{"0ujsszwN8NRY24YaXiTIE2VWDTSD", true, errInvalidDocumentDUID.Error()},
		{"0ujsszwN8NRY24YaXiTIE2VWDT", true, errInvalidDocumentDUID.Error()},
		{"", false, ""},
		{"   0ujsszwN8NRY24YaXiTIE2VWDTS", true, errInvalidDocumentDUID.Error()},
		{"0ujsszwN8NRY24YaXiTIE2VWDTS    ", true, errInvalidDocumentDUID.Error()},
	}

	for _, c := range cases {
		err := ValidateDUID(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}
}

func TestValidateUUID(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"0000XSNJG0MQJHBF4QX1EFD6Y3", false, ""},
		{"0000XSNJG0MQJHBF4QX1EFD6Y33", true, errInvalidDocumentUUID.Error()},
		{"0000XSNJG0MQJHBF4QX1EFD6Y", true, errInvalidDocumentUUID.Error()},
		{"", true, errInvalidDocumentUUID.Error()},
		{"   0000XSNJG0MQJHBF4QX1EFD6Y3", true, errInvalidDocumentUUID.Error()},
		{"0000XSNJG0MQJHBF4QX1EFD6Y3    ", true, errInvalidDocumentUUID.Error()},
	}

	for _, c := range cases {
		err := ValidateUUID(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}
}

func TestValidateFUID(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"4ff30392-8ec8-45a4-ba94-5e22c4a686de", false, ""},
		{"4fg30392-8ec8-45a4-ba94-5e22c4a686de", true, errInvalidDocumentFUID.Error()},
		{"a4ff30392-8ec8-45a4-ba94-5e22c4a686de", true, errInvalidDocumentFUID.Error()},
		{"4ff30392-a8ec8-45a4-ba94-5e22c4a686de", true, errInvalidDocumentFUID.Error()},
		{"4ff30392-8ec8-a45a4-ba94-5e22c4a686de", true, errInvalidDocumentFUID.Error()},
		{"4ff30392-8ec8-45a4-aba94-5e22c4a686de", true, errInvalidDocumentFUID.Error()},
		{"4ff30392-8ec8-45a4-ba94-a5e22c4a686de", true, errInvalidDocumentFUID.Error()},
		{"", true, errInvalidDocumentFUID.Error()},
		{"   4ff30392-8ec8-45a4-ba94-5e22c4a686de", true, errInvalidDocumentFUID.Error()},
		{"4ff30392-8ec8-45a4-ba94-5e22c4a686de    ", true, errInvalidDocumentFUID.Error()},
		{"4ff303928ec8-45a4-ba94-5e22c4a686de", true, errInvalidDocumentFUID.Error()},
		{"4ff30392-8ec845a4-ba94-5e22c4a686de", true, errInvalidDocumentFUID.Error()},
		{"4ff30392-8ec8-45a4ba94-5e22c4a686de", true, errInvalidDocumentFUID.Error()},
		{"4ff30392-8ec8-45a4-ba945e22c4a686de", true, errInvalidDocumentFUID.Error()},
	}

	for _, c := range cases {
		err := ValidateFUID(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidatePublisher(t *testing.T) {
	cases := []struct {
		lastName  string
		firstName string
		isExpErr  bool
		errorStr  string
	}{
		{"Kim", "Lisa", false, ""},
		{"", "Lisa", true, errInvalidDocumentLastName.Error()},
		{"Kim", "123456789123456789012345678901234", true, errInvalidDocumentFirstName.Error()},
		{"       ", "Leesa", true, errInvalidDocumentLastName.Error()},
	}

	for _, c := range cases {
		err := ValidatePublisher(c.lastName, c.firstName)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateLastName(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"Kim", false, ""},
		{"", true, errInvalidDocumentLastName.Error()},
		{"123456789123456789012345678901234", true, errInvalidDocumentLastName.Error()},
		{"       ", true, errInvalidDocumentLastName.Error()},
	}

	for _, c := range cases {
		err := ValidateLastName(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateFirstName(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"Lisa", false, ""},
		{"", true, errInvalidDocumentFirstName.Error()},
		{"123456789123456789012345678901234", true, errInvalidDocumentFirstName.Error()},
		{"       ", true, errInvalidDocumentFirstName.Error()},
	}

	for _, c := range cases {
		err := ValidateFirstName(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateCallTypeName(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"Gunshot", false, ""},
		{"", true, errInvalidDocumentCallTypeName.Error()},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, errInvalidDocumentCallTypeName.Error()},
		{"       ", true, errInvalidDocumentCallTypeName.Error()},
	}

	for _, c := range cases {
		err := ValidateCallTypeName(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateGroundType(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"Gunshot", false, ""},
		{"", true, errInvalidDocumentGroundType.Error()},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, errInvalidDocumentGroundType.Error()},
		{"       ", true, errInvalidDocumentGroundType.Error()},
	}

	for _, c := range cases {
		err := ValidateGroundType(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateStudySite(t *testing.T) {
	cases := []struct {
		city     string
		state    string
		province string
		country  string
		isExpErr bool
		errorStr string
	}{
		{"Tijuana", "Baja California", "", "Mexico", false, ""},
		{"Copenhagen", "", "", "Denmark", false, ""},
		{"Batangas City", "", "Batangas", "Philippines", false, ""},
		{"Batangas City", "", "12345678912345678901234567890123412345678912345678901234567890123", "Philippines", true, errInvalidDocumentProvince.Error()},
		{"Tijuana", "12345678912345678901234567890123412345678912345678901234567890123", "", "Mexico", true, errInvalidDocumentState.Error()},
		{"", "Baja California", "", "Mexico", true, errInvalidDocumentCity.Error()},
		{"Tijuana", "Baja California", "", "", true, errInvalidDocumentCountry.Error()},
		{"San Diego", "CA", "", "12345678912345678901234567890123412345678912345678901234567890123", true, errInvalidDocumentCountry.Error()},
		{"       ", "Baja California", "", "Mexico", true, errInvalidDocumentCity.Error()},
	}

	for _, c := range cases {
		err := ValidateStudySite(c.city, c.state, c.province, c.country)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateCity(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"Mexico City", false, ""},
		{"", true, errInvalidDocumentCity.Error()},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, errInvalidDocumentCity.Error()},
		{"       ", true, errInvalidDocumentCity.Error()},
	}

	for _, c := range cases {
		err := ValidateCity(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateState(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"California", false, ""},
		{"", false, ""},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, errInvalidDocumentState.Error()},
	}

	for _, c := range cases {
		err := ValidateState(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateProvince(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"Batangas City", false, ""},
		{"", false, ""},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, errInvalidDocumentProvince.Error()},
	}

	for _, c := range cases {
		err := ValidateProvince(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateCountry(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"Mexico", false, ""},
		{"", true, errInvalidDocumentCountry.Error()},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, errInvalidDocumentCountry.Error()},
		{"       ", true, errInvalidDocumentCountry.Error()},
	}

	for _, c := range cases {
		err := ValidateCountry(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateOcean(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"Pacific", false, ""},
		{"Atlantic", false, ""},
		{"Arctic", false, ""},
		{"Indian", false, ""},
		{"Southern", false, ""},

		{"Pacific Ocean", false, ""},
		{"Atlantic ocean", false, ""},
		{"Arctic oceaN", false, ""},
		{"INDIAN OCEAN", false, ""},
		{"SoutherN      OCEAN", false, ""},

		{"      Pacific Ocean", false, ""},
		{"Atlantic ocean     ", false, ""},

		{"Atlantic ocean    hello ", true, errInvalidDocumentOcean.Error()},
		{"Atlantic oceans", true, errInvalidDocumentOcean.Error()},
		{"", true, errInvalidDocumentOcean.Error()},
		{"      ", true, errInvalidDocumentOcean.Error()},
		{"idonotexist", true, errInvalidDocumentOcean.Error()},
		{"Indian 1 Ocean", true, errInvalidDocumentOcean.Error()},
	}

	for _, c := range cases {
		err := ValidateOcean(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateSensorType(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"Bprobe", false, ""},
		{"", true, errInvalidDocumentSensorType.Error()},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, errInvalidDocumentSensorType.Error()},
		{"       ", true, errInvalidDocumentSensorType.Error()},
	}

	for _, c := range cases {
		err := ValidateSensorType(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}
}

func TestValidateSensorName(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"Tag", false, ""},
		{"", true, errInvalidDocumentSensorName.Error()},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, errInvalidDocumentSensorName.Error()},
		{"       ", true, errInvalidDocumentSensorName.Error()},
	}

	for _, c := range cases {
		err := ValidateSensorName(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}
}

func TestValidateSamplingRate(t *testing.T) {
	cases := []struct {
		input    uint32
		isExpErr bool
		errorStr string
	}{
		{0, false, ""},
		{1000, false, ""},
		{maxSamplingRate, false, ""},
		{maxSamplingRate + 1, true, errInvalidDocumentSamplingRate.Error()},
	}

	for _, c := range cases {
		err := ValidateSamplingRate(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}
}

func TestValidateLatitude(t *testing.T) {
	cases := []struct {
		input    float32
		isExpErr bool
		errorStr string
	}{
		{minLatitude - 1, true, errInvalidDocumentLatitude.Error()},
		{minLatitude, false, ""},
		{0, false, ""},
		{45, false, ""},
		{maxLatitude, false, ""},
		{maxLatitude + 1, true, errInvalidDocumentLatitude.Error()},
	}

	for _, c := range cases {
		err := ValidateLatitude(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}
}

func TestValidateLongitude(t *testing.T) {
	cases := []struct {
		input    float32
		isExpErr bool
		errorStr string
	}{
		{minLongitude - 1, true, errInvalidDocumentLongitude.Error()},
		{minLongitude, false, ""},
		{0, false, ""},
		{150, false, ""},
		{maxLongitude, false, ""},
		{maxLongitude + 1, true, errInvalidDocumentLongitude.Error()},
	}

	for _, c := range cases {
		err := ValidateLongitude(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}
}

func TestValidateImageURLs(t *testing.T) {
	cases := []struct {
		input    map[string]string
		isExpErr bool
		errorStr string
	}{
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/images/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.jpg",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/imimagesage/hulkgif.png",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/images/Rotating_earth_(large).gif",
		}, false, ""},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba9a4-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/images/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.jpg",
		}, true, errInvalidDocumentFUID.Error()},
		{nil,
			true, errInvalidDocumentImageURLs.Error(),
		},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/images/Rotating_earth_(large).gif",
		}, true, errInvalidDocumentImageURL.Error()},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "hwscdevstorage.blob.core.windows.net/images/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.jpg",
		}, true, "invalid Document ImageURL: hwscdevstorage.blob.core.windows.net/images/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.jpg"},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3",
		}, true, "invalid Document image type ImageURL: https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3"},
		{map[string]string{}, false, ""},
	}

	for _, c := range cases {
		err := ValidateImageURLs(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestValidateAudioURLs(t *testing.T) {
	cases := []struct {
		input    map[string]string
		isExpErr bool
		errorStr string
	}{
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/audios/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.wav",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3",
		}, false, ""},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba9a4-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3",
		}, true, errInvalidDocumentFUID.Error(),
		},
		{nil,
			true, errInvalidDocumentAudioURLs.Error(),
		},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/audios/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.wav",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3",
		}, true, errInvalidDocumentAudioURL.Error()},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3",
		}, true, "invalid Document AudioURL: hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3"},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png",
		}, true, "invalid Document audio type AudioURL: https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png"},
		{map[string]string{}, false, ""},
	}

	for _, c := range cases {
		err := ValidateAudioURLs(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestValidateVideoURLs(t *testing.T) {
	cases := []struct {
		input    map[string]string
		isExpErr bool
		errorStr string
	}{
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4",
		}, false, ""},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba9a4-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4",
		}, true, errInvalidDocumentFUID.Error(),
		},
		{nil,
			true, errInvalidDocumentVideoURLs.Error(),
		},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4",
		}, true, errInvalidDocumentVideoURL.Error()},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
		}, true, "invalid Document VideoURL: hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv"},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png",
		}, true, "invalid Document video type VideoURL: https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png"},
		{map[string]string{}, false, ""},
	}

	for _, c := range cases {
		err := ValidateVideoURLs(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestValidateFileURLs(t *testing.T) {
	cases := []struct {
		input    map[string]string
		isExpErr bool
		errorStr string
	}{
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4",
		}, false, ""},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba9a4-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4",
		}, true, errInvalidDocumentFUID.Error(),
		},
		{nil,
			true, errInvalidDocumentFileURLs.Error(),
		},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4",
		}, true, errInvalidDocumentFileURL.Error()},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
		}, true, "invalid Document FileURL: hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv"},
		{map[string]string{}, false, ""},
	}

	for _, c := range cases {
		err := ValidateFileURLs(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestValidateRecordTimestamp(t *testing.T) {
	cases := []struct {
		input    int64
		isExpErr bool
		errorStr string
	}{
		{1514764800, false, ""},
		{minTimestamp, false, ""},
		{minTimestamp - 1, true, errInvalidDocumentRecordTimestamp.Error()},
		{time.Now().UTC().Unix() + 100, true, errInvalidDocumentRecordTimestamp.Error()},
	}

	for _, c := range cases {
		err := ValidateRecordTimestamp(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateCreateTimestamp(t *testing.T) {
	cases := []struct {
		inputCreateTimestamp int64
		inputRecordTimestamp int64
		isExpErr             bool
		errorStr             string
	}{
		{0, 1514764800, false, ""},
		{1539831496, 1514764800, false, ""},
		{1514764800, 1539831496, true,
			errInvalidDocumentCreateTimestamp.Error()},
		{time.Now().UTC().Unix() + 100, 1539831496, true,
			errInvalidDocumentCreateTimestamp.Error()},
	}

	for _, c := range cases {
		err := ValidateCreateTimestamp(c.inputCreateTimestamp, c.inputRecordTimestamp)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestValidateUpdateTimestamp(t *testing.T) {
	cases := []struct {
		inputUpdateTimestamp int64
		inputCreateTimestamp int64
		isExpErr             bool
		errorStr             string
	}{
		{0, 1514764800, false, ""},
		{1514764801, 1514764800, false, ""},
		{1514764800, 1539831496, true,
			errInvalidUpdateTimestamp.Error()},
		{time.Now().UTC().Unix() + 100, 1539831496, true,
			errInvalidUpdateTimestamp.Error()},
	}

	for _, c := range cases {
		err := ValidateUpdateTimestamp(c.inputUpdateTimestamp, c.inputCreateTimestamp)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}

	}

}

func TestIsStateAvailable(t *testing.T) {
	cases := []struct {
		serverState state
		expRet      bool
	}{
		{available, true},
		{unavailable, false},
	}

	for _, c := range cases {
		serviceStateLocker.currentServiceState = c.serverState
		assert.Equal(t, c.expRet, isStateAvailable())
	}

}

func TestBuildAggregatePipeline(t *testing.T) {
	cases := []struct {
		input     *pb.QueryTransaction
		expOutput *bson.Array
		isExpErr  bool
		errorStr  string
	}{
		{nil, nil, true, errNilQueryTransaction.Error()},
		{
			&pb.QueryTransaction{
				Publishers: []*pb.Publisher{
					{
						LastName:  "Seger",
						FirstName: "Kerri",
					},
					{
						LastName:  "Abadi",
						FirstName: "Shima",
					},
				},
				StudySites: []*pb.StudySite{
					{
						City:     "San Diego",
						State:    "California",
						Province: "",
						Country:  "USA",
					},
					{
						City:     "Batangas City",
						State:    "",
						Province: "Batangas",
						Country:  "Philippines",
					},
					{
						City:     "Some City",
						State:    "",
						Province: "",
						Country:  "Some Country",
					},
				},
				CallTypeNames: []string{},
				GroundTypes:   []string{"Wookie"},
				SensorTypes:   []string{"BProbe"},
				SensorNames:   []string{"Moto"},
			},
			bson.NewArray(
				bson.VC.DocumentFromElements(
					bson.EC.SubDocumentFromElements(
						"$match",
						bson.EC.SubDocumentFromElements("publisherName.lastName",
							bson.EC.ArrayFromElements("$in",
								bson.VC.String("Seger"),
								bson.VC.String("Abadi"),
							)),
						bson.EC.SubDocumentFromElements("publisherName.firstName",
							bson.EC.ArrayFromElements("$in",
								bson.VC.String("Kerri"),
								bson.VC.String("Shima"),
							)),
						bson.EC.SubDocumentFromElements("studySite.city",
							bson.EC.ArrayFromElements("$in",
								bson.VC.String("San Diego"),
								bson.VC.String("Batangas City"),
								bson.VC.String("Some City"),
							)),
						bson.EC.SubDocumentFromElements("studySite.state",
							bson.EC.ArrayFromElements("$in",
								bson.VC.String("California"),
							)),
						bson.EC.SubDocumentFromElements("studySite.province",
							bson.EC.ArrayFromElements("$in",
								bson.VC.String("Batangas"),
							)),
						bson.EC.SubDocumentFromElements("studySite.country",
							bson.EC.ArrayFromElements("$in",
								bson.VC.String("USA"),
								bson.VC.String("Philippines"),
								bson.VC.String("Some Country"),
							)),
						bson.EC.SubDocumentFromElements("callTypeName",
							bson.EC.ArrayFromElements("$in", bson.VC.Regex(".*", ""))),
						bson.EC.SubDocumentFromElements("groundType",
							bson.EC.ArrayFromElements("$in",
								bson.VC.String("Wookie"),
							)),
						bson.EC.SubDocumentFromElements("sensorType",
							bson.EC.ArrayFromElements("$in",
								bson.VC.String("BProbe"),
							)),
						bson.EC.SubDocumentFromElements("sensorName",
							bson.EC.ArrayFromElements("$in",
								bson.VC.String("Moto"),
							)),
					),
				),
			),
			false,
			"",
		},
	}

	for _, c := range cases {
		_, err := buildAggregatePipeline(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestBuildArrayFromElements(t *testing.T) {
	cases := []struct {
		input     []string
		expOutput *bson.Element
	}{
		{nil, bson.EC.ArrayFromElements("$in", bson.VC.Regex(".*", ""))},
		{[]string{}, bson.EC.ArrayFromElements("$in", bson.VC.Regex(".*", ""))},

		{[]string{
			"Seger",
			"Abadi",
		},
			bson.EC.ArrayFromElements("$in",
				bson.VC.String("Seger"),
				bson.VC.String("Abadi")),
		},
	}

	for _, c := range cases {
		elems := buildArrayFromElements(c.input)
		assert.True(t, elems.Equal(c.expOutput))
	}
}

func TestExtractPublishersFields(t *testing.T) {
	cases := []struct {
		input      []*pb.Publisher
		lastNames  []string
		firstNames []string
	}{
		{nil, []string{}, []string{}},
		{[]*pb.Publisher{}, []string{}, []string{}},
		{
			[]*pb.Publisher{
				{
					LastName:  "Seger",
					FirstName: "Kerri",
				},
				{
					LastName:  "Abadi",
					FirstName: "Shima",
				},
			},
			[]string{
				"Seger",
				"Abadi",
			},
			[]string{
				"Kerri",
				"Shima",
			},
		},
	}

	for _, c := range cases {
		lastNames, firstNames := extractPublishersFields(c.input)
		assert.True(t, areSliceEqual(c.lastNames, lastNames))
		assert.True(t, areSliceEqual(c.firstNames, firstNames))
	}
}

func TestExtractStudySitesFields(t *testing.T) {
	cases := []struct {
		input     []*pb.StudySite
		cities    []string
		states    []string
		provinces []string
		countries []string
	}{
		{nil, []string{}, []string{}, []string{}, []string{}},
		{[]*pb.StudySite{}, []string{}, []string{}, []string{}, []string{}},
		{
			[]*pb.StudySite{
				{
					City:     "San Diego",
					State:    "California",
					Province: "",
					Country:  "USA",
				},
				{
					City:     "Batangas City",
					State:    "",
					Province: "Batangas",
					Country:  "Philippines",
				},
				{
					City:     "Some City",
					State:    "",
					Province: "",
					Country:  "Some Country",
				},
			},
			[]string{
				"San Diego",
				"Batangas City",
				"Some City",
			},
			[]string{
				"California",
			},
			[]string{
				"Batangas",
			},
			[]string{
				"USA",
				"Philippines",
				"Some Country",
			},
		},
	}

	for _, c := range cases {
		cities, states, provinces, countries := extractStudySitesFields(c.input)
		assert.True(t, areSliceEqual(c.cities, cities))
		assert.True(t, areSliceEqual(c.states, states))
		assert.True(t, areSliceEqual(c.provinces, provinces))
		assert.True(t, areSliceEqual(c.countries, countries))
	}
}

func areSliceEqual(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
