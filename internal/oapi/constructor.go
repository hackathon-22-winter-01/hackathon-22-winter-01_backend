package oapi

func WsResponseFromType(typ WsResponseType) *WsResponse {
	return &WsResponse{
		Type: typ,
	}
}
