package main

import (
    "github.com/cojac/tango"
)

// Index Page
type IndexHandler struct{ tango.BaseHandler }

func (h *IndexHandler) New() tango.HandlerInterface {
    return &IndexHandler{}
}

func (h *IndexHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
    return tango.NewHttpResponse("Hello! Now visit a non existent <a href=\"/404/\">page</a>.")
}

// Custom 404 Page
type NotFoundHandler struct{ tango.BaseHandler }

func (h *NotFoundHandler) New() tango.HandlerInterface {
    return &NotFoundHandler{}
}

func (h *NotFoundHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
    return tango.NewHttpResponse("Sorry... this page does not not exist. Don't forget to set this handler as <strong>404</strong>!!", 404)
}

func init() {
    tango.Settings.Set("debug", true)
    tango.Settings.Set("serve_address", ":8000")

    tango.Pattern("/", &IndexHandler{})
    tango.SetNotFoundHandler(&NotFoundHandler{})
}

func main() {
    tango.ListenAndServe()
}
