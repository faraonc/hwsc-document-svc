package service

import (
	"fmt"
	"github.com/google/uuid"
	pbdoc "github.com/hwsc-org/hwsc-api-blocks/protobuf/lib"
	"github.com/hwsc-org/hwsc-document-svc/consts"
	log "github.com/hwsc-org/hwsc-lib/logger"
	"github.com/segmentio/ksuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	maxLastNameLength     = 32
	maxFirstNameLength    = 32
	maxCallTypeNameLength = 64
	maxGroundTypeLength   = 64
	maxCityLength         = 64
	maxStateLength        = 32
	maxProvinceLength     = 48
	maxCountryLength      = 64
	maxSensorTypeLength   = 64
	maxSensorNameLength   = 64
	// KiloHertz
	maxSamplingRate = 4000000000
	maxLatitude     = 90
	minLatitude     = -90
	maxLongitude    = 180
	minLongitude    = -180
	// Minimum timestamp in seconds (Jan 1, 1990)
	minTimestamp = 631152000
	// maxTimestamp = 32503680524 in MongoDB
)

// TODO regex to point to the proper storage in Azure
var (
	// https://godoc.org/github.com/segmentio/ksuid
	duidRegex = regexp.MustCompile("^[[:digit:][:alpha:]]{27}$")

	// https://godoc.org/github.com/oklog/ulid
	uuidRegex = regexp.MustCompile("^[[:digit:][:alpha:]]{26}$")

	//https://godoc.org/github.com/google/uuid
	fuidRegex = regexp.MustCompile("^[[:digit:]a-f]{8}-[[:digit:]a-f]{4}-[[:digit:]a-f]{4}-[[:digit:]a-f]{4}-[[:digit:]a-f]{12}$")

	imageRegex = regexp.MustCompile("^.*\\.(jpg|jpeg|png|bmp|tif|gif|tiff)$")
	audioRegex = regexp.MustCompile("^.*\\.(wav|wma|ogg|m4a|mp3)$")
	videoRegex = regexp.MustCompile("^.*\\.(flv|wmv|mov|avi|mp4)$")
	oceanMap   = map[string]bool{
		"pacific":  true,
		"atlantic": true,
		"indian":   true,
		"southern": true,
		"arctic":   true,
	}
	distinctSearchFieldNames = map[int]string{
		0: "publisherName",
		1: "studySite",
		2: "callTypeName",
		3: "groundType",
		4: "sensorType",
		5: "sensorName",
	}
	distinctResultFieldIndices = map[string]int{
		"Publishers":    0,
		"StudySites":    1,
		"CallTypeNames": 2,
		"GroundTypes":   3,
		"SensorTypes":   4,
		"SensorNames":   5,
	}
	mongoDBPatternAll = primitive.Regex{Pattern: ".*", Options: ""}
)

// NewDUID generates a new document unique ID.
func (d *duidLocker) NewDUID() string {
	d.lock.Lock()
	defer d.lock.Unlock()
	duid := ksuid.New().String()
	return duid
}

// NewFUID generates a new file metadata unique ID.
func (d *fuidLocker) NewFUID() string {
	d.lock.Lock()
	defer d.lock.Unlock()
	newUUID := uuid.New().String()
	return newUUID
}

// ValidateDocument validates the Document.
// Returns an error if field fails validation.
func ValidateDocument(meta *pbdoc.Document) error {
	if err := ValidateDUID(meta.GetDuid()); err != nil {
		return err
	}
	if err := ValidateUUID(meta.GetUuid()); err != nil {
		return err
	}
	if err := ValidatePublisher(meta.GetPublisherName().GetLastName(),
		meta.GetPublisherName().GetFirstName()); err != nil {
		return err
	}
	if err := ValidateCallTypeName(meta.GetCallTypeName()); err != nil {
		return err
	}
	if err := ValidateGroundType(meta.GetGroundType()); err != nil {
		return err
	}
	if err := ValidateStudySite(meta.GetStudySite().GetCity(), meta.GetStudySite().GetState(),
		meta.GetStudySite().GetProvince(), meta.GetStudySite().GetCountry()); err != nil {
		return err
	}
	if err := ValidateOcean(meta.GetOcean()); err != nil {
		return err
	}
	if err := ValidateSensorType(meta.GetSensorType()); err != nil {
		return err
	}
	if err := ValidateSensorName(meta.GetSensorName()); err != nil {
		return err
	}
	if err := ValidateSamplingRate(meta.GetSamplingRate()); err != nil {
		return err
	}
	if err := ValidateLatitude(meta.GetLatitude()); err != nil {
		return err
	}
	if err := ValidateLongitude(meta.GetLongitude()); err != nil {
		return err
	}
	if err := ValidateImageURLs(meta.GetImageUrlsMap()); err != nil {
		return err
	}
	if err := ValidateAudioURLs(meta.GetAudioUrlsMap()); err != nil {
		return err
	}
	if err := ValidateVideoURLs(meta.GetVideoUrlsMap()); err != nil {
		return err
	}
	if err := ValidateFileURLs(meta.GetFileUrlsMap()); err != nil {
		return err
	}
	if err := ValidateRecordTimestamp(meta.GetRecordTimestamp()); err != nil {
		return err
	}
	if err := ValidateCreateTimestamp(meta.GetCreateTimestamp(), meta.GetRecordTimestamp()); err != nil {
		return err
	}
	if err := ValidateUpdateTimestamp(meta.GetUpdateTimestamp(), meta.GetCreateTimestamp()); err != nil {
		return err
	}
	if len(meta.GetImageUrlsMap()) == 0 && len(meta.GetAudioUrlsMap()) == 0 {
		return consts.ErrAtLeastOneImageAudioURL
	}

	return nil
}

// ValidateDUID validates duid.
// Returns an error if duid fails regex validation, and it is not an empty string.
func ValidateDUID(duid string) error {
	if !duidRegex.MatchString(duid) && duid != "" {
		return consts.ErrInvalidDocumentDUID
	}
	return nil
}

// ValidateUUID validates uuid.
// Returns an error if uuid fails regex validation.
func ValidateUUID(uuid string) error {
	if !uuidRegex.MatchString(uuid) {
		return consts.ErrInvalidDocumentUUID
	}
	return nil
}

// ValidateFUID validates fuid.
// Returns an error if fuid fails regex validation.
func ValidateFUID(fuid string) error {
	if !fuidRegex.MatchString(fuid) {
		return consts.ErrInvalidDocumentFUID
	}
	return nil
}

// ValidatePublisher validates the publisher name.
// Returns an error if last name or first name fails validation
func ValidatePublisher(lastName string, firstName string) error {
	if err := ValidateLastName(lastName); err != nil {
		return err
	}
	if err := ValidateFirstName(firstName); err != nil {
		return err
	}

	return nil
}

// ValidateLastName validates last name.
// Returns an error if last name is an empty string or exceeds 32 chars.
func ValidateLastName(lastName string) error {
	if strings.TrimSpace(lastName) == "" || len(lastName) > maxLastNameLength {
		return consts.ErrInvalidDocumentLastName
	}
	return nil
}

// ValidateFirstName validates first name.
// Returns an error if first name is an empty string or exceeds 32 chars.
func ValidateFirstName(firstName string) error {
	if strings.TrimSpace(firstName) == "" || len(firstName) > maxFirstNameLength {
		return consts.ErrInvalidDocumentFirstName
	}
	return nil
}

// ValidateCallTypeName validates call type name.
// Returns an error if call type name is an empty string or exceeds 64 chars.
func ValidateCallTypeName(callTypeName string) error {
	if strings.TrimSpace(callTypeName) == "" || len(callTypeName) > maxCallTypeNameLength {
		return consts.ErrInvalidDocumentCallTypeName
	}
	return nil
}

// ValidateGroundType validates ground type.
// Returns an error if ground type is an empty string or exceeds 64 chars.
func ValidateGroundType(groundType string) error {
	if strings.TrimSpace(groundType) == "" || len(groundType) > maxGroundTypeLength {
		return consts.ErrInvalidDocumentGroundType
	}
	return nil
}

// ValidateStudySite validates study site.
// Returns an error if city, state, province, or country fails validation
func ValidateStudySite(city string, state string, province string, country string) error {
	if err := ValidateCity(city); err != nil {
		return err
	}
	if err := ValidateState(state); err != nil {
		return err
	}
	if err := ValidateProvince(province); err != nil {
		return err
	}
	if err := ValidateCountry(country); err != nil {
		return err
	}

	return nil
}

// ValidateCity validates city study site.
// Returns an error if city is an empty string or exceeds 64 chars.
func ValidateCity(city string) error {
	if strings.TrimSpace(city) == "" || len(city) > maxCityLength {
		return consts.ErrInvalidDocumentCity
	}
	return nil
}

// ValidateState validates state study site.
// Returns an error if state exceeds 32 chars.
func ValidateState(state string) error {
	if len(state) > maxStateLength {
		return consts.ErrInvalidDocumentState
	}
	return nil
}

// ValidateProvince validates province study site.
// Returns an error if province exceeds 48 chars.
func ValidateProvince(province string) error {
	if len(province) > maxProvinceLength {
		return consts.ErrInvalidDocumentProvince
	}
	return nil
}

// ValidateCountry validates country study site.
// Returns an error if country is an empty string or exceeds 64 chars.
func ValidateCountry(country string) error {
	if strings.TrimSpace(country) == "" || len(country) > maxCountryLength {
		return consts.ErrInvalidDocumentCountry
	}
	return nil
}

// ValidateOcean validates ocean.
// Returns an error if ocean is an empty string or exceeds 64 chars.
func ValidateOcean(ocean string) error {
	if strings.TrimSpace(ocean) == "" {
		return consts.ErrInvalidDocumentOcean
	}

	w := strings.Fields(ocean)
	switch len(w) {
	case 1:
		{
			if !oceanMap[strings.ToLower(w[0])] {
				return consts.ErrInvalidDocumentOcean
			}
		}
	case 2:
		{
			if !strings.EqualFold(strings.ToLower(w[1]), "ocean") || !oceanMap[strings.ToLower(w[0])] {
				return consts.ErrInvalidDocumentOcean
			}
		}
	default:
		{
			return consts.ErrInvalidDocumentOcean
		}

	}

	return nil
}

// ValidateSensorType validates sensor type.
// Returns an error if sensor type is an empty string or exceeds 64 chars.
func ValidateSensorType(sensorType string) error {
	if strings.TrimSpace(sensorType) == "" || len(sensorType) > maxSensorTypeLength {
		return consts.ErrInvalidDocumentSensorType
	}
	return nil
}

// ValidateSensorName validates sensor name.
// Returns an error if sensor name is an empty string or exceeds 64 chars.
func ValidateSensorName(sensorName string) error {
	if strings.TrimSpace(sensorName) == "" || len(sensorName) > maxSensorNameLength {
		return consts.ErrInvalidDocumentSensorName
	}
	return nil
}

// ValidateSamplingRate validates sampling rate.
// Returns an error if sampling rate exceeds max sampling rate of 4000000000 KHz.
func ValidateSamplingRate(samplingRate uint32) error {
	if samplingRate > maxSamplingRate {
		return consts.ErrInvalidDocumentSamplingRate
	}
	return nil
}

// ValidateLatitude validates latitude.
// Returns an error if latitude is not within [-90,90].
func ValidateLatitude(latitude float32) error {
	if latitude > maxLatitude || latitude < minLatitude {
		return consts.ErrInvalidDocumentLatitude
	}
	return nil
}

// ValidateLongitude validates longitude.
// Returns an error if a longitude is not within [-180,180].
func ValidateLongitude(longitude float32) error {
	if longitude > maxLongitude || longitude < minLongitude {
		return consts.ErrInvalidDocumentLongitude
	}
	return nil
}

// ValidateImageURLs validates image urls.
// Returns an error if a url is an empty string, unsupported format, or unreachable.
func ValidateImageURLs(imageURLs map[string]string) error {
	if imageURLs == nil {
		return consts.ErrInvalidDocumentImageURLs
	}

	for k, v := range imageURLs {
		if !fuidRegex.MatchString(k) {
			return consts.ErrInvalidDocumentFUID
		}

		if strings.TrimSpace(v) == "" {
			return consts.ErrInvalidDocumentImageURL
		}
		if !imageRegex.MatchString(strings.ToLower(v)) {
			return fmt.Errorf("invalid Document image type ImageURL: %s", v)
		}
		if err := ValidateURL(v); err != nil {
			return fmt.Errorf("invalid Document ImageURL: %s", v)
		}
	}
	return nil
}

// ValidateAudioURLs validates audio urls.
// Returns an error if a url is an empty string, unsupported format, or unreachable.
func ValidateAudioURLs(audioURLs map[string]string) error {
	if audioURLs == nil {
		return consts.ErrInvalidDocumentAudioURLs
	}

	for k, v := range audioURLs {
		if !fuidRegex.MatchString(k) {
			return consts.ErrInvalidDocumentFUID
		}

		if strings.TrimSpace(v) == "" {
			return consts.ErrInvalidDocumentAudioURL
		}
		if !audioRegex.MatchString(strings.ToLower(v)) {
			return fmt.Errorf("invalid Document audio type AudioURL: %s", v)
		}
		if err := ValidateURL(v); err != nil {
			return fmt.Errorf("invalid Document AudioURL: %s", v)
		}
	}
	return nil
}

// ValidateVideoURLs validates video urls.
// Returns an error if a url is an empty string, unsupported format, or unreachable.
func ValidateVideoURLs(videoURLs map[string]string) error {
	if videoURLs == nil {
		return consts.ErrInvalidDocumentVideoURLs
	}

	for k, v := range videoURLs {
		if !fuidRegex.MatchString(k) {
			return consts.ErrInvalidDocumentFUID
		}

		if strings.TrimSpace(v) == "" {
			return consts.ErrInvalidDocumentVideoURL
		}
		if !videoRegex.MatchString(strings.ToLower(v)) {
			return fmt.Errorf("invalid Document video type VideoURL: %s", v)
		}
		if err := ValidateURL(v); err != nil {
			return fmt.Errorf("invalid Document VideoURL: %s", v)
		}
	}
	return nil
}

// ValidateFileURLs validates video urls.
// Returns an error if a url is an empty string, or unreachable.
func ValidateFileURLs(fileURLs map[string]string) error {
	if fileURLs == nil {
		return consts.ErrInvalidDocumentFileURLs
	}

	for k, v := range fileURLs {
		if !fuidRegex.MatchString(k) {
			return consts.ErrInvalidDocumentFUID
		}

		if strings.TrimSpace(v) == "" {
			return consts.ErrInvalidDocumentFileURL
		}
		if err := ValidateURL(v); err != nil {
			return fmt.Errorf("invalid Document FileURL: %s", v)
		}
	}
	return nil
}

// ValidateRecordTimestamp validates record timestamp.
// Returns an error if timestamp is set before Jan 1, 1990 or now.
func ValidateRecordTimestamp(timestamp int64) error {
	if timestamp < minTimestamp || timestamp > time.Now().UTC().Unix() {
		return consts.ErrInvalidDocumentRecordTimestamp
	}

	return nil
}

// ValidateCreateTimestamp validates create timestamp.
// Returns an error if create timestamp is set before record timestamp, or create timestamp is set after now.
func ValidateCreateTimestamp(createTimestamp int64, recordTimeStamp int64) error {
	if createTimestamp == 0 {
		return nil
	}

	if createTimestamp < recordTimeStamp || createTimestamp > time.Now().UTC().Unix() {
		return consts.ErrInvalidDocumentCreateTimestamp
	}

	return nil
}

// ValidateUpdateTimestamp validates create timestamp.
// Returns an error if create timestamp is set after update timestamp, or update timestamp is set after now.
func ValidateUpdateTimestamp(updateTimestamp int64, createTimestamp int64) error {
	if updateTimestamp == 0 {
		return nil
	}

	if createTimestamp > updateTimestamp || updateTimestamp > time.Now().UTC().Unix() {
		return consts.ErrInvalidUpdateTimestamp
	}

	return nil
}

func isStateAvailable() bool {
	// Lock the state for reading
	serviceStateLocker.lock.RLock()
	// Unlock the state before function exits
	defer serviceStateLocker.lock.RUnlock()

	log.Info(consts.ServiceStateTag, serviceStateLocker.currentServiceState.String())
	if serviceStateLocker.currentServiceState != available {
		return false
	}

	return true
}

func buildAggregatePipeline(queryParams *pbdoc.QueryTransaction) (bson.A, error) {
	if queryParams == nil {
		return nil, consts.ErrNilQueryTransaction
	}

	lastNames, firstNames := extractPublishersFields(queryParams.GetPublishers())
	cities, states, provinces, countries := extractStudySitesFields(queryParams.GetStudySites())
	pipeline := bson.A{
		bson.M{"$match": bson.M{
			"$and": bson.A{
				bson.M{"publisherName.lastName": bson.M{"$in": buildArrayFromElements(lastNames)}},
				bson.M{"publisherName.firstName": bson.M{"$in": buildArrayFromElements(firstNames)}},

				bson.M{"studySite.city": bson.M{"$in": buildArrayFromElements(cities)}},
				bson.M{"studySite.state": bson.M{"$in": buildArrayFromElements(states)}},
				bson.M{"studySite.province": bson.M{"$in": buildArrayFromElements(provinces)}},
				bson.M{"studySite.country": bson.M{"$in": buildArrayFromElements(countries)}},

				bson.M{"callTypeName": bson.M{"$in": buildArrayFromElements(queryParams.GetCallTypeNames())}},
				bson.M{"groundType": bson.M{"$in": buildArrayFromElements(queryParams.GetGroundTypes())}},
				bson.M{"sensorType": bson.M{"$in": buildArrayFromElements(queryParams.GetSensorTypes())}},
				bson.M{"sensorName": bson.M{"$in": buildArrayFromElements(queryParams.GetSensorNames())}},
				bson.M{"recordTimestamp": bson.M{"$gte": queryParams.GetMinRecordTimestamp(),
					"$lte": queryParams.GetMaxRecordTimestamp()}},
			},
		},
		},
	}

	return pipeline, nil
}

func buildArrayFromElements(elems []string) bson.A {
	if elems == nil || len(elems) == 0 {
		return bson.A{mongoDBPatternAll}
	}

	result := bson.A{}
	for _, e := range elems {
		result = append(result, e)
	}

	return result
}

func extractPublishersFields(publishers []*pbdoc.Publisher) ([]string, []string) {
	if publishers == nil || len(publishers) == 0 {
		return []string{}, []string{}
	}

	lastNames := make([]string, 0)
	firstNames := make([]string, 0)

	for i := 0; i < len(publishers); i++ {
		tempLastName := strings.TrimSpace(publishers[i].GetLastName())
		if tempLastName != "" {
			lastNames = append(lastNames, tempLastName)
		}

		tempFirstName := strings.TrimSpace(publishers[i].GetFirstName())
		if tempFirstName != "" {
			firstNames = append(firstNames, tempFirstName)
		}
	}

	return lastNames, firstNames
}

func extractStudySitesFields(studySites []*pbdoc.StudySite) ([]string, []string, []string, []string) {
	if studySites == nil || len(studySites) == 0 {
		return []string{}, []string{}, []string{}, []string{}
	}

	cities := make([]string, 0)
	states := make([]string, 0)
	provinces := make([]string, 0)
	countries := make([]string, 0)

	for i := 0; i < len(studySites); i++ {
		tempCity := strings.TrimSpace(studySites[i].GetCity())
		if tempCity != "" {
			cities = append(cities, tempCity)
		}

		tempState := strings.TrimSpace(studySites[i].GetState())
		if tempState != "" {
			states = append(states, tempState)
		}

		tempProvince := strings.TrimSpace(studySites[i].GetProvince())
		if tempProvince != "" {
			provinces = append(provinces, tempProvince)
		}

		tempCountry := strings.TrimSpace(studySites[i].GetCountry())
		if tempCountry != "" {
			countries = append(countries, tempCountry)
		}
	}

	return cities, states, provinces, countries
}

func extractDistinctResults(queryResult *pbdoc.QueryTransaction, fieldName string, distinctResult []interface{}) error {
	if queryResult == nil {
		return consts.ErrNilQueryResult
	}

	if distinctResult == nil || len(distinctResult) == 0 {
		return consts.ErrInvalidDistinctResult
	}

	var err error
	switch fieldName {
	case "Publishers":
		queryResult.Publishers, err = extractDistinctPublishers(distinctResult)
	case "StudySites":
		queryResult.StudySites, err = extractDistinctStudySites(distinctResult)
	case "CallTypeNames":
		queryResult.CallTypeNames, err = extractDistinct(distinctResult)
	case "GroundTypes":
		queryResult.GroundTypes, err = extractDistinct(distinctResult)
	case "SensorTypes":
		queryResult.SensorTypes, err = extractDistinct(distinctResult)
	case "SensorNames":
		queryResult.SensorNames, err = extractDistinct(distinctResult)
	default:
		err = consts.ErrInvalidDistinctFieldName
	}

	return err
}

func extractDistinctPublishers(distinctResult []interface{}) ([]*pbdoc.Publisher, error) {
	if distinctResult == nil || len(distinctResult) == 0 {
		return nil, consts.ErrInvalidDistinctResult
	}
	publishers := make([]*pbdoc.Publisher, 0)
	for _, v := range distinctResult {
		publishers = append(publishers, &pbdoc.Publisher{
			LastName:  v.(bson.D).Map()["lastName"].(string),
			FirstName: v.(bson.D).Map()["firstName"].(string),
		})
	}
	return publishers, nil
}

func extractDistinctStudySites(distinctResult []interface{}) ([]*pbdoc.StudySite, error) {
	if distinctResult == nil || len(distinctResult) == 0 {
		return nil, consts.ErrInvalidDistinctResult
	}
	studySites := make([]*pbdoc.StudySite, 0)
	for _, v := range distinctResult {
		studySites = append(studySites, &pbdoc.StudySite{
			City:     v.(bson.D).Map()["city"].(string),
			State:    v.(bson.D).Map()["state"].(string),
			Province: v.(bson.D).Map()["province"].(string),
			Country:  v.(bson.D).Map()["country"].(string),
		})
	}
	return studySites, nil
}

func extractDistinct(distinctResult []interface{}) ([]string, error) {
	if distinctResult == nil || len(distinctResult) == 0 {
		return nil, consts.ErrInvalidDistinctResult
	}
	distinct := make([]string, len(distinctResult))
	for i := 0; i < len(distinctResult); i++ {
		distinct[i] = distinctResult[i].(string)
	}
	return distinct, nil
}

func areSlicesEqual(a, b []string) bool {

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

// ValidateURL confirms if the address is reachable and valid.
// Return an error if the address is unreachable and invalid.
func ValidateURL(addr string) error {
	if _, err := url.ParseRequestURI(addr); err != nil {
		return consts.ErrUnreachableURI
	}
	resp, err := http.Get(addr)
	if err != nil || resp.StatusCode >= 400 {
		return consts.ErrUnreachableURI
	}
	_ = resp.Body.Close()
	return nil
}
