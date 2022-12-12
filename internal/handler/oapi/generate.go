//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config=oapi_server.yml ../../../docs/openapi/v1.yml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config=oapi_models.yml ../../../docs/openapi/v1.yml

package oapi
