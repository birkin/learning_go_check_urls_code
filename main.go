package main

import (
	"fmt"
	"time"

	"github.com/romana/rlog"
)

/*
TODO Next:
- √ temporarily put email-addresses string into a setting and use it.
- √ log sql update querystring
- √ save to db
- √ refactor initial db call
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
	calculated_seconds          int
	next_check_time             time.Time
	custom_time_taken           time.Duration
}

var sites []Site // i think this declares a slice, not an array
// var db *sql.DB
var now_string string

// var send_email bool

func main() {
	/* Loads settings, initializes sites array, calls worker function. */

	rlog.Info("\n\nstarting")

	/// initialize settings
	settings := load_settings() // settings.go
	rlog.Debug(fmt.Sprintf("settings, ```%#v```", settings))

	/// access db
	// db = setup_db(settings.DB_USERNAME, settings.DB_PASSWORD, settings.DB_HOST, settings.DB_PORT, settings.DB_NAME) // db.go

	/// prepare current-time
	t := time.Now()
	now_string = fmt.Sprintf("%v", t.Format("2006-01-02 15:04:05"))
	rlog.Debug(fmt.Sprintf("now_string, ```%v```", now_string))

	/// initialize sites
	// sites := initialize_sites_from_db(db) // db.go
	sites := initialize_sites_from_db(settings.DB_USERNAME, settings.DB_PASSWORD, settings.DB_HOST, settings.DB_PORT, settings.DB_NAME) // db.go
	rlog.Debug("sites from db initialized")
	defer db.Close()

	/// TEMP email hijack
	for idx, _ := range sites {
		sites[idx].email_addresses = "birkin_diana@brown.edu, birkin.diana@gmail.com"
	}
	rlog.Info(fmt.Sprintf("updated sites, ```%#v```", sites)) // prints, eg, `{name:"clusters api", url:"etc...`

	/// call worker function
	check_sites_with_goroutines(sites) // check.go

} // end func main()

/// EOF
