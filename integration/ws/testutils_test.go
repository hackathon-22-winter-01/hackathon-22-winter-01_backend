package ws_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/usecases/services/ws"
	"github.com/shiguredo/websocket"
	"github.com/stretchr/testify/require"
)

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

func readWsResponse[T any](t *testing.T, c *websocket.Conn) *oapi.WsResponseWrapper[T] {
	t.Helper()

	var w oapi.WsResponseWrapper[T]

	res := new(oapi.WsResponse)
	require.NoError(t, c.ReadJSON(res))

	w.T = t
	w.Res = res

	return &w
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

func connectToWs(t *testing.T, streamer ws.Streamer, playerID uuid.UUID) *websocket.Conn {
	server := httptest.NewServer(&httpHandler{t, streamer, playerID})
	server.URL = "ws" + strings.TrimPrefix(server.URL, "http")
	c, _, err := websocket.DefaultDialer.Dial(server.URL, nil)
	require.NoError(t, err)

	return c
}

// randint 乱数であることを明示する
type randint = int
