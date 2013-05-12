// Run cmd: go run settings_flags.go -debug=true -num=1234 -serve=0.0.0.0:8000
package main

import (
    "flag"
    "fmt"
    "github.com/cojac/tango"
)

type IndexHandler struct {
    tango.BaseHandler
}

func (h *IndexHandler) New() tango.HandlerInterface {
    return &IndexHandler{}
}

func (h *IndexHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
    debug := fmt.Sprintf("Debug: <b>%v</b>", tango.Debug)
    number := fmt.Sprintf("Some Number: <b>%d</b>", tango.Settings.Int("some_number"))
    bad := fmt.Sprintf("Not Set: <b>%s</b>", tango.Settings.String("not_set", "but we supplied a default!"))

    output := fmt.Sprintf("%s<br />%s<br />%s", debug, number, bad)

    return tango.NewHttpResponse(output)
}

func init() {
    // Path is relative. Place this file and settings.json in the same directory.
    serve := flag.String("serve", ":8000", "Address to serve this app.")
    debug := flag.Bool("debug", false, "Should we debug?")
    random := flag.Int("num", -1, "Random int number.")
    flag.Parse()

    tango.Settings.Set("debug", *debug)
    tango.Settings.Set("serve_address", *serve)
    tango.Settings.Set("some_number", *random)
}

func main() {
    tango.Pattern("/", &IndexHandler{})
    tango.ListenAndServe()
}
