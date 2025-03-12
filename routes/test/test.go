package test

import (
	"clawmark/types"
	"clawmark/uapi"
	"net/http"

	docs "clawmark/doclib"
)

func TestDocs() *docs.Doc {
	return &docs.Doc{
		Summary:     "Test Documentation",
		Description: "This endpoint tests our documentation page.",
		Params:      []docs.Parameter{},
		Resp:        types.Response{},
	}
}

func TestRoute(d uapi.RouteData, r *http.Request) uapi.HttpResponse {
	msg := "Hello, there. This is a test endpoint."

	return uapi.HttpResponse{
		Status: http.StatusOK,
		Json: types.Response{
			Success: true,
			Message: &msg,
			JSON: map[string]interface{}{
				"username": "johndoe",
				"age":      25,
			},
		},
	}
}
