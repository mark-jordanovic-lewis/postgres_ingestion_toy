package pgx_conn

import (
	"database/sql"
	"fmt"
	gen "generator"
	"logger"
	"time"

	pgx "github.com/jackc/pgx"
)

// PqConnection : simple DB connection model - refactor to Connection in pq_conn too
type PgxConnection struct {
	DbName           string
	Table            pgx.Identifier
	Log              *logger.Logger
	Conn             *pgx.ConnPool
	Txn              *sql.Tx
	connectionOpen   bool
	connectionBroken bool
	batchIngested    bool
}

// MakeConnection : actually builds a connection pool to the DB and returns
// a newly instantiated PqxConnection object
func MakeConnection(dbname, table string, n_connections int) *PgxConnection {
	var tableId pgx.Identifier
	log := logger.InitLog("PGX_CONN")
	poolConfig := makeConnectionPool(makeConfig(dbname), n_connections)
	db_conn, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		log.LogError(
			fmt.Sprintf("Could not make connection pool: %v", err.Error()))
	}
	tableId = append(tableId, table)

	conn := PgxConnection{
		DbName:         dbname,
		Table:          tableId,
		Log:            &log,
		Conn:           db_conn,
		Txn:            nil,
		connectionOpen: false,
		batchIngested:  false}

	return &conn
}

// CheckConnectionState : Pings DB to check connection status
func (conn *PgxConnection) CheckConnectionState() {
	status := conn.Conn.Stat()
	// in the future use a slice for each conneciton in the pool
	conn.connectionBroken = status.MaxConnections == 0
	conn.connectionOpen = status.CurrentConnections > 0
	// status.AvailableConnections can be used to streamline
}

// CopyData : uses CopyFrom to move data into the DB
func (conn *PgxConnection) CopyData(data [][]interface{}) (float64, error) {
	conn.batchIngested = false
	cols := []string{"src", "dst", "flags"}
	source := pgx.CopyFromRows(data)

	start := time.Now().UnixNano()
	i, err := conn.Conn.CopyFrom(conn.Table, cols, source)
	end := time.Now().UnixNano()
	dt := float64(end - start)

	if err != nil {
		conn.Log.LogError(
			fmt.Sprintf("An error occurred in CopyData: %v", err.Error()))
	} else if i != len(data) {
		conn.Log.LogError(
			fmt.Sprintf("Not all Rows ingested: %v of %v", i, len(data)))
	} else {
		conn.Log.LogError(
			fmt.Sprintf("Rows ingested: %v of %v", i, len(data)))
		conn.batchIngested = true
	}
	return dt, err
}

// DropAllRows : removes all rows from table
func (conn *PgxConnection) DropAllRows() {
	del_conn, err := conn.Conn.Acquire()
	if err != nil {
		conn.Log.LogError(
			fmt.Sprintf("Could not acquire connection: %v", err.Error()))
	}
	_, err = del_conn.Exec(fmt.Sprintf("DELETE FROM %v", conn.Table[0]))
	if err != nil {
		conn.Log.LogError(
			fmt.Sprintf("Could not delete rows: %v", err.Error()))
	}
}

// getters: this is code smell, don't do this.
func (conn PgxConnection) BatchIngested() bool {
	return conn.batchIngested
}

func (conn PgxConnection) ConnectionBroken() bool {
	return conn.connectionBroken
}

func (conn PgxConnection) ConnectionOpen() bool {
	return conn.connectionOpen
}
func tmp_pgx(conn *pgx.ConnPool) {
	txn, err := conn.Begin()
	rows := []gen.DataFields{}
	fmt.Println("txn: ", txn, " rows:", rows, " err:", err)
}
