package oapi

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hackathon-22-winter-01/hackathon-22-winter-01_backend/pkg/assert"
)

// for tests
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
