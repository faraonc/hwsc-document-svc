package service

import (
	pbdoc "github.com/hwsc-org/hwsc-api-blocks/lib"
	"github.com/hwsc-org/hwsc-document-svc/consts"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"testing"
	"time"
)

func TestNewDUID(t *testing.T) {
	const count = 100
	var tokens sync.Map

	var wg sync.WaitGroup // waits until all goroutines finish before executing code below wg.Wait()
	wg.Add(count)         // indicate we are going to wait for 100 go routines

	// start is a signal channel
	// channel of empty structs is used to indicate that this channel
	// will only be used for signalling and not for passing data
	start := make(chan struct{})
	for i := 0; i < count; i++ {
		go func() {
			// <-start blocks code below, waiting until the for loop is finished
			// it waits for the start channel to be closed,
			// once closed, all goroutines will execute(start) almost simultaneously
			<-start

			// decrement wg.Add, indicates 1 go routine has finished
			// defer will call wg.Done() at end of go routine
			defer wg.Done()

			// store tokens in map to check for duplicates
			duid := duidGenerator.NewDUID()
			_, ok := tokens.Load(duid)
			assert.Equal(t, false, ok)

			tokens.Store(duid, true)
		}()
	}

	// closing this channel, will unblock it,
	// allowing execution to continue
	close(start)

	// wait for all 100 go routines to finish (when wg.Add reaches 0)
	// blocks from running any code below it
	wg.Wait()
}

func TestNewFUID(t *testing.T) {
	const count = 100
	var tokens sync.Map

	var wg sync.WaitGroup // waits until all goroutines finish before executing code below wg.Wait()
	wg.Add(count)         // indicate we are going to wait for 100 go routines

	// start is a signal channel
	// channel of empty structs is used to indicate that this channel
	// will only be used for signalling and not for passing data
	start := make(chan struct{})
	for i := 0; i < count; i++ {
		go func() {
			// <-start blocks code below, waiting until the for loop is finished
			// it waits for the start channel to be closed,
			// once closed, all goroutines will execute(start) almost simultaneously
			<-start

			// decrement wg.Add, indicates 1 go routine has finished
			// defer will call wg.Done() at end of go routine
			defer wg.Done()

			// store tokens in map to check for duplicates
			fuid := fuidGenerator.NewFUID()
			_, ok := tokens.Load(fuid)
			assert.Equal(t, false, ok)

			tokens.Store(fuid, true)
		}()
	}

	// closing this channel, will unblock it,
	// allowing execution to continue
	close(start)

	// wait for all 100 go routines to finish (when wg.Add reaches 0)
	// blocks from running any code below it
	wg.Wait()
}

func TestValidateDocument(t *testing.T) {

	cases := []struct {
		input    *pbdoc.Document
		isExpErr bool
		errorStr string
	}{
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrAtLeastOneImageAudioURL.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentFUID.Error()},
		{&pbdoc.Document{
			Duid: "0ujssszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentDUID.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "000s0XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentUUID.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentLastName.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentFirstName.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentCallTypeName.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentGroundType.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentCity.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentState.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentProvince.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentCountry.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentOcean.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentSensorType.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentSensorName.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentSamplingRate.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentLatitude.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentLongitude.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentImageURL.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentAudioURL.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentVideoURL.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentFileURL.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentRecordTimestamp.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidDocumentCreateTimestamp.Error()},
		{&pbdoc.Document{
			Duid: "0ujsszwN8NRY24YaXiTIE2VWDTS",
			Uuid: "0000XSNJG0MQJHBF4QX1EFD6Y3",
			PublisherName: &pbdoc.Publisher{
				LastName:  "Kim",
				FirstName: "Lisa",
			},
			CallTypeName: "some call type name",
			GroundType:   "some ground type",
			StudySite: &pbdoc.StudySite{
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
			true, consts.ErrInvalidUpdateTimestamp.Error()},
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
		{"0ujsszwN8NRY24YaXiTIE2VWDTSD", true, consts.ErrInvalidDocumentDUID.Error()},
		{"0ujsszwN8NRY24YaXiTIE2VWDT", true, consts.ErrInvalidDocumentDUID.Error()},
		{"", false, ""},
		{"   0ujsszwN8NRY24YaXiTIE2VWDTS", true, consts.ErrInvalidDocumentDUID.Error()},
		{"0ujsszwN8NRY24YaXiTIE2VWDTS    ", true, consts.ErrInvalidDocumentDUID.Error()},
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
		{"0000XSNJG0MQJHBF4QX1EFD6Y33", true, consts.ErrInvalidDocumentUUID.Error()},
		{"0000XSNJG0MQJHBF4QX1EFD6Y", true, consts.ErrInvalidDocumentUUID.Error()},
		{"", true, consts.ErrInvalidDocumentUUID.Error()},
		{"   0000XSNJG0MQJHBF4QX1EFD6Y3", true, consts.ErrInvalidDocumentUUID.Error()},
		{"0000XSNJG0MQJHBF4QX1EFD6Y3    ", true, consts.ErrInvalidDocumentUUID.Error()},
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
		{"4fg30392-8ec8-45a4-ba94-5e22c4a686de", true, consts.ErrInvalidDocumentFUID.Error()},
		{"a4ff30392-8ec8-45a4-ba94-5e22c4a686de", true, consts.ErrInvalidDocumentFUID.Error()},
		{"4ff30392-a8ec8-45a4-ba94-5e22c4a686de", true, consts.ErrInvalidDocumentFUID.Error()},
		{"4ff30392-8ec8-a45a4-ba94-5e22c4a686de", true, consts.ErrInvalidDocumentFUID.Error()},
		{"4ff30392-8ec8-45a4-aba94-5e22c4a686de", true, consts.ErrInvalidDocumentFUID.Error()},
		{"4ff30392-8ec8-45a4-ba94-a5e22c4a686de", true, consts.ErrInvalidDocumentFUID.Error()},
		{"", true, consts.ErrInvalidDocumentFUID.Error()},
		{"   4ff30392-8ec8-45a4-ba94-5e22c4a686de", true, consts.ErrInvalidDocumentFUID.Error()},
		{"4ff30392-8ec8-45a4-ba94-5e22c4a686de    ", true, consts.ErrInvalidDocumentFUID.Error()},
		{"4ff303928ec8-45a4-ba94-5e22c4a686de", true, consts.ErrInvalidDocumentFUID.Error()},
		{"4ff30392-8ec845a4-ba94-5e22c4a686de", true, consts.ErrInvalidDocumentFUID.Error()},
		{"4ff30392-8ec8-45a4ba94-5e22c4a686de", true, consts.ErrInvalidDocumentFUID.Error()},
		{"4ff30392-8ec8-45a4-ba945e22c4a686de", true, consts.ErrInvalidDocumentFUID.Error()},
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
		{"", "Lisa", true, consts.ErrInvalidDocumentLastName.Error()},
		{"Kim", "123456789123456789012345678901234", true, consts.ErrInvalidDocumentFirstName.Error()},
		{"       ", "Leesa", true, consts.ErrInvalidDocumentLastName.Error()},
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
		{"", true, consts.ErrInvalidDocumentLastName.Error()},
		{"123456789123456789012345678901234", true, consts.ErrInvalidDocumentLastName.Error()},
		{"       ", true, consts.ErrInvalidDocumentLastName.Error()},
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
		{"", true, consts.ErrInvalidDocumentFirstName.Error()},
		{"123456789123456789012345678901234", true, consts.ErrInvalidDocumentFirstName.Error()},
		{"       ", true, consts.ErrInvalidDocumentFirstName.Error()},
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
		{"", true, consts.ErrInvalidDocumentCallTypeName.Error()},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, consts.ErrInvalidDocumentCallTypeName.Error()},
		{"       ", true, consts.ErrInvalidDocumentCallTypeName.Error()},
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
		{"", true, consts.ErrInvalidDocumentGroundType.Error()},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, consts.ErrInvalidDocumentGroundType.Error()},
		{"       ", true, consts.ErrInvalidDocumentGroundType.Error()},
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
		{"Batangas City", "", "12345678912345678901234567890123412345678912345678901234567890123", "Philippines", true, consts.ErrInvalidDocumentProvince.Error()},
		{"Tijuana", "12345678912345678901234567890123412345678912345678901234567890123", "", "Mexico", true, consts.ErrInvalidDocumentState.Error()},
		{"", "Baja California", "", "Mexico", true, consts.ErrInvalidDocumentCity.Error()},
		{"Tijuana", "Baja California", "", "", true, consts.ErrInvalidDocumentCountry.Error()},
		{"San Diego", "CA", "", "12345678912345678901234567890123412345678912345678901234567890123", true, consts.ErrInvalidDocumentCountry.Error()},
		{"       ", "Baja California", "", "Mexico", true, consts.ErrInvalidDocumentCity.Error()},
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
		{"", true, consts.ErrInvalidDocumentCity.Error()},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, consts.ErrInvalidDocumentCity.Error()},
		{"       ", true, consts.ErrInvalidDocumentCity.Error()},
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
		{"12345678912345678901234567890123412345678912345678901234567890123", true, consts.ErrInvalidDocumentState.Error()},
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
		{"12345678912345678901234567890123412345678912345678901234567890123", true, consts.ErrInvalidDocumentProvince.Error()},
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
		{"", true, consts.ErrInvalidDocumentCountry.Error()},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, consts.ErrInvalidDocumentCountry.Error()},
		{"       ", true, consts.ErrInvalidDocumentCountry.Error()},
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

		{"Atlantic ocean    hello ", true, consts.ErrInvalidDocumentOcean.Error()},
		{"Atlantic oceans", true, consts.ErrInvalidDocumentOcean.Error()},
		{"", true, consts.ErrInvalidDocumentOcean.Error()},
		{"      ", true, consts.ErrInvalidDocumentOcean.Error()},
		{"idonotexist", true, consts.ErrInvalidDocumentOcean.Error()},
		{"Indian 1 Ocean", true, consts.ErrInvalidDocumentOcean.Error()},
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
		{"", true, consts.ErrInvalidDocumentSensorType.Error()},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, consts.ErrInvalidDocumentSensorType.Error()},
		{"       ", true, consts.ErrInvalidDocumentSensorType.Error()},
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
		{"", true, consts.ErrInvalidDocumentSensorName.Error()},
		{"12345678912345678901234567890123412345678912345678901234567890123", true, consts.ErrInvalidDocumentSensorName.Error()},
		{"       ", true, consts.ErrInvalidDocumentSensorName.Error()},
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
		{maxSamplingRate + 1, true, consts.ErrInvalidDocumentSamplingRate.Error()},
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
		{minLatitude - 1, true, consts.ErrInvalidDocumentLatitude.Error()},
		{minLatitude, false, ""},
		{0, false, ""},
		{45, false, ""},
		{maxLatitude, false, ""},
		{maxLatitude + 1, true, consts.ErrInvalidDocumentLatitude.Error()},
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
		{minLongitude - 1, true, consts.ErrInvalidDocumentLongitude.Error()},
		{minLongitude, false, ""},
		{0, false, ""},
		{150, false, ""},
		{maxLongitude, false, ""},
		{maxLongitude + 1, true, consts.ErrInvalidDocumentLongitude.Error()},
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
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/images/Rotating_earth_(large).gif",
		}, false, ""},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba9a4-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/images/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.jpg",
		}, true, consts.ErrInvalidDocumentFUID.Error()},
		{nil,
			true, consts.ErrInvalidDocumentImageURLs.Error(),
		},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/images/Rotating_earth_(large).gif",
		}, true, consts.ErrInvalidDocumentImageURL.Error()},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "hwscdevstorage.blob.core.windows.net/images/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.jpg",
		}, true, "invalid Document ImageURL: hwscdevstorage.blob.core.windows.net/images/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.jpg"},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3",
		}, true, "invalid Document image type ImageURL: https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3"},
		{map[string]string{}, false, ""},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/images/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigratingjpg",
		}, true, "invalid Document image type ImageURL: https://hwscdevstorage.blob.core.windows.net/images/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigratingjpg"},
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
		}, true, consts.ErrInvalidDocumentFUID.Error(),
		},
		{nil,
			true, consts.ErrInvalidDocumentAudioURLs.Error(),
		},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/audios/Seger_Conga_CaboMexico_Tag_Acousonde_20140313_112313_8000_3_BreedingMigrating.wav",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3",
		}, true, consts.ErrInvalidDocumentAudioURL.Error()},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3",
		}, true, "invalid Document AudioURL: hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128].mp3"},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png",
		}, true, "invalid Document audio type AudioURL: https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png"},
		{map[string]string{}, false, ""},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128]mp3",
		}, true, "invalid Document audio type AudioURL: https://hwscdevstorage.blob.core.windows.net/audios/Milad Hosseini - Deli Asheghetam [128]mp3"},
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
		}, true, consts.ErrInvalidDocumentFUID.Error(),
		},
		{nil,
			true, consts.ErrInvalidDocumentVideoURLs.Error(),
		},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4",
		}, true, consts.ErrInvalidDocumentVideoURL.Error()},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
		}, true, "invalid Document VideoURL: hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv"},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png",
		}, true, "invalid Document video type VideoURL: https://hwscdevstorage.blob.core.windows.net/images/hulkgif.png"},
		{map[string]string{}, false, ""},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "https://hwscdevstorage.blob.core.windows.net/videos/videoplaybackwmv",
		}, true, "invalid Document video type VideoURL: https://hwscdevstorage.blob.core.windows.net/videos/videoplaybackwmv"},
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
		}, true, consts.ErrInvalidDocumentFUID.Error(),
		},
		{nil,
			true, consts.ErrInvalidDocumentFileURLs.Error(),
		},
		{map[string]string{
			"4ff30392-8ec8-45a4-ba94-5e22c4a686de": "",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d1": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.wmv",
			"4ff30392-8ec8-45a4-ba94-5e22c4a686d2": "https://hwscdevstorage.blob.core.windows.net/videos/videoplayback.mp4",
		}, true, consts.ErrInvalidDocumentFileURL.Error()},
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
		{minTimestamp - 1, true, consts.ErrInvalidDocumentRecordTimestamp.Error()},
		{time.Now().UTC().Unix() + 100, true, consts.ErrInvalidDocumentRecordTimestamp.Error()},
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
			consts.ErrInvalidDocumentCreateTimestamp.Error()},
		{time.Now().UTC().Unix() + 100, 1539831496, true,
			consts.ErrInvalidDocumentCreateTimestamp.Error()},
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
			consts.ErrInvalidUpdateTimestamp.Error()},
		{time.Now().UTC().Unix() + 100, 1539831496, true,
			consts.ErrInvalidUpdateTimestamp.Error()},
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
	// NOTE: force a race condition by commenting out the locks inside isStateAvailable()

	// test for unavailbility
	serviceStateLocker.currentServiceState = unavailable
	assert.Equal(t, unavailable, serviceStateLocker.currentServiceState)

	ok := isStateAvailable()
	assert.Equal(t, false, ok)

	// test for availability
	serviceStateLocker.currentServiceState = available
	assert.Equal(t, available, serviceStateLocker.currentServiceState)

	ok = isStateAvailable()
	assert.Equal(t, true, ok)

	// test race conditions
	const count = 20
	var wg sync.WaitGroup
	start := make(chan struct{}) // signal channel

	wg.Add(count) // #count go routines to wait for

	for i := 0; i < count; i++ {
		go func() {
			<-start // blocks code below, until channel is closed

			defer wg.Done()
			_ = isStateAvailable()
		}()
	}

	close(start) // starts executing blocked goroutines almost at the same time

	// test that read-lock inside isStateAvailable() blocks this write-lock
	serviceStateLocker.lock.Lock()
	serviceStateLocker.currentServiceState = available
	serviceStateLocker.lock.Unlock()

	wg.Wait() // wait until all goroutines finish executing
}

func TestBuildAggregatePipeline(t *testing.T) {
	cases := []struct {
		input     *pbdoc.QueryTransaction
		expOutput bson.A
		isExpErr  bool
		errorStr  string
	}{
		{nil, nil, true, consts.ErrNilQueryTransaction.Error()},
		{
			&pbdoc.QueryTransaction{
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
				StudySites: []*pbdoc.StudySite{
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
			bson.A{
				bson.M{"$match": bson.M{
					"$and": bson.A{
						bson.M{"publisherName.lastName": bson.M{"$in": bson.A{"Seger", "Abadi"}}},
						bson.M{"publisherName.firstName": bson.M{"$in": bson.A{"Kerri", "Shima"}}},
						bson.M{"studySite.city": bson.M{"$in": bson.A{"San Diego", "Batangas City", "Some City"}}},
						bson.M{"studySite.state": bson.M{"$in": bson.A{"California"}}},
						bson.M{"studySite.province": bson.M{"$in": bson.A{"Batangas"}}},
						bson.M{"studySite.country": bson.M{"$in": bson.A{"USA", "Philippines", "Some Country"}}},
						bson.M{"callTypeName": bson.M{"$in": bson.A{mongoDBPatternAll}}},
						bson.M{"groundType": bson.M{"$in": bson.A{"Wookie"}}},
						bson.M{"sensorType": bson.M{"$in": bson.A{"BProbe"}}},
						bson.M{"sensorName": bson.M{"$in": bson.A{"Moto"}}},
						bson.M{"recordTimestamp": bson.M{"$gte": int64(0), "$lte": int64(0)}},
					},
				},
				},
			},
			false,
			"",
		},
	}

	for _, c := range cases {
		output, err := buildAggregatePipeline(c.input)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, c.expOutput, output)
		}
	}
}

func TestBuildArrayFromElements(t *testing.T) {
	cases := []struct {
		input     []string
		expOutput bson.A
	}{
		{nil, bson.A{mongoDBPatternAll}},
		{[]string{}, bson.A{mongoDBPatternAll}},

		{[]string{
			"Seger",
			"Abadi",
		},
			bson.A{"Seger", "Abadi"},
		},
	}

	for _, c := range cases {
		elems := buildArrayFromElements(c.input)
		assert.Equal(t, c.expOutput, elems)
	}
}

func TestExtractPublishersFields(t *testing.T) {
	cases := []struct {
		input      []*pbdoc.Publisher
		lastNames  []string
		firstNames []string
	}{
		{nil, []string{}, []string{}},
		{[]*pbdoc.Publisher{}, []string{}, []string{}},
		{
			[]*pbdoc.Publisher{
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
		assert.True(t, areSlicesEqual(c.lastNames, lastNames))
		assert.True(t, areSlicesEqual(c.firstNames, firstNames))
	}
}

func TestExtractStudySitesFields(t *testing.T) {
	cases := []struct {
		input     []*pbdoc.StudySite
		cities    []string
		states    []string
		provinces []string
		countries []string
	}{
		{nil, []string{}, []string{}, []string{}, []string{}},
		{[]*pbdoc.StudySite{}, []string{}, []string{}, []string{}, []string{}},
		{
			[]*pbdoc.StudySite{
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
		assert.True(t, areSlicesEqual(c.cities, cities))
		assert.True(t, areSlicesEqual(c.states, states))
		assert.True(t, areSlicesEqual(c.provinces, provinces))
		assert.True(t, areSlicesEqual(c.countries, countries))
	}
}

func TestExtractDistinctResults(t *testing.T) {
	cases := []struct {
		input       []interface{}
		queryResult *pbdoc.QueryTransaction
		expOutput   *pbdoc.QueryTransaction
		fieldName   string
		isExpErr    bool
		errorStr    string
	}{
		{nil, nil, nil, "", true, consts.ErrNilQueryResult.Error()},
		{nil, &pbdoc.QueryTransaction{}, nil, "", true, consts.ErrInvalidDistinctResult.Error()},
		{
			[]interface{}{
				bson.D{
					{"lastName", "Seger"},
					{"firstName", "Kerri"},
				},
				bson.D{
					{"lastName", "Abadi"},
					{"firstName", "Shima"},
				},
			},
			&pbdoc.QueryTransaction{},
			&pbdoc.QueryTransaction{
				Publishers: []*pbdoc.Publisher{
					{LastName: "Seger", FirstName: "Kerri"},
					{LastName: "Abadi", FirstName: "Shima"},
				},
			},
			"Publishers",
			false,
			"",
		},
		{
			[]interface{}{
				bson.D{
					{"city", "Batangas City"},
					{"state", ""},
					{"province", "Batangas"},
					{"country", "Philippines"},
				},
				bson.D{
					{"city", "San Diego"},
					{"state", "California"},
					{"province", ""},
					{"country", "USA"},
				},
			},
			&pbdoc.QueryTransaction{},
			&pbdoc.QueryTransaction{
				StudySites: []*pbdoc.StudySite{
					{City: "Batangas City", Province: "Batangas", Country: "Philippines"},
					{City: "San Diego", State: "California", Country: "USA"},
				},
			},
			"StudySites",
			false,
			"",
		},
		{
			[]interface{}{
				"ab",
				"cd",
			},
			&pbdoc.QueryTransaction{},
			&pbdoc.QueryTransaction{
				CallTypeNames: []string{
					"ab",
					"cd",
				},
			},
			"CallTypeNames",
			false,
			"",
		},
		{
			[]interface{}{
				"ef",
				"gh",
			},
			&pbdoc.QueryTransaction{},
			&pbdoc.QueryTransaction{
				GroundTypes: []string{
					"ef",
					"gh",
				},
			},
			"GroundTypes",
			false,
			"",
		},
		{
			[]interface{}{
				"ij",
				"kl",
			},
			&pbdoc.QueryTransaction{},
			&pbdoc.QueryTransaction{
				SensorTypes: []string{
					"ij",
					"kl",
				},
			},
			"SensorTypes",
			false,
			"",
		},
		{
			[]interface{}{
				"mn",
				"op",
			},
			&pbdoc.QueryTransaction{},
			&pbdoc.QueryTransaction{
				SensorNames: []string{
					"mn",
					"op",
				},
			},
			"SensorNames",
			false,
			"",
		},
		{
			[]interface{}{
				"mn",
			},
			&pbdoc.QueryTransaction{},
			&pbdoc.QueryTransaction{
				SensorNames: []string{
					"mn",
				},
			},
			"default",
			true,
			consts.ErrInvalidDistinctFieldName.Error(),
		},
	}

	for _, c := range cases {
		err := extractDistinctResults(c.queryResult, c.fieldName, c.input)
		if !c.isExpErr {
			assert.Equal(t, c.expOutput, c.queryResult)
		} else {
			assert.EqualError(t, err, c.errorStr)
		}
	}
}

func TestExtractDistinctPublishers(t *testing.T) {
	cases := []struct {
		input     []interface{}
		expOutput []*pbdoc.Publisher
		isExpErr  bool
		errorStr  string
	}{
		{nil, nil, true, consts.ErrInvalidDistinctResult.Error()},
		{[]interface{}{}, nil, true, consts.ErrInvalidDistinctResult.Error()},
		{
			[]interface{}{
				bson.D{
					{"lastName", "Seger"},
					{"firstName", "Kerri"},
				},
				bson.D{
					{"lastName", "Abadi"},
					{"firstName", "Shima"},
				},
			},
			[]*pbdoc.Publisher{
				{LastName: "Seger", FirstName: "Kerri"},
				{LastName: "Abadi", FirstName: "Shima"},
			},
			false,
			"",
		},
	}

	for _, c := range cases {
		ret, err := extractDistinctPublishers(c.input)
		if !c.isExpErr {
			assert.Equal(t, c.expOutput, ret)
		} else {
			assert.EqualError(t, err, c.errorStr)
		}
	}
}

func TestExtractDistinctStudySites(t *testing.T) {
	cases := []struct {
		input     []interface{}
		expOutput []*pbdoc.StudySite
		isExpErr  bool
		errorStr  string
	}{
		{nil, nil, true, consts.ErrInvalidDistinctResult.Error()},
		{[]interface{}{}, nil, true, consts.ErrInvalidDistinctResult.Error()},
		{
			[]interface{}{
				bson.D{
					{"city", "Batangas City"},
					{"state", ""},
					{"province", "Batangas"},
					{"country", "Philippines"},
				},
				bson.D{
					{"city", "San Diego"},
					{"state", "California"},
					{"province", ""},
					{"country", "USA"},
				},
			},
			[]*pbdoc.StudySite{
				{City: "Batangas City", Province: "Batangas", Country: "Philippines"},
				{City: "San Diego", State: "California", Country: "USA"},
			},
			false,
			"",
		},
	}

	for _, c := range cases {
		ret, err := extractDistinctStudySites(c.input)
		if !c.isExpErr {
			assert.Equal(t, c.expOutput, ret)
		} else {
			assert.EqualError(t, err, c.errorStr)
		}
	}
}

func TestExtractDistinct(t *testing.T) {
	cases := []struct {
		input     []interface{}
		expOutput []string
		isExpErr  bool
		errorStr  string
	}{
		{nil, nil, true, consts.ErrInvalidDistinctResult.Error()},
		{[]interface{}{}, nil, true, consts.ErrInvalidDistinctResult.Error()},
		{[]interface{}{"a", "b", "c"}, []string{"a", "b", "c"}, false, ""},
	}

	for _, c := range cases {
		ret, err := extractDistinct(c.input)
		if !c.isExpErr {
			assert.Equal(t, c.expOutput, ret)
		} else {
			assert.EqualError(t, err, c.errorStr)
		}
	}
}

func TestAreSlicesEqual(t *testing.T) {
	cases := []struct {
		inputA    []string
		inputB    []string
		expOutput bool
	}{
		{nil, nil, true},
		{[]string{}, []string{}, true},
		{nil, []string{}, false},
		{[]string{"a"}, []string{"a"}, true},
		{[]string{"a"}, []string{"b"}, false},
		{[]string{"a"}, []string{"a", "b"}, false},
	}

	for _, c := range cases {
		assert.Equal(t, c.expOutput, areSlicesEqual(c.inputA, c.inputB))
	}
}

func TestValidateURL(t *testing.T) {
	cases := []struct {
		input    string
		isExpErr bool
		errorStr string
	}{
		{"", true, consts.ErrUnreachableURI.Error()},
		{"https://hwscdevstorage.blob.core.windows.net/imag/Rotating_earth_(large).gif", true, consts.ErrUnreachableURI.Error()},
		{"https://hwscdevstorage.blob.core.windows.net/images/Rotating_earth_(large).gif", false, ""},
	}

	for _, c := range cases {
		err := ValidateURL(c.input)
		if !c.isExpErr {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, c.errorStr)
		}
	}

}
