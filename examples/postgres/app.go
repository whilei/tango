package main

import (
    "database/sql"
    "fmt"
    "github.com/cojac/tango"
    _ "github.com/lib/pq"
    "time"
)

var Db *sql.DB

// Index Page
type IndexHandler struct {
    tango.BaseHandler
}

func (h *IndexHandler) New() tango.HandlerInterface {
    return &IndexHandler{}
}

func (h *IndexHandler) Get(request *tango.HttpRequest) *tango.HttpResponse {
    query := "SELECT clock_timestamp();"

    stmt, err := Db.Prepare(query)
    if err != nil {
        tango.LogError.Panicf("DB 1 error:", err)
    }
    defer stmt.Close()

    var datetime time.Time
    err = stmt.QueryRow().Scan(&datetime)
    if err != nil {
        tango.LogError.Panicf("DB 2 error:", err)
    }

    return tango.NewHttpResponse(fmt.Sprintf("Postgres clock_timestamp is: %s", datetime))
}

func init() {
    tango.Settings.Set("debug", true)
    tango.Settings.Set("serve_address", ":8000")

    // Setup PG Connection
    conf := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=%s",
        tango.Settings.String("db_name", "postgres"),
        tango.Settings.String("db_user", "postgres"),
        tango.Settings.String("db_password", ""),
        tango.Settings.String("db_host", "127.0.0.1"),
        tango.Settings.Int("db_port", 5432),
        tango.Settings.String("db_sslmode", "disable"),
    )

    var err error
    Db, err = sql.Open("postgres", conf)
    if err != nil {
        panic(err)
    }

    // Now setup as you normally would.
    tango.Pattern("/", &IndexHandler{})
}

func main() {
    tango.ListenAndServe()
}
