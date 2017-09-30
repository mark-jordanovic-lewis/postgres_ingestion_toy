package main

import (
	"database/sql"
	"os"
)

//  "github.com/lib/pq"

func main() {
	db, err := sql.Open(
		"postgres", "user=%s dbname=%s password=%s",
		os.Arg[1], os.Arg[2], os.Arg[3])
	txn, err := db.Begin()
}
