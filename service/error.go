package service

import (
	"errors"
)

var (
	errServiceUnavailable             = errors.New("service unavailable")
	errNilMongoDBClient               = errors.New("nil MongoDB client")
	errNilRequest                     = errors.New("nil request")
	errNilRequestData                 = errors.New("nil request data")
	errMissingDUID                    = errors.New("missing DUID")
	errNilQueryArgs                   = errors.New("nil query arguments")
	errNilQueryResult                 = errors.New("nil query result")
	errInvalidDistinctResult          = errors.New("invalid distinct result")
	errInvalidDistinctFieldName       = errors.New("invalid distinct field name")
	errAtLeastOneImageAudioURL        = errors.New("requires at least 1 valid Document ImageURL or AudioURL")
	errInvalidDocumentDUID            = errors.New("invalid Document duid")
	errInvalidDocumentUUID            = errors.New("invalid Document uuid")
	errInvalidDocumentFUID            = errors.New("invalid Document fuid")
	errInvalidDocumentLastName        = errors.New("invalid Document LastName")
	errInvalidDocumentFirstName       = errors.New("invalid Document FirstName")
	errInvalidDocumentCallTypeName    = errors.New("invalid Document CallTypeName")
	errInvalidDocumentGroundType      = errors.New("invalid Document GroundType")
	errInvalidDocumentCity            = errors.New("invalid Document City")
	errInvalidDocumentState           = errors.New("invalid Document State")
	errInvalidDocumentProvince        = errors.New("invalid Document Province")
	errInvalidDocumentCountry         = errors.New("invalid Document Country")
	errInvalidDocumentOcean           = errors.New("invalid Document Ocean")
	errInvalidDocumentSensorType      = errors.New("invalid Document SensorType")
	errInvalidDocumentSensorName      = errors.New("invalid Document SensorName")
	errInvalidDocumentSamplingRate    = errors.New("invalid Document SamplingRate")
	errInvalidDocumentLatitude        = errors.New("invalid Document Latitude")
	errInvalidDocumentLongitude       = errors.New("invalid Document Longitude")
	errInvalidDocumentImageURLs       = errors.New("nil Document ImageURLs")
	errInvalidDocumentImageURL        = errors.New("invalid Document ImageURL")
	errInvalidDocumentAudioURLs       = errors.New("nil Document AudioURLs")
	errInvalidDocumentAudioURL        = errors.New("invalid Document AudioURL")
	errInvalidDocumentVideoURLs       = errors.New("nil Document VideoURLs")
	errInvalidDocumentVideoURL        = errors.New("invalid Document VideoURL")
	errInvalidDocumentFileURLs        = errors.New("nil Document FileURLs")
	errInvalidDocumentFileURL         = errors.New("invalid Document FileURL")
	errInvalidDocumentRecordTimestamp = errors.New("invalid Document RecordTimestamp")
	errInvalidDocumentCreateTimestamp = errors.New("invalid Document CreateTimestamp")
	errInvalidUpdateTimestamp         = errors.New("invalid Document UpdateTimestamp")
	errNilQueryTransaction            = errors.New("nil QueryTransaction")
)
