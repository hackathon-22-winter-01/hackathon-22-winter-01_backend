package ws_test

import (
	"net/http"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/stretchr/testify/require"
)

func mustWriteWsRequest(t *testing.T, c *websocket.Conn, typ oapi.WsRequestType, body oapi.WsRequest_Body) {
	t.Helper()

	require.NoError(t, c.WriteJSON(oapi.WsRequest{
		Type: typ,
		Body: body,
	}))
}

func forEachClientAsync(t *testing.T, wg *sync.WaitGroup, c []*websocket.Conn, f func(int, *websocket.Conn)) {
	t.Helper()

	for i, c := range c {
		wg.Add(1)

		go func(i int, c *websocket.Conn) {
			defer wg.Done()
			f(i, c)
		}(i, c)
	}

	wg.Wait()
}

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
