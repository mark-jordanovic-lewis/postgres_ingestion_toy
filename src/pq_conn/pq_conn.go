package pq_conn

import (
	"database/sql"
	"fmt"
	"logger"
	"os"
	"time"

	pq "github.com/lib/pq"
)

type PqConnection struct {
	Log      logger.Logger
	Conn     *sql.DB
	Listener *pq.Listener
}

func MakeConnection() *PqConnection {
	conn_url := "postgres://maruko:@localhost/swarmtest?sslmode=disable"
	conn_opts := "user=maruko dbname=swarmtest sslmode=disable"
	db_conn, _ := sql.Open("postgres", conn_url) // no conn occurs, err always nil
	log := logger.InitLog("pq_conn")

	conn := PqConnection{
		Log:  log,
		Conn: db_conn,
		Listener: pq.NewListener(
			conn_url,
			time.Duration(50)*time.Millisecond,
			time.Duration(1)*time.Second,
			listenerCallback(&log))}
	// if !conn.PingDB() {
	// 	return nil
	// }
	return &conn
}

func (conn PqConnection) PingDB() bool {
	if err := conn.Listener.Ping(); err != nil {
		conn.Log.LogError(
			fmt.Sprintf(
				"Connection could not be established\n%v\n",
				err))
		return false
	}
	return true
}

func (conn PqConnection) OpenTransaction() *sql.Tx {
	txn, err := conn.Conn.Begin()
	if err != nil {
		conn.Log.LogError("Could not open transaction")
		errStr := fmt.Sprintf("%v\n", err.Error())
		conn.Log.LogError(errStr)
		os.Exit(1)
	}
	return txn
}

// just to keep the libs imported
func tmp_pq(txn *sql.Tx) {
	txn.Prepare(pq.CopyIn("ingestion_test", "src", "dst", "flags"))
}

func listenerCallback(l *logger.Logger) func(event pq.ListenerEventType, err error) {
	return func(event pq.ListenerEventType, err error) {
		switch event {
		case pq.ListenerEventDisconnected:
			l.LogError(
				fmt.Sprintf("Connection Disconnected %v\n%v\n", time.Now(), err.Error()))
		case pq.ListenerEventReconnected:
			l.LogError(fmt.Sprintf("Connection to db re-established %v\n", time.Now()))
		case pq.ListenerEventConnectionAttemptFailed:
			l.LogError(fmt.Sprintf("Could not establish DB connection %v\n", time.Now()))
		}
	}
}
