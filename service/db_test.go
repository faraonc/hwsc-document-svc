package service

import (
	"github.com/hwsc-org/hwsc-document-svc/conf"
	"github.com/hwsc-org/hwsc-document-svc/consts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDialMongoDB(t *testing.T) {
	cases := []struct {
		desc     string
		uri      string
		isExpErr bool
		errorStr string
	}{
		{"test for valid reader", conf.DocumentDB.Reader, false, ""},
		{"test for valid writer", conf.DocumentDB.Writer, false, ""},
		{"test for empty empty uri", "", true, consts.ErrEmptyMongoDBURI.Error()},
	}

	for _, c := range cases {
		client, err := dialMongoDB(&c.uri)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr, c.desc)
		} else {
			assert.Nil(t, err, c.desc)
			assert.NotNil(t, client, c.desc)
		}
	}

}

func TestDisconnectMongoDBClient(t *testing.T) {
	cases := []struct {
		desc     string
		uri      string
		isExpErr bool
		errorStr string
	}{
		{"test for valid reader", conf.DocumentDB.Reader, false, ""},
		{"test for valid writer", conf.DocumentDB.Writer, false, ""},
		{"test for empty nil", "", true, consts.ErrNilMongoDBClient.Error()},
	}

	for _, c := range cases {
		client, _ := dialMongoDB(&c.uri)
		err := disconnectMongoDBClient(client)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr, c.desc)
		} else {
			assert.Nil(t, err, c.desc)
		}
	}
}

func TestRefreshMongoDBConnection(t *testing.T) {
	mongoDBReader = nil
	assert.Nil(t, mongoDBReader)
	err := refreshMongoDBConnection(mongoDBReader, &conf.DocumentDB.Reader)
	assert.Nil(t, err)
	assert.NotNil(t, mongoDBReader)
	mongoDBWriter = nil
	assert.Nil(t, mongoDBWriter)
	err = refreshMongoDBConnection(mongoDBWriter, &conf.DocumentDB.Writer)
	assert.Nil(t, err)
	assert.NotNil(t, mongoDBWriter)

	mongoDBReader, err = dialMongoDB(&conf.DocumentDB.Reader)
	assert.Nil(t, err)
	err = refreshMongoDBConnection(mongoDBReader, &conf.DocumentDB.Reader)
	assert.Nil(t, err)
	err = disconnectMongoDBClient(mongoDBReader)
	assert.Nil(t, err)
	err = refreshMongoDBConnection(mongoDBReader, &conf.DocumentDB.Reader)
	assert.Nil(t, err)
	assert.NotNil(t, mongoDBReader)
	err = disconnectMongoDBClient(mongoDBReader)
	assert.Nil(t, err)
	err = refreshMongoDBConnection(mongoDBReader, &conf.DocumentDB.Reader)
	assert.Nil(t, err)
	assert.NotNil(t, mongoDBReader)
	err = disconnectMongoDBClient(mongoDBReader)
	assert.Nil(t, err)

	mongoDBWriter, err = dialMongoDB(&conf.DocumentDB.Writer)
	assert.Nil(t, err)
	err = refreshMongoDBConnection(mongoDBWriter, &conf.DocumentDB.Writer)
	assert.Nil(t, err)
	err = disconnectMongoDBClient(mongoDBWriter)
	assert.Nil(t, err)
	err = refreshMongoDBConnection(mongoDBWriter, &conf.DocumentDB.Writer)
	assert.Nil(t, err)
	assert.NotNil(t, mongoDBWriter)
	err = disconnectMongoDBClient(mongoDBWriter)
	assert.Nil(t, err)
	err = refreshMongoDBConnection(mongoDBWriter, &conf.DocumentDB.Writer)
	assert.Nil(t, err)
	assert.NotNil(t, mongoDBWriter)
	err = disconnectMongoDBClient(mongoDBWriter)
	assert.Nil(t, err)
}

func TestAssignMongoDBClient(t *testing.T) {
	assignMongoDBClient(nil, &conf.DocumentDB.Reader)
	assert.Nil(t, mongoDBReader)
	assignMongoDBClient(nil, &conf.DocumentDB.Writer)
	assert.Nil(t, mongoDBWriter)

	newReader, err := dialMongoDB(&conf.DocumentDB.Reader)
	assert.Nil(t, err)
	assignMongoDBClient(newReader, &conf.DocumentDB.Reader)
	assert.NotNil(t, mongoDBReader)

	newWriter, err := dialMongoDB(&conf.DocumentDB.Writer)
	assert.Nil(t, err)
	assignMongoDBClient(newWriter, &conf.DocumentDB.Writer)
	assert.NotNil(t, mongoDBWriter)
}
