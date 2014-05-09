package main

import (
	"github.com/unrolled/tango"
)

type IndexHandler struct{ tango.BaseHandler }

func (h *IndexHandler) New() tango.HandlerInterface {
	return &IndexHandler{}
}

func (h *IndexHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
	return tango.NewHttpResponse("Hello, world.")
}

func init() {
	tango.Settings.Set("debug", true)
	tango.Settings.Set("serve_address", ":8000")

	tango.Pattern("/", &IndexHandler{})
}

func main() {
	tango.ListenAndServe()
}
