package service

import (
	"github.com/hwsc-org/hwsc-document-svc/conf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDialMongoDB(t *testing.T) {
	cases := []struct {
		uri      string
		isExpErr bool
		errorStr string
	}{
		{conf.DocumentDB.Reader, false, ""},
		{"", true, "error parsing uri (): scheme must be \"mongodb\" or \"mongodb+srv\""},
	}

	for _, c := range cases {
		client, err := dialMongoDB(c.uri)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
			assert.NotNil(t, client)
		}
	}
}

func TestDisconnectMongoDBClient(t *testing.T) {
	cases := []struct {
		uri      string
		isExpErr bool
		errorStr string
	}{
		{conf.DocumentDB.Reader, false, ""},
		{"", true, errNilMongoDBClient.Error()},
	}

	for _, c := range cases {
		client, _ := dialMongoDB(c.uri)
		err := disconnectMongoDBClient(client)
		if c.isExpErr {
			assert.EqualError(t, err, c.errorStr)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestRefreshMongoDBConnection(t *testing.T) {
	client, err := dialMongoDB(conf.DocumentDB.Reader)
	assert.Nil(t, err)
	err = refreshMongoDBConnection(client)
	assert.Nil(t, err)
	err = disconnectMongoDBClient(client)
	assert.Nil(t, err)
	err = refreshMongoDBConnection(client)
	assert.Nil(t, err)
	err = disconnectMongoDBClient(client)
	assert.Nil(t, err)
}
