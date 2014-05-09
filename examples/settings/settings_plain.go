// Run cmd: go run settings_plain.go
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
	project := fmt.Sprintf("Project Name: <b>%s</b>", tango.Settings.String("project_name"))
	number := fmt.Sprintf("Some Number: <b>%d</b>", tango.Settings.Int("some_number"))
	other := fmt.Sprintf("Some Other Number: <b>%f</b>", tango.Settings.Float("some_other_number"))
	ready := fmt.Sprintf("Is Production Ready: <b>%v</b>", tango.Settings.Bool("is_production_ready"))
	bad := fmt.Sprintf("Not Set: <b>%s</b>", tango.Settings.String("not_set", "but we supplied a default!"))

	output := fmt.Sprintf("%s<br />%s<br />%s<br />%s<br />%s", project, number, other, ready, bad)

	return tango.NewHttpResponse(output)
}

func init() {
	// "debug" is the only special var that will be accessible via "tango.Debug".
	tango.Settings.Set("debug", true)
	tango.Settings.Set("serve_address", ":8000")
	tango.Settings.Set("project_name", "Tango")
	tango.Settings.Set("some_number", 123)
	tango.Settings.Set("some_other_number", 456.78)
	tango.Settings.Set("is_production_ready", false)
}

func main() {
	tango.Pattern("/", &IndexHandler{})
	tango.ListenAndServe()
}
