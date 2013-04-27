package main

import (
    "github.com/cojac/tango"
    "github.com/cojac/tango/middleware"
)

type IndexHandler struct {
    tango.BaseHandler
}

func (h IndexHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
    return tango.NewHttpResponse("Hello, world")
}

func init() {
    tango.Settings.Set("debug", true)
    tango.Settings.Set("serve_address", ":8000")
}

func main() {
    tango.Pattern("/", IndexHandler{})
    tango.Pattern("/hello/world/", IndexHandler{})
    tango.Pattern("/reg/{ex}/{id:[0-9]+}/", IndexHandler{})

    tango.Middleware(middleware.RuntimeProfile{})

    tango.ListenAndServe()
}
