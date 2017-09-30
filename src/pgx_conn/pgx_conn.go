package pgx_conn

import (
	"database/sql"

	pgx "github.com/jackc/pgx"
)

func tmp_pgx(db *sql.DB) {
	txn, err := db.Begin()
	rows := []DataFields{}

	copyCount, err := conn.CopyFrom(
		pgx.Identifier{"ingestion_test"},
		[]string{"src", "dst", "flags"},
		pgx.CopyFromRows(rows),
	)
}
