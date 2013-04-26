package main

import (
    "github.com/cojac/tango"
    "github.com/cojac/tango/middleware"
)

// Setup our only handler.
type IndexHandler struct {
    tango.BaseHandler
}

func (h IndexHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
    return tango.NewHttpResponse("Hello, world")
}

func main() {
    tango.Pattern("/", IndexHandler{})
    tango.Pattern("/hello/world/", IndexHandler{})
    tango.Pattern("/reg/{ex}/{id:[0-9]+}/", IndexHandler{})

    tango.Middleware(middleware.RuntimeProfile{})

    tango.ListenAndServe()
}
