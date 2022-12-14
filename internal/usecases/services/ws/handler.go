package ws

import (
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/internal/oapi"
)

func (c *Client) callEventHandler(req *oapi.WsRequest) error {
	c.send <- &oapi.WsResponse{Type: oapi.WsResponseType(req.Type), Body: oapi.WsResponse_Body{}}

	return nil
}
