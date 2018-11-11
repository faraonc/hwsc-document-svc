package service

import (
	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-document-svc/proto"
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
			true, "requires at least 1 valid Document ImageURL or AudioURL"},
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
			true, "invalid Document fuid"},
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
			true, "invalid Document duid"},
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
			true, "invalid Document uuid"},
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
			true, "invalid Document LastName"},
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
			true, "invalid Document FirstName"},
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
			true, "invalid Document CallTypeName"},
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
			true, "invalid Document GroundType"},
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
			true, "invalid Document City"},
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
			true, "invalid Document State"},
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
			true, "invalid Document Province"},
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
			true, "invalid Document Country"},
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
			true, "invalid Document Ocean"},
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
			true, "invalid Document SensorType"},
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
			true, "invalid Document SensorName"},
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
			true, "invalid Document SamplingRate"},
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
			true, "invalid Document Latitude"},
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
			true, "invalid Document Longitude"},
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
			true, "invalid Document ImageURL"},
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
			true, "invalid Document AudioURL"},
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
			true, "invalid Document VideoURL"},
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
			true, "invalid Document FileURL"},
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
			true, "invalid Document RecordTimestamp"},
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
			true, "invalid Document CreateTimestamp"},
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
			true, "invalid Document UpdateTimestamp"},
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
		{"0ujsszwN8NRY24YaXiTIE2VWDTSD", true, "invalid Document duid"},
		{"0ujsszwN8NRY24YaXiTIE2VWDT", true, "invalid Document duid"},
		{"", false, ""},
		{"   0ujsszwN8NRY24YaXiTIE2VWDTS", true, "invalid Document duid"},
		{"0ujsszwN8NRY24YaXiTIE2VWDTS    ", true, "invalid Document duid"},
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
		{"0000XSNJG0MQJHBF4QX1EFD6Y33", true, "invalid Document uuid"},
		{"0000XSNJG0MQJHBF4QX1EFD6Y", true, "invalid Document uuid"},
		{"", true, "invalid Document uuid"},
		{"   0000XSNJG0MQJHBF4QX1EFD6Y3", true, "invalid Document uuid"},
		{"0000XSNJG0MQJHBF4QX1EFD6Y3    ", true, "invalid Document uuid"},
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
		{"4fg30392-8ec8-45a4-ba94-5e22c4a686de", true, "invalid Document fuid"},
		{"a4ff30392-8ec8-45a4-ba94-5e22c4a686de", true, "invalid Document fuid"},
		{"4ff30392-a8ec8-45a4-ba94-5e22c4a686de", true, "invalid Document fuid"},
		{"4ff30392-8ec8-a45a4-ba94-5e22c4a686de", true, "invalid Document fuid"},
		{"4ff30392-8ec8-45a4-aba94-5e22c4a686de", true, "invalid Document fuid"},
		{"4ff30392-8ec8-45a4-ba94-a5e22c4a686de", true, "invalid Document fuid"},
		{"", true, "invalid Document fuid"},
		{"   4ff30392-8ec8-45a4-ba94-5e22c4a686de", true, "invalid Document fuid"},
		{"4ff30392-8ec8-45a4-ba94-5e22c4a686de    ", true, "invalid Document fuid"},
		{"4ff303928ec8-45a4-ba94-5e22c4a686de", true, "invalid Document fuid"},
		{"4ff30392-8ec845a4-ba94-5e22c4a686de", true, "invalid Document fuid"},
		{"4ff30392-8ec8-45a4ba94-5e22c4a686de", true, "invalid Document fuid"},
		{"4ff30392-8ec8-45a4-ba945e22c4a686de", true, "invalid Document fuid"},
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
		{"", "Lisa", true, "invalid Document LastName"},
		{"Kim", "123456789123456789012345678901234", true, "invalid Document FirstName"},
		{"       ", "Leesa", true, "invalid Document LastName"},
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
		{"", true, "invalid Document LastName"},
		{"123456789123456789012345678901234", true, "invalid Document LastName"},
		{"       ", true, "invalid Document LastName"},
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
		{"", true, "invalid Document FirstName"},
		{"123456789123456789012345678901234", true, "invalid Document FirstName"},
		{"       ", true, "invalid Document FirstName"},
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
		{"", true, "invalid Document CallTypeName"},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, "invalid Document CallTypeName"},
		{"       ", true, "invalid Document CallTypeName"},
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
		{"", true, "invalid Document GroundType"},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, "invalid Document GroundType"},
		{"       ", true, "invalid Document GroundType"},
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
		{"Batangas City", "", "12345678912345678901234567890123412345678912345678901234567890123", "Philippines", true, "invalid Document Province"},
		{"Tijuana", "12345678912345678901234567890123412345678912345678901234567890123", "", "Mexico", true, "invalid Document State"},
		{"", "Baja California", "", "Mexico", true, "invalid Document City"},
		{"Tijuana", "Baja California", "", "", true, "invalid Document Country"},
		{"San Diego", "CA", "", "12345678912345678901234567890123412345678912345678901234567890123", true, "invalid Document Country"},
		{"       ", "Baja California", "", "Mexico", true, "invalid Document City"},
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
		{"", true, "invalid Document City"},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, "invalid Document City"},
		{"       ", true, "invalid Document City"},
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
		{"12345678912345678901234567890123412345678912345678901234567890123", true, "invalid Document State"},
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
		{"12345678912345678901234567890123412345678912345678901234567890123", true, "invalid Document Province"},
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
		{"", true, "invalid Document Country"},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, "invalid Document Country"},
		{"       ", true, "invalid Document Country"},
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

		{"Atlantic ocean    hello ", true, "invalid Document Ocean"},
		{"Atlantic oceans", true, "invalid Document Ocean"},
		{"", true, "invalid Document Ocean"},
		{"      ", true, "invalid Document Ocean"},
		{"idonotexist", true, "invalid Document Ocean"},
		{"Indian 1 Ocean", true, "invalid Document Ocean"},
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
		{"", true, "invalid Document SensorType"},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, "invalid Document SensorType"},
		{"       ", true, "invalid Document SensorType"},
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
		{"", true, "invalid Document SensorName"},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, "invalid Document SensorName"},
		{"       ", true, "invalid Document SensorName"},
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
		{maxSamplingRate + 1, true, "invalid Document SamplingRate"},
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
		{minLatitude - 1, true, "invalid Document Latitude"},
		{minLatitude, false, ""},
		{0, false, ""},
		{45, false, ""},
		{maxLatitude, false, ""},
		{maxLatitude + 1, true, "invalid Document Latitude"},
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
		{minLongitude - 1, true, "invalid Document Longitude"},
		{minLongitude, false, ""},
		{0, false, ""},
		{150, false, ""},
		{maxLongitude, false, ""},
		{maxLongitude + 1, true, "invalid Document Longitude"},
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
		}, true, "invalid Document fuid"},
		{nil,
			true, "nil Document ImageURLs",
		},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/images/Rotating_earth_(large).gif",
		}, true, "invalid Document ImageURL"},
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
		}, true, "invalid Document fuid",
		},
		{nil,
			true, "nil Document AudioURLs",
		},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/audios/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.wav",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3",
		}, true, "invalid Document AudioURL"},
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
		}, true, "invalid Document fuid",
		},
		{nil,
			true, "nil Document VideoURLs",
		},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4",
		}, true, "invalid Document VideoURL"},
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
		}, true, "invalid Document fuid",
		},
		{nil,
			true, "nil Document FileURLs",
		},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4",
		}, true, "invalid Document FileURL"},
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
		{minTimestamp - 1, true, "invalid Document RecordTimestamp"},
		{time.Now().UTC().Unix() + 100, true, "invalid Document RecordTimestamp"},
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
		{1514764800, 1539831496, true, "invalid Document CreateTimestamp"},
		{time.Now().UTC().Unix() + 100, 1539831496, true, "invalid Document CreateTimestamp"},
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
		{1514764800, 1539831496, true, "invalid Document UpdateTimestamp"},
		{time.Now().UTC().Unix() + 100, 1539831496, true, "invalid Document UpdateTimestamp"},
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
