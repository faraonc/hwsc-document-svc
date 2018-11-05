package service

import (
	"errors"
	"fmt"
	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-document-svc/proto"
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
	if err := ValidateStudySite(meta.GetStudySite().GetCity(), meta.GetStudySite().GetState(), //TODO unit test
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
		return errors.New("requires at least 1 valid Document ImageURL or AudioURL")
	}

	return nil
}

// ValidateDUID validates duid.
// Returns an error if duid fails regex validation, and it is not an empty string.
func ValidateDUID(duid string) error {
	if !duidRegex.MatchString(duid) && duid != "" {
		return errors.New("invalid Document duid")
	}
	return nil
}

// ValidateUUID validates uuid.
// Returns an error if uuid fails regex validation.
func ValidateUUID(uuid string) error {
	if !uuidRegex.MatchString(uuid) {
		return errors.New("invalid Document uuid")
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

// ValidateFUID validates fuid.
// Returns an error if fuid fails regex validation.
func ValidateFUID(fuid string) error {
	if !fuidRegex.MatchString(fuid) {
		return errors.New("invalid Document fuid")
	}
	return nil
}

// ValidateLastName validates last name.
// Returns an error if last name is an empty string or exceeds 32 chars.
func ValidateLastName(lastName string) error {
	if strings.TrimSpace(lastName) == "" || len(lastName) > maxLastNameLength {
		return errors.New("invalid Document LastName")
	}
	return nil
}

// ValidateFirstName validates first name.
// Returns an error if first name is an empty string or exceeds 32 chars.
func ValidateFirstName(firstName string) error {
	if strings.TrimSpace(firstName) == "" || len(firstName) > maxFirstNameLength {
		return errors.New("invalid Document FirstName")
	}
	return nil
}

// ValidateCallTypeName validates call type name.
// Returns an error if call type name is an empty string or exceeds 64 chars.
func ValidateCallTypeName(callTypeName string) error {
	if strings.TrimSpace(callTypeName) == "" || len(callTypeName) > maxCallTypeNameLength {
		return errors.New("invalid Document CallTypeName")
	}
	return nil
}

// ValidateGroundType validates ground type.
// Returns an error if ground type is an empty string or exceeds 64 chars.
func ValidateGroundType(groundType string) error {
	if strings.TrimSpace(groundType) == "" || len(groundType) > maxGroundTypeLength {
		return errors.New("invalid Document GroundType")
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
		return errors.New("invalid Document City")
	}
	return nil
}

// ValidateState validates state study site.
// Returns an error if state exceeds 32 chars. //TODO unit test
func ValidateState(state string) error {
	if len(state) > maxStateLength {
		return errors.New("invalid Document State")
	}
	return nil
}

// ValidateProvince validates province study site.
// Returns an error if province exceeds 48 chars. //TODO unit test
func ValidateProvince(province string) error {
	if len(province) > maxProvinceLength {
		return errors.New("invalid Document Province")
	}
	return nil
}

// ValidateCountry validates country study site.
// Returns an error if country is an empty string or exceeds 64 chars.
func ValidateCountry(country string) error {
	if strings.TrimSpace(country) == "" || len(country) > maxCountryLength {
		return errors.New("invalid Document Country")
	}
	return nil
}

// ValidateOcean validates ocean.
// Returns an error if ocean is an empty string or exceeds 64 chars.
func ValidateOcean(ocean string) error {
	if strings.TrimSpace(ocean) == "" {
		return errors.New("invalid Document Ocean")
	}

	w := strings.Fields(ocean)
	switch len(w) {
	case 1:
		{
			if !oceanMap[strings.ToLower(w[0])] {
				return errors.New("invalid Document Ocean")
			}
		}
	case 2:
		{
			if !strings.EqualFold(strings.ToLower(w[1]), "ocean") || !oceanMap[strings.ToLower(w[0])] {
				return errors.New("invalid Document Ocean")
			}
		}
	default:
		{
			return errors.New("invalid Document Ocean")
		}

	}

	return nil
}

// ValidateSensorType validates sensor type.
// Returns an error if sensor type is an empty string or exceeds 64 chars.
func ValidateSensorType(sensorType string) error {
	if strings.TrimSpace(sensorType) == "" || len(sensorType) > maxSensorTypeLength {
		return errors.New("invalid Document SensorType")
	}
	return nil
}

// ValidateSensorName validates sensor name.
// Returns an error if sensor name is an empty string or exceeds 64 chars.
func ValidateSensorName(sensorName string) error {
	if strings.TrimSpace(sensorName) == "" || len(sensorName) > maxSensorNameLength {
		return errors.New("invalid Document SensorName")
	}
	return nil
}

// ValidateSamplingRate validates sampling rate.
// Returns an error if sampling rate exceeds max sampling rate of 4000000000 KHz.
func ValidateSamplingRate(samplingRate uint32) error {
	if samplingRate > maxSamplingRate {
		return errors.New("invalid Document SamplingRate")
	}
	return nil
}

// ValidateLatitude validates latitude.
// Returns an error if latitude is not within [-90,90].
func ValidateLatitude(latitude float32) error {
	if latitude > maxLatitude || latitude < minLatitude {
		return errors.New("invalid Document Latitude")
	}
	return nil
}

// ValidateLongitude validates longitude.
// Returns an error if a longitude is not within [-180,180].
func ValidateLongitude(longitude float32) error {
	if longitude > maxLongitude || longitude < minLongitude {
		return errors.New("invalid Document Longitude")
	}
	return nil
}

// ValidateImageURLs validates image urls.
// Returns an error if a url is an empty string, unsupported format, or unreachable.
func ValidateImageURLs(imageURLs map[string]string) error {
	if imageURLs == nil {
		return errors.New("nil Document ImageURLs")
	}

	for k, v := range imageURLs {
		if !fuidRegex.MatchString(k) {
			return errors.New("invalid Document fuid")
		}

		if strings.TrimSpace(v) == "" {
			return errors.New("invalid Document ImageURL")
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
		return errors.New("nil Document AudioURLs")
	}

	for k, v := range audioURLs {
		if !fuidRegex.MatchString(k) {
			return errors.New("invalid Document fuid")
		}

		if strings.TrimSpace(v) == "" {
			return errors.New("invalid Document AudioURL")
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
		return errors.New("nil Document VideoURLs")
	}

	for k, v := range videoURLs {
		if !fuidRegex.MatchString(k) {
			return errors.New("invalid Document fuid")
		}

		if strings.TrimSpace(v) == "" {
			return errors.New("invalid Document VideoURL")
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
		return errors.New("nil Document FileURLs")
	}

	for k, v := range fileURLs {
		if !fuidRegex.MatchString(k) {
			return errors.New("invalid Document fuid")
		}

		if strings.TrimSpace(v) == "" {
			return errors.New("invalid Document FileURL")
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
		return errors.New("invalid Document RecordTimestamp")
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
		return errors.New("invalid Document CreateTimestamp")
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
		return errors.New("invalid Document UpdateTimestamp")
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
