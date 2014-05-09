// Run cmd: go run settings_env.go
package main

import (
	"fmt"

	"github.com/unrolled/tango"
)

type IndexHandler struct {
	tango.BaseHandler
}

func (h *IndexHandler) New() tango.HandlerInterface {
	return &IndexHandler{}
}

func (h *IndexHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
	debug := fmt.Sprintf("Debug: <b>%v</b>", tango.Debug)
	gopath := fmt.Sprintf("Go Path: <b>%s</b>", tango.Settings.String("gopath", "why this not set?"))
	bad := fmt.Sprintf("Not Set: <b>%s</b>", tango.Settings.String("not_set", "but we supplied a default!"))

	output := fmt.Sprintf("%s<br />%s<br />%s", debug, gopath, bad)

	return tango.NewHttpResponse(output)
}

func init() {
	tango.Settings.Set("serve_address", ":8000")

	tango.Settings.SetFromEnv("gopath", "GOPATH")
	tango.Settings.SetFromEnv("debug", "DOES_NOT_EXIST_SO_DEFAULT_TRUE", true)
}

func main() {
	tango.Pattern("/", &IndexHandler{})
	tango.ListenAndServe()
}
