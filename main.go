package main

import (
	"fmt"
	"time"

	"github.com/romana/rlog"
)

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
	check_frequency_number      int
	check_frequency_unit        string
	custom_time_taken           time.Duration
}

var sites []Site // i think this declares a slice, not an array

func main() {
	/* Loads settings, initializes sites array, calls worker function. */

	/// start-tracking
	rlog.Info("\n\nstarting")
	strt_tm := time.Now()
	strt_str := fmt.Sprintf("%v", strt_tm.Format("2006-01-02 15:04:05"))
	rlog.Debug(fmt.Sprintf("main() start time, ```%v```", strt_str))

	/// initialize settings
	settings := load_settings() // settings.go
	rlog.Debug(fmt.Sprintf("settings, ```%#v```", settings))

	/// initialize sites
	// sites := initialize_sites_from_db(db) // db.go
	sites := initialize_sites_from_db(settings.DB_USERNAME, settings.DB_PASSWORD, settings.DB_HOST, settings.DB_PORT, settings.DB_NAME, settings.DB_TABLE) // db.go
	rlog.Debug("sites from db initialized")
	defer db.Close()

	/// TEMP email hijack
	for idx, _ := range sites {
		sites[idx].email_addresses = "birkin_diana@brown.edu, birkin.diana@gmail.com"
	}
	rlog.Info(fmt.Sprintf("updated sites, ```%#v```", sites)) // prints, eg, `{name:"clusters api", url:"etc...`

	/// call worker function
	check_sites_with_goroutines(sites) // check.go

	/// end-tracking
	main_elapsed := time.Since(strt_tm)
	rlog.Info(fmt.Sprintf("main() elapsed time, ```%v```", main_elapsed))
	rlog.Debug("end of main()")

} // end func main()

/// EOF
