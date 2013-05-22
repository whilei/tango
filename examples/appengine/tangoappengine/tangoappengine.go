// Run cmd (from one level up): 'dev_appserver.py app.yaml'
package tangoappengine

import (
    "github.com/cojac/tango"
)

type IndexHandler struct{ tango.BaseHandler }

func (h *IndexHandler) New() tango.HandlerInterface {
    return &IndexHandler{}
}

func (h *IndexHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
    return tango.NewHttpResponse("Hello, appengine.")
}

func init() {
    tango.Settings.Set("debug", true)

    tango.Pattern("/", &IndexHandler{})
}

// For AppEngine, just leave out the main func().
// func main() {
//     tango.ListenAndServe()
// }
