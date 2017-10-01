package pq_conn

import (
	"database/sql"
	"fmt"
	"generator"
	"logger"
	"time"

	pq "github.com/lib/pq"
)

// PqConnection : simple DB connection model
type PqConnection struct {
	Log              logger.Logger
	Conn             *sql.DB
	Listener         *pq.Listener
	Txn              *sql.Tx
	ConnectionOpen   bool
	ConnectionBroken bool
}

// SwarmRow : row return struct
type SwarmRow struct {
	Ts    *time.Time
	Src   *int64
	Dst   *int64
	Flags *int64
}

// MakeConnection : build PqConnection object with connection to DB
func MakeConnection(dbname string) *PqConnection {
	// conn_url := "postgres://swarm64:swarm64@localhost/swarmtest?sslmode=require"
	conn_opts := fmt.Sprintf(
		"user=swarm64 password=swarm64 dbname=%v sslmode=disable", dbname)
	db_conn, _ := sql.Open("postgres", conn_opts) // no conn occurs, err always nil
	log := logger.InitLog("PQ_CONN")

	conn := PqConnection{
		Log:  log,
		Conn: db_conn,
		Listener: pq.NewListener(
			conn_opts,
			time.Duration(50)*time.Millisecond,
			time.Duration(1)*time.Second,
			listenerCallback(&log)),
		Txn:              nil,
		ConnectionOpen:   false,
		ConnectionBroken: false}
	return &conn
}

// CheckConnectionState : Ping db using PqConnection.Listener to establish connection state
func (conn *PqConnection) CheckConnectionState() {
	if err := conn.Listener.Ping(); err != nil {
		if errStr := err.Error(); errStr == "no connection" {
			conn.Log.LogError(fmt.Sprintf("No currently active connection\n"))
		} else {
			conn.Log.LogError(fmt.Sprintf("Connection has issues: %v\n", err))
			conn.ConnectionBroken = true
		}
		conn.ConnectionOpen = false
	} else {
		conn.ConnectionOpen = true
	}
}

// OpenTransaction : Opens up a new transaction with the DB
func (conn *PqConnection) OpenTransaction() {
	txn, err := conn.Conn.Begin()
	if err != nil {
		conn.Log.LogError("Could not open transaction")
		errStr := fmt.Sprintf("%v\n", err.Error())
		conn.Log.LogError(errStr)
	}
	// allow some time for the transaction to connect, 2 ms seems enough, 1 ms too short
	time.Sleep(2 * time.Millisecond)
	conn.CheckConnectionState()
	conn.Txn = txn
}

// IngestData : takes data and ingests it into DB, returning true, or, returns false on errors
func (conn *PqConnection) IngestData(data []generator.DataFields) (complete bool) {
	defer func() {
		if exit := conn.Txn.Rollback(); exit != nil {
			conn.Log.LogError(fmt.Sprintf("Rollback Message: %v", exit.Error()))
		}
		if r := recover(); r != nil {
			conn.Log.LogError(fmt.Sprintln(r))
			conn.CheckConnectionState()
			complete = false
		} else {
			complete = true
		}
	}()
	stmnt := conn.prepareStatement(data)
	conn.applyStatement(stmnt)
	if exit := conn.Txn.Commit(); exit != nil {
		panic(
			fmt.Sprintf("Problem committing transaction: %v", exit.Error()))
	}
	return
}

// IngestData helper methods - all panics caught in InjestData defer
func (conn PqConnection) prepareStatement(data []generator.DataFields) *sql.Stmt {
	var errStr string
	stmnt, err := conn.Txn.Prepare(pq.CopyIn("ingestion_test", "src", "dst", "flags"))
	if err != nil {
		errStr = fmt.Sprintf("Could not generate transaction statement: %v", err.Error())
		panic(errStr)
	}

	for _, dat := range data {
		_, err := stmnt.Exec(dat.Src, dat.Dst, dat.Flags)
		if err != nil {
			errStr = fmt.Sprintf("Problem adding %v to txn statement: %v", dat, err.Error())
			if exit := stmnt.Close(); exit != nil {
				errStr = fmt.Sprintf(
					"%v\n\t\tProblem closing transaction statement: %v", errStr, exit.Error())
			}
			panic(errStr)
		}
	}
	return stmnt
}

func (conn PqConnection) applyStatement(stmnt *sql.Stmt) {
	var errStr string
	if _, execExit := stmnt.Exec(); execExit != nil {
		errStr = fmt.Sprintf("Problem submitting transaction statement: %v", execExit.Error())
		if exit := stmnt.Close(); exit != nil {
			errStr = fmt.Sprintf(
				"%v\n\t\tProblem closing transaction statement: %v", errStr, exit.Error())
		}
		panic(errStr)
	}
	if exit := stmnt.Close(); exit != nil {
		errStr = fmt.Sprintf("Problem closing transaction statement: %v", exit.Error())
		panic(errStr)
	}
}

// PqConnection Listener Callback
func listenerCallback(l *logger.Logger) func(event pq.ListenerEventType, err error) {
	return func(event pq.ListenerEventType, err error) {
		switch event {
		case pq.ListenerEventDisconnected:
			l.LogError(
				fmt.Sprintf("Connection Disconnected: %v\n", err.Error()))
		case pq.ListenerEventReconnected:
			l.LogError(fmt.Sprintf("Connection to db re-established\n"))
		case pq.ListenerEventConnectionAttemptFailed:
			l.LogError(fmt.Sprintf("Could not establish DB connection\n"))
		case pq.ListenerEventConnected:
			l.LogError(fmt.Sprintf("DB connection established\n"))
		}
	}
}
