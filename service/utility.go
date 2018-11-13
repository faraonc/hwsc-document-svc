package service

import (
	"fmt"
	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-document-svc/proto"
	"github.com/mongodb/mongo-go-driver/bson"
	"log"
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
)

// TODO regex to point to the proper storage in Azure
var (
	// https://godoc.org/github.com/segmentio/ksuid
	duidRegex = regexp.MustCompile("^[[:digit:][:alpha:]]{27}$")

	// https://godoc.org/github.com/oklog/ulid
	uuidRegex = regexp.MustCompile("^[[:digit:][:alpha:]]{26}$")

	//https://godoc.org/github.com/google/uuid
	fuidRegex = regexp.MustCompile("^[[:digit:]a-f]{8}-[[:digit:]a-f]{4}-[[:digit:]a-f]{4}-[[:digit:]a-f]{4}-[[:digit:]a-f]{12}$")

	imageRegex = regexp.MustCompile("^.*(jpg|jpeg|png|bmp|tif|gif|tiff).*$")
	audioRegex = regexp.MustCompile("^.*(wav|wma|ogg|m4a|mp3).*$")
	videoRegex = regexp.MustCompile("^.*(flv|wmv|mov|avi|mp4).*$")
	oceanMap   = map[string]bool{
		"pacific":  true,
		"atlantic": true,
		"indian":   true,
		"southern": true,
		"arctic":   true,
	}
)

// ValidateDocument validates the Document.
// Returns an error if field fails validation.
func ValidateDocument(meta *pb.Document) error {
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
		return errAtLeastOneImageAudioURL
	}

	return nil
}

// ValidateDUID validates duid.
// Returns an error if duid fails regex validation, and it is not an empty string.
func ValidateDUID(duid string) error {
	if !duidRegex.MatchString(duid) && duid != "" {
		return errInvalidDocumentDUID
	}
	return nil
}

// ValidateUUID validates uuid.
// Returns an error if uuid fails regex validation.
func ValidateUUID(uuid string) error {
	if !uuidRegex.MatchString(uuid) {
		return errInvalidDocumentUUID
	}
	return nil
}

// ValidateFUID validates fuid.
// Returns an error if fuid fails regex validation.
func ValidateFUID(fuid string) error {
	if !fuidRegex.MatchString(fuid) {
		return errInvalidDocumentFUID
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
		return errInvalidDocumentLastName
	}
	return nil
}

// ValidateFirstName validates first name.
// Returns an error if first name is an empty string or exceeds 32 chars.
func ValidateFirstName(firstName string) error {
	if strings.TrimSpace(firstName) == "" || len(firstName) > maxFirstNameLength {
		return errInvalidDocumentFirstName
	}
	return nil
}

// ValidateCallTypeName validates call type name.
// Returns an error if call type name is an empty string or exceeds 64 chars.
func ValidateCallTypeName(callTypeName string) error {
	if strings.TrimSpace(callTypeName) == "" || len(callTypeName) > maxCallTypeNameLength {
		return errInvalidDocumentCallTypeName
	}
	return nil
}

// ValidateGroundType validates ground type.
// Returns an error if ground type is an empty string or exceeds 64 chars.
func ValidateGroundType(groundType string) error {
	if strings.TrimSpace(groundType) == "" || len(groundType) > maxGroundTypeLength {
		return errInvalidDocumentGroundType
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
		return errInvalidDocumentCity
	}
	return nil
}

// ValidateState validates state study site.
// Returns an error if state exceeds 32 chars.
func ValidateState(state string) error {
	if len(state) > maxStateLength {
		return errInvalidDocumentState
	}
	return nil
}

// ValidateProvince validates province study site.
// Returns an error if province exceeds 48 chars.
func ValidateProvince(province string) error {
	if len(province) > maxProvinceLength {
		return errInvalidDocumentProvince
	}
	return nil
}

// ValidateCountry validates country study site.
// Returns an error if country is an empty string or exceeds 64 chars.
func ValidateCountry(country string) error {
	if strings.TrimSpace(country) == "" || len(country) > maxCountryLength {
		return errInvalidDocumentCountry
	}
	return nil
}

// ValidateOcean validates ocean.
// Returns an error if ocean is an empty string or exceeds 64 chars.
func ValidateOcean(ocean string) error {
	if strings.TrimSpace(ocean) == "" {
		return errInvalidDocumentOcean
	}

	w := strings.Fields(ocean)
	switch len(w) {
	case 1:
		{
			if !oceanMap[strings.ToLower(w[0])] {
				return errInvalidDocumentOcean
			}
		}
	case 2:
		{
			if !strings.EqualFold(strings.ToLower(w[1]), "ocean") || !oceanMap[strings.ToLower(w[0])] {
				return errInvalidDocumentOcean
			}
		}
	default:
		{
			return errInvalidDocumentOcean
		}

	}

	return nil
}

// ValidateSensorType validates sensor type.
// Returns an error if sensor type is an empty string or exceeds 64 chars.
func ValidateSensorType(sensorType string) error {
	if strings.TrimSpace(sensorType) == "" || len(sensorType) > maxSensorTypeLength {
		return errInvalidDocumentSensorType
	}
	return nil
}

// ValidateSensorName validates sensor name.
// Returns an error if sensor name is an empty string or exceeds 64 chars.
func ValidateSensorName(sensorName string) error {
	if strings.TrimSpace(sensorName) == "" || len(sensorName) > maxSensorNameLength {
		return errInvalidDocumentSensorName
	}
	return nil
}

// ValidateSamplingRate validates sampling rate.
// Returns an error if sampling rate exceeds max sampling rate of 4000000000 KHz.
func ValidateSamplingRate(samplingRate uint32) error {
	if samplingRate > maxSamplingRate {
		return errInvalidDocumentSamplingRate
	}
	return nil
}

// ValidateLatitude validates latitude.
// Returns an error if latitude is not within [-90,90].
func ValidateLatitude(latitude float32) error {
	if latitude > maxLatitude || latitude < minLatitude {
		return errInvalidDocumentLatitude
	}
	return nil
}

// ValidateLongitude validates longitude.
// Returns an error if a longitude is not within [-180,180].
func ValidateLongitude(longitude float32) error {
	if longitude > maxLongitude || longitude < minLongitude {
		return errInvalidDocumentLongitude
	}
	return nil
}

// ValidateImageURLs validates image urls.
// Returns an error if a url is an empty string, unsupported format, or unreachable.
func ValidateImageURLs(imageURLs map[string]string) error {
	if imageURLs == nil {
		return errInvalidDocumentImageURLs
	}

	for k, v := range imageURLs {
		if !fuidRegex.MatchString(k) {
			return errInvalidDocumentFUID
		}

		if strings.TrimSpace(v) == "" {
			return errInvalidDocumentImageURL
		}
		if !imageRegex.MatchString(strings.ToLower(v)) {
			return fmt.Errorf("invalid Document image type ImageURL: %s", v)
		}
		if _, err := url.ParseRequestURI(v); err != nil {
			return fmt.Errorf("invalid Document ImageURL: %s", v)
		}
	}
	return nil
}

// ValidateAudioURLs validates audio urls.
// Returns an error if a url is an empty string, unsupported format, or unreachable.
func ValidateAudioURLs(audioURLs map[string]string) error {
	if audioURLs == nil {
		return errInvalidDocumentAudioURLs
	}

	for k, v := range audioURLs {
		if !fuidRegex.MatchString(k) {
			return errInvalidDocumentFUID
		}

		if strings.TrimSpace(v) == "" {
			return errInvalidDocumentAudioURL
		}
		if !audioRegex.MatchString(strings.ToLower(v)) {
			return fmt.Errorf("invalid Document audio type AudioURL: %s", v)
		}
		if _, err := url.ParseRequestURI(v); err != nil {
			return fmt.Errorf("invalid Document AudioURL: %s", v)
		}
	}
	return nil
}

// ValidateVideoURLs validates video urls.
// Returns an error if a url is an empty string, unsupported format, or unreachable.
func ValidateVideoURLs(videoURLs map[string]string) error {
	if videoURLs == nil {
		return errInvalidDocumentVideoURLs
	}

	for k, v := range videoURLs {
		if !fuidRegex.MatchString(k) {
			return errInvalidDocumentFUID
		}

		if strings.TrimSpace(v) == "" {
			return errInvalidDocumentVideoURL
		}
		if !videoRegex.MatchString(strings.ToLower(v)) {
			return fmt.Errorf("invalid Document video type VideoURL: %s", v)
		}
		if _, err := url.ParseRequestURI(v); err != nil {
			return fmt.Errorf("invalid Document VideoURL: %s", v)
		}
	}
	return nil
}

// ValidateFileURLs validates video urls.
// Returns an error if a url is an empty string, or unreachable.
func ValidateFileURLs(fileURLs map[string]string) error {
	if fileURLs == nil {
		return errInvalidDocumentFileURLs
	}

	for k, v := range fileURLs {
		if !fuidRegex.MatchString(k) {
			return errInvalidDocumentFUID
		}

		if strings.TrimSpace(v) == "" {
			return errInvalidDocumentFileURL
		}
		if _, err := url.ParseRequestURI(v); err != nil {
			return fmt.Errorf("invalid Document FileURL: %s", v)
		}
	}
	return nil
}

// ValidateRecordTimestamp validates record timestamp.
// Returns an error if timestamp is set before Jan 1, 1990 or now.
func ValidateRecordTimestamp(timestamp int64) error {
	if timestamp < minTimestamp || timestamp > time.Now().UTC().Unix() {
		return errInvalidDocumentRecordTimestamp
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
		return errInvalidDocumentCreateTimestamp
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
		return errInvalidUpdateTimestamp
	}

	return nil
}

func isStateAvailable() bool {
	// Lock the state for reading
	serviceStateLocker.lock.RLock()
	// Unlock the state before function exits
	defer serviceStateLocker.lock.RUnlock()

	log.Printf("[INFO] Service State: %s\n", serviceStateLocker.currentServiceState)
	if serviceStateLocker.currentServiceState != available {
		return false
	}

	return true
}

//TODO unit test
//TODO error checking
func buildAggregatePipeline(queryParams *pb.QueryTransaction) (*bson.Array, error) {
	if queryParams == nil {
		return nil, errNilQueryTransaction
	}

	if queryParams.GetPublishers() == nil &&
		queryParams.GetStudySites() == nil &&
		queryParams.GetCallTypes() == nil &&
		queryParams.GetGroundTypes() == nil &&
		queryParams.GetSensorTypes() == nil &&
		queryParams.GetSensorNames() == nil {

		return nil, errNilQueryTransactionFields
	}

	lastNames, firstNames := extractPublishersFields(queryParams.GetPublishers())
	cities, states, provinces, countries := extractStudySitesFields(queryParams.GetStudySites())

	pipeline := bson.NewArray(
		bson.VC.DocumentFromElements(
			bson.EC.SubDocumentFromElements(
				"$match",
				bson.EC.SubDocumentFromElements("publisherName.lastName",
					buildArrayFromElements(lastNames)),
				bson.EC.SubDocumentFromElements("publisherName.firstName",
					buildArrayFromElements(firstNames)),

				bson.EC.SubDocumentFromElements("studySite.city",
					buildArrayFromElements(cities)),
				bson.EC.SubDocumentFromElements("studySite.state",
					buildArrayFromElements(states)),
				bson.EC.SubDocumentFromElements("studySite.province",
					buildArrayFromElements(provinces)),
				bson.EC.SubDocumentFromElements("studySite.country",
					buildArrayFromElements(countries)),

				bson.EC.SubDocumentFromElements("callTypeName",
					buildArrayFromElements(queryParams.GetCallTypes())),

				bson.EC.SubDocumentFromElements("groundType",
					buildArrayFromElements(queryParams.GetGroundTypes())),

				bson.EC.SubDocumentFromElements("sensorType",
					buildArrayFromElements(queryParams.GetSensorTypes())),

				bson.EC.SubDocumentFromElements("sensorName",
					buildArrayFromElements(queryParams.GetSensorNames())),
			),
		),
	)

	return pipeline, nil
}

func buildArrayFromElements(elems []string) *bson.Element {
	elemVals := make([]*bson.Value, len(elems))
	for i := 0; i < len(elems); i++ {
		elemVals[i] = bson.VC.String(elems[i])
	}

	if len(elemVals) == 0{
		return bson.EC.ArrayFromElements("$in", bson.VC.Regex(".*", ""))
	}

	return bson.EC.ArrayFromElements("$in", elemVals...)
}

func extractPublishersFields(publishers []*pb.Publisher) (lastNames, firstNames []string) {
	if publishers == nil {
		return []string{}, []string{}
	}

	lastNames = make([]string, len(publishers))
	firstNames = make([]string, len(publishers))

	for i := 0; i < len(publishers); i++ {
		lastNames[i] = publishers[i].GetLastName()
		firstNames[i] = publishers[i].GetFirstName()
	}

	return
}

func extractStudySitesFields(studySites []*pb.StudySite) (cities, states, provinces, countries []string) {
	if studySites == nil {
		return []string{}, []string{}, []string{}, []string{}
	}

	cities = make([]string, len(studySites))
	states = make([]string, 0)
	provinces = make([]string, 0)
	countries = make([]string, len(studySites))

	for i := 0; i < len(studySites); i++ {
		cities[i] = studySites[i].GetCity()

		tempState := strings.TrimSpace(studySites[i].GetState())
		if tempState != "" {
			states = append(states, tempState)
		}

		tempProvince := strings.TrimSpace(studySites[i].GetProvince())
		if tempState != "" {
			provinces = append(provinces, tempProvince)
		}

		countries[i] = studySites[i].GetCountry()
	}

	return
}
