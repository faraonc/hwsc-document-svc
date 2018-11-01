package service

import (
	pb "github.com/faraonc/hwsc-api-blocks/int/hwsc-document-svc/proto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
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
