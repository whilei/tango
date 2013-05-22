package main

import (
    "fmt"
    "github.com/cojac/tango"
    "time"
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
    tango.Middleware(&RuntimeProfile{})

    tango.Pattern("/", &IndexHandler{})
}

func main() {
    tango.ListenAndServe()
}

// ---------------------------------------------------------------------------------
// This should probably live somewhere else. I'll leave it in the examples for now.
// But ideally we should have some kind of "module/mixin/middleware" repo (or individual repos?)

// RuntimeProfile Middleware
type RuntimeProfile struct {
    tango.BaseMiddleware
}

func (m *RuntimeProfile) ProcessRequest(request *tango.HttpRequest, response *tango.HttpResponse) {
    request.Registry[runTimeContextKey] = time.Now()
}

func (m *RuntimeProfile) ProcessResponse(request *tango.HttpRequest, response *tango.HttpResponse) {
    started := request.Registry[runTimeContextKey]
    response.Header.Set("X-Runtime", fmt.Sprintf("%s", time.Since(started.(time.Time))))
}

const runTimeContextKey string = "__middleware_run_time_profile_start_key__"
