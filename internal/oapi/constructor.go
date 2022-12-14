package oapi

func NewWsResponse(typ WsResponseType) *WsResponse {
	return &WsResponse{
		Type: typ,
	}
}
