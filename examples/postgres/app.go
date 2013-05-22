package main

import (
    "database/sql"
    "fmt"
    _ "github.com/bmizerany/pq"
    "github.com/cojac/tango"
    "time"
)

type IndexHandler struct {
    tango.BaseHandler
    PostgresMixin // This allows us to reference 'DB' within our handler (see below).
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
    tango.Mixin(&PostgresMixin{})

    tango.Pattern("/", &IndexHandler{})
}

func main() {
    tango.ListenAndServe()
}

// ---------------------------------------------------------------------------------
// This should probably live somewhere else. I'll leave it in the examples for now.
// But ideally we should have some kind of "module/mixin/middleware" repo (or individual repos?)

// Postgres Mixin
type PostgresMixin struct {
    tango.BaseMixin
    DB  *sql.DB
}

func (m *PostgresMixin) InitMixin() {
    err := DbPool.InitPool(tango.Settings.Int("db_pool_size", 3), initPostgresConnection)
    if err != nil {
        tango.LogError.Panicln("Database init error:", err)
    }
}

func (m *PostgresMixin) PreparePostgresMixin() {
    m.DB = DbPool.GetConn().(*sql.DB)
}

func (m *PostgresMixin) FinishPostgresMixin() {
    DbPool.ReleaseConn(m.DB)
}

func initPostgresConnection() (interface{}, error) {
    tango.LogInfo.Println("Creating Postgres Connection")
    conf := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=%s",
        tango.Settings.String("db_name"),
        tango.Settings.String("db_user"),
        tango.Settings.String("db_password", ""),
        tango.Settings.String("db_host", "127.0.0.1"),
        tango.Settings.Int("db_port", 5432),
        tango.Settings.String("db_sslmode", "disable"))
    return sql.Open("postgres", conf)
}

// Setup the connection pool.
var DbPool = &ConnectionPoolWrapper{}

type InitFunction func() (interface{}, error)

type ConnectionPoolWrapper struct {
    size int
    conn chan interface{}
}

func (p *ConnectionPoolWrapper) InitPool(size int, initfn InitFunction) error {
    // Create a buffered channel allowing size senders
    p.conn = make(chan interface{}, size)
    for x := 0; x < size; x++ {
        conn, err := initfn()
        if err != nil {
            return err
        }

        p.conn <- conn
    }
    p.size = size
    return nil
}

func (p *ConnectionPoolWrapper) GetConn() interface{} {
    return <-p.conn
}

func (p *ConnectionPoolWrapper) ReleaseConn(conn interface{}) {
    p.conn <- conn
}
