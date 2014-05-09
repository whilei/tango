// Run cmd: go run settings_json.go
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
	// Path is relative. Place this file and settings.json in the same directory.
	tango.Settings.LoadFromFile("settings.json")
}

func main() {
	tango.Pattern("/", &IndexHandler{})
	tango.ListenAndServe()
}
