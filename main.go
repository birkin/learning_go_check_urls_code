package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
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

func check_site(site Site, dbwriter_channel chan Site) {
	/* Checks site, stores data to updated-site, & writes updated-site to channel. */
	rlog.Debug(fmt.Sprintf("go routine started for site, ```%v```", site.name))

	/// check site
	var site_check_result string = "init"
	mini_start_time := time.Now()
	resp, err := http.Get(site.url)
	if err != nil {
		rlog.Info(fmt.Sprintf("error accessing site, `%v`; error, ```%v```", site.name, err))
		site_check_result = "url_not_accessible"
	} else {
		body_bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			rlog.Info(fmt.Sprintf("error reading response from site, `%v`; error, ```%v```", site.name, err))
			site_check_result = "unable_to_read_response"
		}
		text := string(body_bytes)
		if strings.Contains(text, site.text_expected) {
			site_check_result = "passed"
		} else {
			site_check_result = "text_not_found"
		}
	}

	/// update site-object
	site.pre_previous_checked_result = site.previous_checked_result
	site.previous_checked_result = site.recent_checked_result
	site.recent_checked_result = site_check_result
	site.recent_checked_time = time.Now()

	/// determine whether to send email
	// var bool_val bool = run_email_check(site)
	// rlog.Debug(fmt.Sprintf("bool_val, `%v`", bool_val))

	/// determine next time-check -- TODO

	/// store other info to site
	/* TODO, update site object with next time-check */
	mini_elapsed := time.Since(mini_start_time)
	site.custom_time_taken = mini_elapsed

	/// write info to channel for db save
	dbwriter_channel <- site
	rlog.Info(fmt.Sprintf("site-info after write to channel, ```%#v```", site))

} // end func check_site()

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
