package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/romana/rlog"
)

func check_sites_with_goroutines(sites []Site) {
	/* Flow:
	   - create channel,
	   - kick off go-routines; each will...
	   		- run web-check
	   		- ascertain next-check-time
	   - channel:
	   		- writes each result & next-check-time to db
	   		- starts `check_email_need()` go-routine. possibilities: do-nothing, send_failure_email, send_success_email
	   - channel closes
	   Called by main() */

	rlog.Debug(fmt.Sprintf("starting check_sites"))
	main_start := time.Now()

	/// initialize channel
	dbwriter_channel := make(chan Site)

	/// start go routines
	for _, site_element := range sites {
		// rlog.Debug(fmt.Sprintf("here"))
		go check_site(site_element, dbwriter_channel)
	}

	/// handle channel data
	var counter int = 0
	var channel_site_data Site
	for channel_site_data = range dbwriter_channel {
		counter++
		time.Sleep(50 * time.Millisecond)
		/// save check-result to db -- TODO
		save_check_result(channel_site_data) // db.go

		/// go routine for checking whether email should be sent, and sending it if necessary
		go run_email_check(channel_site_data)
		rlog.Info("just called run_email_check()")
		rlog.Info(fmt.Sprintf("channel-data, ```%#v```", channel_site_data))
		if counter == len(sites) {
			close(dbwriter_channel)
			rlog.Info("channel closed")
			break // shouldn't be needed
		}
	}
	main_elapsed := time.Since(main_start)
	rlog.Info(fmt.Sprintf("main_elapsed, ```%v```", main_elapsed))

} // end func check_sites_with_goroutines()

func check_site(site Site, dbwriter_channel chan Site) {
	/*	Checks site, calculates next-check-time, updates site-object, & writes updated-site to channel.
		Called as go-routine by check_sites_with_goroutines()  */
	rlog.Debug(fmt.Sprintf("go routine started for site, ```%v```", site.name))

	/// check site
	var site_check_result string = "site_check_result_init"
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

	/// calculate next-check-time
	original_next_check_time := site.next_check_time
	rlog.Debug(fmt.Sprintf("original_next_check_time, ```%v```", original_next_check_time))

	/// update site-object
	site.pre_previous_checked_result = site.previous_checked_result
	site.previous_checked_result = site.recent_checked_result
	site.recent_checked_result = site_check_result
	site.recent_checked_time = time.Now()
	site.next_check_time = calc_next_check_time(site)
	rlog.Debug(fmt.Sprintf("calculated next_check_time, ```%v```", site.next_check_time))

	/// determine whether to send email
	// var bool_val bool = run_email_check(site)
	// rlog.Debug(fmt.Sprintf("bool_val, `%v`", bool_val))

	/// store other info to site
	/* TODO, update site object with next time-check */
	mini_elapsed := time.Since(mini_start_time)
	site.custom_time_taken = mini_elapsed

	/// write info to channel for db save
	dbwriter_channel <- site
	rlog.Info(fmt.Sprintf("site-info after write to channel, ```%#v```", site))

} // end func check_site()

func calc_next_check_time(site Site) time.Time {
	/*	Calculates next check time.
		Called by check_site()  */
	rlog.Debug(fmt.Sprintf("original site.calculated_seconds, ```%v```", site.calculated_seconds))
	t := time.Now()
	rlog.Debug(fmt.Sprintf("now-time, ```%v```", t))
	duration := time.Second * time.Duration(site.calculated_seconds)
	rlog.Debug(fmt.Sprintf("duration, ```%v```", duration))
	next_check_time := t.Add(duration)
	rlog.Debug(fmt.Sprintf("next_check_time, ```%v```", next_check_time))
	return next_check_time
}

func run_email_check(site Site) bool {
	/*	Determines whether email should be sent.
		Called as go-routine by check_sites_with_goroutines()  */
	rlog.Debug("checking whether to send email")
	var bool_val bool = false
	rand.Seed(time.Now().UnixNano()) // initialize global pseudo random generator
	num := rand.Intn(2)              // so will be 0 or 1
	rlog.Info(fmt.Sprintf("num, `%v`", num))
	if num == 1 {
		bool_val = true
	}
	rlog.Info(fmt.Sprintf("bool_val, `%v`", bool_val))
	return bool_val
}
