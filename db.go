package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // package is imported only for its `side-effects`; it gets registered as the driver for the regular database/sql package
	"github.com/romana/rlog"
)

func setup_db(user string, pass string, host string, port string, name string) *sql.DB {
	/* Initializes db object and confirms connection.
	   Called by main() */
	var connect_str string = fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?parseTime=true",
		user, pass, host, port, name) // user:password@tcp(host:port)/dbname
	db, err := sql.Open("mysql", connect_str)
	if err != nil {
		msg := fmt.Sprintf("error connecting to db, ```%v```", err)
		rlog.Error(msg)
		panic(msg)
	}
	/// sql.Open doesn't open a connection, so validate DSN (data source name) data
	err = db.Ping()
	if err != nil {
		msg := fmt.Sprintf("error accessing db, ```%v```", err)
		rlog.Error(msg)
		panic(msg)
	}
	return db
} // end func setup_db()
