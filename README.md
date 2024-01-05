# Gitlab Service



## Getting started

To Generate Server-side boilerplate from OpenAPI specifications :
```bash
go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
oapi-codegen -package api -generate gorilla ./specs/api.yml > ./cmd/api/gitlab-service.gen.go
```