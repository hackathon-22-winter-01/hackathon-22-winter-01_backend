package oapi

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/assert"
	"github.com/shiguredo/websocket"
	"github.com/stretchr/testify/require"
)

// WsResponse for tests
type WsResponseWrapper[T any] struct {
	T    *testing.T
	Res  *WsResponse
	Body T
}

func (w *WsResponseWrapper[T]) Equal(expectedType WsResponseType, expectedBody T, opts ...cmp.Option) {
	w.T.Helper()

	assert.NoError(w.T, json.Unmarshal(w.Res.Body.union, &w.Body))

	assert.Equal(w.T, expectedType, w.Res.Type)
	assert.Equal(w.T, expectedBody, w.Body, opts...)
}

// WsRequest for tests
type WsRequestWrapper struct {
	T   *testing.T
	Req *WsRequest
}

func NewWsRequestForTest(t *testing.T, typ WsRequestType, body any) *WsRequestWrapper {
	t.Helper()

	b, err := json.Marshal(body)
	assert.NoError(t, err)

	return &WsRequestWrapper{
		T: t,
		Req: &WsRequest{
			Type: typ,
			Body: WsRequest_Body{union: b},
		},
	}
}

func WriteWsRequest(t *testing.T, c *websocket.Conn, typ WsRequestType, body any) {
	t.Helper()

	b, err := json.Marshal(body)
	assert.NoError(t, err)

	require.NoError(t, c.WriteJSON(WsRequest{
		Type: typ,
		Body: WsRequest_Body{union: b},
	}))
}
