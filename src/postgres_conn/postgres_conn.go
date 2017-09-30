package postgres_conn

import (
	"database/sql"
	"fmt"
	"os"

	pgx "github.com/jackc/pgx"
	pq "github.com/lib/pq"
)

func MakeConnection() *sql.DB {
	conn, err := sql.Open("postgres", "user=maruko dbname=swarmtest")
	if err != nil {
		fmt.Println("Could not connect to database")
		os.Exit(1)
	}
	return conn
}

// just to keep the libs imported
func tmp_pq(db *sql.DB) {
	txn, err := db.Begin()
	txn.Prepare(pq.CopyIn("ingestion_test", "src", "dst", "flags"))
}

func tmp_pgx(db *sql.DB) {
	txn, err := db.Begin()
	rows := []DataFields{}

	copyCount, err := conn.CopyFrom(
		pgx.Identifier{"ingestion_test"},
		[]string{"src", "dst", "flags"},
		pgx.CopyFromRows(rows),
	)
}
