package main

import (
    "github.com/cojac/tango"
    "github.com/cojac/tango/middleware"
)

// Using the NewHttpResponse method to generate an HttpResponse.
type NewHandler struct{ tango.BaseHandler }

func (h NewHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
    return tango.NewHttpResponse("Hello, world.")
}

// Using the literal HttpResponse invokation.
type LiteralHandler struct{ tango.BaseHandler }

func (h LiteralHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
    return &tango.HttpResponse{Content: "Hello, world... again!", StatusCode: 200}
}

func init() {
    tango.Settings.Set("debug", true)
    tango.Settings.Set("serve_address", ":8000")
}

func main() {
    tango.Pattern("/", NewHandler{})
    tango.Pattern("/lit/{id:[0-9]+}/", LiteralHandler{})

    tango.Middleware(middleware.RuntimeProfile{})

    tango.ListenAndServe()
}
