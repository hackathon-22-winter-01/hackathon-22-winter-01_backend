package ws_test

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/stretchr/testify/require"
)

func readWsResponse(t *testing.T, c *websocket.Conn) *oapi.WsResponse {
	t.Helper()

	res := new(oapi.WsResponse)
	require.NoError(t, c.ReadJSON(res))

	return res
}
