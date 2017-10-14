package main

import (
	"fmt"
	"github.com/romana/rlog"
	"time"
)

func check_sites_with_goroutines(sites []Site) {
	/* Flow:
	   - creates channel,
	   - kicks off go-routines to run the web-checks,
	   - channel writes each check-result to db,
	   - closes channel.
	   Called by main() */

	/*
		TODO flow...
		- initialize db-writer channel, which will get a Site, not a Result
		- for each site:
			- start `check_site()` go-routine
		- have the channel write the results of each updated site to the db
		- for each updated site, start `check_email_need()` go-routine, which will:
			- get email_flag
			- if email_flag is `send_failure_email` or `send_success_email`:
				- send email
	*/

	rlog.Debug(fmt.Sprintf("starting check_sites"))
	main_start := time.Now()

	/// initialize channel
	dbwriter_channel := make(chan Site)

	/// start go routines
	for _, site_element := range sites {
		// rlog.Debug(fmt.Sprintf("here"))
		go check_site(site_element, dbwriter_channel)
	}

	// rlog.Info(fmt.Sprintf("len(dbwriter_channel), ```%v```", len(dbwriter_channel)))

	/// handle channel data
	var counter int = 0
	var channel_site_data Site
	for channel_site_data = range dbwriter_channel {
		counter++
		time.Sleep(50 * time.Millisecond)
		/// save site info to db -- TODO
		/// go routine for checking whether email should be sent, and sending it if necessary
		go run_email_check(channel_site_data)
		rlog.Info("just called run_email_check()")
		rlog.Info(fmt.Sprintf("channel-data, ```%#v```", channel_site_data))
		if counter == len(sites) {
			// rlog.Info("about to close channel")
			close(dbwriter_channel)
			rlog.Info("channel closed")
			break // shouldn't be needed
		}
	}
	main_elapsed := time.Since(main_start)
	rlog.Info(fmt.Sprintf("main_elapsed, ```%v```", main_elapsed))

} // end func check_sites_with_goroutines()
