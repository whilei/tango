package tangoappengine

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

    tango.Pattern("/", IndexHandler{})

    tango.Middleware(middleware.RuntimeProfile{})
}

// For AppEngine, just leave out the main func().
// func main() {
//     tango.ListenAndServe()
// }
