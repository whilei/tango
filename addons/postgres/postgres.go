package postgres

import (
    "database/sql"
    "fmt"
    "github.com/cojac/tango"
    _ "github.com/lib/pq"
    "runtime"
)

type PostgresMixin struct {
    tango.BaseMixin
    Db  *sql.DB
}

func (m *PostgresMixin) InitMixin() {
    err := DbPool.InitPool(tango.Settings.Int("db_pool_size", runtime.NumCPU()*2), InitPostgresConnection)
    if err != nil {
        tango.LogError.Panicln("Database init error:", err)
    }
}

func (m *PostgresMixin) PreparePostgresMixin() {
    m.Db = DbPool.GetConn().(*sql.DB)
}

func (m *PostgresMixin) FinishPostgresMixin() {
    DbPool.ReleaseConn(m.Db)
}

func InitPostgresConnection() (interface{}, error) {
    tango.LogDebug.Println("Creating Postgres Connection")

    // Used if we supply a full string in our conf.
    conf := tango.Settings.String("db_dsn", "")

    if conf == "" {
        conf = fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=%s",
            tango.Settings.String("db_name"),
            tango.Settings.String("db_user"),
            tango.Settings.String("db_password", ""),
            tango.Settings.String("db_host", "127.0.0.1"),
            tango.Settings.Int("db_port", 5432),
            tango.Settings.String("db_sslmode", "disable"))
    }

    db, err := sql.Open("postgres", conf)

    return db, err
}

// TODO: This needs some clean up yet... and maybe a testcase or two.
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