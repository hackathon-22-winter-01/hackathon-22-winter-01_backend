package ws_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/stretchr/testify/require"
)

func readWsResponse(t *testing.T, c *websocket.Conn) *oapi.WsResponse {
	t.Helper()

	res := new(oapi.WsResponse)
	require.NoError(t, c.ReadJSON(res))

	return res
}

type httpHandler struct {
	t   *testing.T
	s   ws.Streamer
	pid uuid.UUID
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.s.ServeWS(w, r, h.pid); err != nil {
		h.t.Error(err)
	}
}
