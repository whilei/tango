package tango

import (
    "database/sql"
    "fmt"
    _ "github.com/bmizerany/pq"
)

// Postgres Mixin
type PostgresMixin struct {
    BaseMixin
    DB  *sql.DB
}

func (m *PostgresMixin) InitMixin() {
    err := DbPool.InitPool(Settings.Int("db_pool_size", 3), initPostgresConnection)
    if err != nil {
        LogError.Panicln("Database init error:", err)
    }
}

func (m *PostgresMixin) PreparePostgresMixin() {
    m.DB = DbPool.GetConn().(*sql.DB)
}

func (m *PostgresMixin) FinishPostgresMixin() {
    DbPool.ReleaseConn(m.DB)
}

func initPostgresConnection() (interface{}, error) {
    LogInfo.Println("Creating Postgres Connection")
    conf := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=%s",
        Settings.String("db_name"),
        Settings.String("db_user"),
        Settings.String("db_password", ""),
        Settings.String("db_host", "127.0.0.1"),
        Settings.Int("db_port", 5432),
        Settings.String("db_sslmode", "disable"))
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
