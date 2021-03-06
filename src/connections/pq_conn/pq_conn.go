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
	DbName           string
	Table            string
	Log              *logger.Logger
	Conn             *sql.DB
	Listener         *pq.Listener
	Txn              *sql.Tx
	connectionOpen   bool
	connectionBroken bool
	batchIngested    bool
}

// MakeConnection : build PqConnection object with connection to DB
func MakeConnection(dbname, table string) *PqConnection {
	conn_opts := fmt.Sprintf(
		"user=swarm64 password=swarm64 dbname=%v sslmode=disable", dbname)
	db_conn, _ := sql.Open("postgres", conn_opts) // no conn occurs, err always nil
	log := logger.InitLog("PQ_CONN")

	conn := PqConnection{
		DbName: dbname,
		Table:  table,
		Log:    &log,
		Conn:   db_conn,
		Listener: pq.NewListener(
			conn_opts,
			time.Duration(50)*time.Millisecond,
			time.Duration(1)*time.Second,
			listenerCallback(&log)),
		Txn:              nil,
		connectionOpen:   false,
		connectionBroken: false,
		batchIngested:    false}

	return &conn
}

// CheckConnectionState : Ping db using PqConnection.Listener to establish connection state
func (conn *PqConnection) CheckConnectionState() {
	if err := conn.Listener.Ping(); err != nil {
		if errStr := err.Error(); errStr == "no connection" {
			conn.Log.LogError(fmt.Sprintf("No currently active connection"))
		} else {
			conn.Log.LogError(fmt.Sprintf("Connection has issues: %v", err))
			conn.connectionBroken = true
		}
		conn.connectionOpen = false
	} else {
		conn.connectionOpen = true
	}
}

// OpenTransaction : Opens up a new transaction with the DB
func (conn *PqConnection) OpenTransaction() {
	txn, err := conn.Conn.Begin()
	if err != nil {
		conn.Log.LogError("Could not open transaction")
		errStr := fmt.Sprintf("%v", err.Error())
		conn.Log.LogError(errStr)
	}
	// allow some time for the transaction to connect
	// 3 ms seems enough, 1&2 ms too short
	time.Sleep(3 * time.Millisecond)
	conn.CheckConnectionState()
	conn.Txn = txn
}

// IngestData : takes data and ingests it into DB, returning true, or, returns false on errors
func (conn *PqConnection) IngestData(data []generator.DataFields) {
	conn.batchIngested = false
	conn.connectionOpen = false
	defer func() {
		switch exit := conn.Txn.Rollback(); exit.Error() {
		case "sql: Transaction has already been committed or rolled back":
			conn.Log.LogError("Successful transaction.")
		default:
			conn.Log.LogError(fmt.Sprintf("Rollback Message: %v", exit.Error()))
		}
		if r := recover(); r != nil {
			conn.Log.LogError(fmt.Sprintln(r))
		} else {
			conn.batchIngested = true
		}
		conn.CheckConnectionState()
	}()
	// time each of these and make an inline version of this and time that too.
	stmnt := conn.prepareStatement(data)
	conn.applyStatement(stmnt)
	conn.commitTxn()
	return
}

// SelectTimeStamps : select timestamps out of ingested data
func (conn *PqConnection) SelectTimeStamps() (tss []time.Time) {
	for !conn.connectionOpen {
		conn.CheckConnectionState()
	}
	var tmpT *time.Time
	rows, err := conn.Conn.Query(
		fmt.Sprintf("SELECT ts FROM %v ORDER BY ts", conn.Table))
	if err != nil {
		conn.Log.LogError(
			fmt.Sprintf("Error in select ts: %v", err.Error()))
		return
	}
	for rows.Next() {
		if err := rows.Scan(&tmpT); err != nil {
			conn.Log.LogError(
				fmt.Sprintf("Could not scan row: %v", err.Error()))
		} else {
			tss = append(tss, *tmpT)
		}
	}
	return
}

// dropAllRows : cleans the table for another run
func (conn PqConnection) DropAllRows() bool {
	_, err := conn.Conn.Exec(
		fmt.Sprintf("DELETE FROM %v", conn.Table))
	if err != nil {
		conn.Log.LogError(
			fmt.Sprintf("Rows may remain in the DB: %v", err.Error()))
		return false
	}
	return true
}

// IngestData helper methods - all panics caught in InjestData defer
func (conn PqConnection) commitTxn() {
	if exit := conn.Txn.Commit(); exit != nil {
		panic(
			fmt.Sprintf("Problem committing transaction: %v", exit.Error()))
	}
}

func (conn PqConnection) prepareStatement(data []generator.DataFields) *sql.Stmt {
	var errStr string
	stmnt, err := conn.Txn.Prepare(pq.CopyIn(conn.Table, "src", "dst", "flags"))
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

func (conn PqConnection) BatchIngested() bool {
	return conn.batchIngested
}

func (conn PqConnection) ConnectionBroken() bool {
	return conn.connectionBroken
}

func (conn PqConnection) ConnectionOpen() bool {
	return conn.connectionOpen
}

// PqConnection Listener Callback
func listenerCallback(l *logger.Logger) func(event pq.ListenerEventType, err error) {
	return func(event pq.ListenerEventType, err error) {
		switch event {
		case pq.ListenerEventDisconnected:
			l.LogError(
				fmt.Sprintf("Connection Disconnected: %v", err.Error()))
		case pq.ListenerEventReconnected:
			l.LogError(fmt.Sprintf("Connection to db re-established"))
		case pq.ListenerEventConnectionAttemptFailed:
			l.LogError(fmt.Sprintf("Could not establish DB connection"))
		case pq.ListenerEventConnected:
			l.LogError(fmt.Sprintf("DB connection established"))
		default:
			l.LogError(fmt.Sprintf("Unhandled DB response: %v", event))
		}
	}
}
