package pgx_config

import (
  "database/sql"
	"fmt"
	"os"

	pgx "github.com/jackc/pgx"
)
// add a logger

const ConnectionConfig = pgx.ConnConfig {
    Host              "localhost"
    Port              5432
    Database          "swarmtest"
    User              "maruko"
    Logger            Logger
    LogLevel          int
}
