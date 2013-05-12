// Run cmd: go run settings_combo.go settings.json
package main

import (
    "fmt"
    "github.com/cojac/tango"
    "os"
    "path/filepath"
)

type IndexHandler struct {
    tango.BaseHandler
}

func (h *IndexHandler) New() tango.HandlerInterface {
    return &IndexHandler{}
}

func (h *IndexHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
    inline := fmt.Sprintf("Inline: <b>%s</b>", tango.Settings.String("inline"))
    json := fmt.Sprintf("A json val: <b>%s</b>", tango.Settings.String("project_name"))
    env := fmt.Sprintf("DB Pass from Env: <b>%s</b>", tango.Settings.String("db_pass"))
    bad := fmt.Sprintf("Not Set: <b>%s</b>", tango.Settings.String("not_set", "but we supplied a default!"))

    output := fmt.Sprintf("%s<br />%s<br />%s<br />%s", inline, json, env, bad)

    return tango.NewHttpResponse(output)
}

func init() {
    tango.Settings.Set("inline", "test123")

    if len(os.Args) < 2 {
        fmt.Printf("Usage: %s <settings-path>\n", filepath.Base(os.Args[0]))
        os.Exit(1)
    }

    tango.Settings.LoadFromFile(filepath.Base(os.Args[1]))

    tango.Settings.SetFromEnv("gopath", "GOPATH")
    tango.Settings.SetFromEnv("db_pass", "DB_PASSWORD", "default_password")
}

func main() {
    tango.Pattern("/", &IndexHandler{})
    tango.ListenAndServe()
}
