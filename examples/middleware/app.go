package main

import (
	"github.com/unrolled/tango"
	"github.com/unrolled/tango/addons/runtime"
)

type IndexHandler struct{ tango.BaseHandler }

func (h *IndexHandler) New() tango.HandlerInterface {
	return &IndexHandler{}
}

func (h *IndexHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
	return tango.NewHttpResponse("View the headers of this page... you will see 'X-Runtime' in the response headers.")
}

func init() {
	tango.Settings.Set("debug", true)
	tango.Settings.Set("serve_address", ":8000")

	// Add the Runtime Profiler to our middleware stack. Same ordering rules as Django!
	tango.Middleware(&runtime.Profiler{})

	tango.Pattern("/", &IndexHandler{})
}

func main() {
	tango.ListenAndServe()
}
