package main

import (
	"fmt"

	"github.com/romana/rlog"
)

func run_email_check(site Site) {
	/*	Determines whether email should be sent.
		Called as go-routine by check_sites_with_goroutines()  */
	rlog.Debug("starting run_email_check()")
	send, type_send := assess_email_need(site)
	if send == true {
		send_email(site, type_send)
	}
	rlog.Info("end of run_email_check()")
}

func assess_email_need(site Site) (bool, string) {
	/*	Determines if email needs to be sent, and, if so,whether it should be a `service-back-up` or a `service-down` email.
		Called by run_email_check()
		- send_no_email logic:
			- site.recent_checked_result == "passed" && site.previous_checked_result == "passed"  // all's well
			- site.recent_checked_result != "passed" && site.previous_checked_result == "passed" && site.pre_previous_checked_result == "passed"  // possible temporary failure
			- site.recent_checked_result == "passed" && site.previous_checked_result != "passed" && site.pre_previous_checked_result == "passed"  // recovery from temporary temporary failure
			- site.recent_checked_result != "passed" && site.previous_checked_result != "passed" && site.pre_previous_checked_result != "passed"  // repeated failure
			- site.previous_checked_result == "" // new entry  */
	rlog.Debug("starting assess_email_need()")
	var send bool = false
	var send_type string = "send_no_email"
	/// failure email
	if site.recent_checked_result != "passed" &&
		site.previous_checked_result != "passed" &&
		site.pre_previous_checked_result == "passed" {
		send = true
		send_type = "send_failure_email"
	}
	/// success email
	if site.recent_checked_result == "passed" &&
		site.previous_checked_result != "passed" &&
		site.pre_previous_checked_result != "passed" {
		send = true
		send_type = "send_success_email"
	}

	/// TEMP FOR DEVELOPMENT
	send = true
	send_type = "send_failure_email"
	/// END TEMP FOR DEVELOPMENT

	rlog.Debug(fmt.Sprintf("send, `%v`; send_type, `%v`", send, send_type))
	return send, send_type
}

func send_email(site Site, type_send string) {
	/*  Sends email if called.
	    Called by run_email_check()  */
	if type_send == "send_success_email" {
		send_success_email(site) // email_send.go
	} else {
		send_failure_email(site) // email_send.go
	}
	rlog.Debug(fmt.Sprintf("`%v` email sent", type_send))
	return
}
