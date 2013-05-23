package main

import (
    "fmt"
    "github.com/cojac/tango"
    "github.com/cojac/tango-addons/postgres"
    "time"
)

type IndexHandler struct {
    tango.BaseHandler
    postgres.PostgresMixin // This allows us to reference 'DB' within our handler (see below).
}

func (h *IndexHandler) New() tango.HandlerInterface {
    return &IndexHandler{}
}

func (h *IndexHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
    r, err := h.DB.Query("SELECT clock_timestamp()")
    if err != nil {
        tango.LogError.Panicf("DB 1 error:", err)
    }
    defer r.Close()

    if !r.Next() {
        tango.LogError.Panicf("DB 2 error:", err)
    }

    var datetime time.Time
    err = r.Scan(&datetime)
    if err != nil {
        tango.LogError.Panicf("DB 3 error:", err)
    }

    return tango.NewHttpResponse(fmt.Sprintf("Postgres says the timestamp is: %s", datetime))
}

func init() {
    tango.Settings.Set("debug", true)
    tango.Settings.Set("serve_address", ":8000")

    // DB Settings
    tango.Settings.Set("db_pool_size", 3)
    tango.Settings.Set("db_name", "postgres")
    tango.Settings.Set("db_user", "postgres")
    tango.Settings.Set("db_password", "")
    tango.Settings.Set("db_host", "127.0.0.1")
    tango.Settings.Set("db_port", 5432)
    tango.Settings.Set("db_sslmode", "disable")

    // Add the Postgres mixin to our app.
    tango.Mixin(&postgres.PostgresMixin{})

    tango.Pattern("/", &IndexHandler{})
}

func main() {
    tango.ListenAndServe()
}
