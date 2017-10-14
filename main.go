package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	// _ "github.com/go-sql-driver/mysql" // package is imported only for its `side-effects`; it gets registered as the driver for the regular database/sql package
	"github.com/romana/rlog"
)

/*
TODO Next:
- check python code for 'save()' work
	- replicate in go
	- above should set `next-check-time`
	- save to db
- go-routine email call should be to handle_email(), which should:
	- see if email needs to be sent
	- send it
*/

type Site struct {
	id                          int
	name                        string
	url                         string
	text_expected               string
	email_addresses             string
	email_message               string
	recent_checked_time         time.Time
	recent_checked_result       string
	previous_checked_result     string
	pre_previous_checked_result string
	next_check_time             time.Time
	custom_time_taken           time.Duration
}

var sites []Site // i think this declares a slice, not an array
var db *sql.DB
var now_string string
var send_email bool

func main() {
	/* Loads settings, initializes sites array, calls worker function. */

	rlog.Info("\n\nstarting")

	/// initialize settings
	settings := load_settings() // settings.go
	rlog.Debug(fmt.Sprintf("settings, ```%#v```", settings))

	/// access db
	db = setup_db(settings.DB_USERNAME, settings.DB_PASSWORD, settings.DB_HOST, settings.DB_PORT, settings.DB_NAME) // db.go

	/// prepare current-time
	t := time.Now()
	now_string = fmt.Sprintf("%v", t.Format("2006-01-02 15:04:05"))
	rlog.Debug(fmt.Sprintf("now_string, ```%v```", now_string))

	/// initialize sites
	sites := initialize_sites_from_db(db) // db.go
	rlog.Debug("sites from db initialized")
	defer db.Close()

	/// call worker function
	check_sites_with_goroutines(sites)

} // end func main()

/* ----------------------------------------------------------------------
   helper functions
   ---------------------------------------------------------------------- */

func run_email_check(site Site) bool {
	/* Determines whether email should be sent. */
	rlog.Debug("checking whether to send email")
	var bool_val bool = true
	rand.Seed(time.Now().UnixNano()) // initialize global pseudo random generator
	num := rand.Intn(2)              // so will be 0 or 1
	rlog.Info(fmt.Sprintf("num, `%v`", num))
	if num == 1 {
		bool_val = false
	}
	rlog.Info(fmt.Sprintf("bool_val, `%v`", bool_val))
	return bool_val
}

/// EOF
