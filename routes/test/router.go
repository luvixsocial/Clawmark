package test

import (
	"clawmark/uapi"

	"github.com/go-chi/chi/v5"
)

type Router struct{}

func (b Router) Tag() (string, string) {
	return "Test", "Hello, there. This category of test endpoints that allow our developers to test the core of our API."
}

func (b Router) Routes(r *chi.Mux) {
	uapi.Route{
		Pattern: "/test",
		OpId:    "test",
		Method:  uapi.GET,
		Docs:    TestDocs,
		Handler: TestRoute,
	}.Route(r)
}
