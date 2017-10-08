package pgx_conn

import (
	"fmt"
	"log"
	"os"
	"time"

	pgx "github.com/jackc/pgx"
)

func makeConfig(dbname string) pgx.ConnConfig {
	return pgx.ConnConfig{
		Host:      "localhost",
		Database:  dbname,
		User:      "swarm64",
		Password:  "swarm64",
		TLSConfig: nil,
		// RuntimeParams: map[string]string,
		OnNotice: noticeHandler}
}

func makeConnectionPool(config pgx.ConnConfig, n int) pgx.ConnPoolConfig {
	return pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: n,
		AfterConnect:   afterConnectHandler,
		AcquireTimeout: 0}
}

func afterConnectHandler(conn *pgx.Conn) error {
	var logThis string
	var cod error
	f, err := os.OpenFile(
		"pgx_after_connection.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)
	if err != nil {
		log.Fatal(err)
	}
	if conn.IsAlive() {
		logThis = "Connection established\n"
	} else {
		cod = conn.CauseOfDeath()
		logThis = fmt.Sprintf("[%v]\n\tcause of death: %v\n", cod.Error())
	}
	if _, err := f.Write([]byte(logThis)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	return cod
}

func noticeHandler(conn *pgx.Conn, notice *pgx.Notice) {
	var logThis string
	f, err := os.OpenFile(
		"pgx_notices.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)
	if err != nil {
		logThis = fmt.Sprintf("[%v]\n\terr: %v\n", time.Now(), err.Error())
	} else {
		logThis = fmt.Sprintf(
			"[%v]\n\t%v: %v, \n\tdetail: %v\n%v\n",
			time.Now(), notice.Severity, notice.Code, notice.Detail, notice.Hint)
	}
	if _, err := f.Write([]byte(logThis)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
