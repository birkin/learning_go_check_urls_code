package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/romana/rlog"
)

var settings Settings = load_settings() // settings.go

func send_success_email(site Site) {
	rlog.Debug("success email will be sent here")
}

func send_failure_email(site Site) {
	/* Sends email.
	   Called by email_prep.go -> send_email() */
	/// host/port/auth stuff
	host_string := fmt.Sprintf("%v", settings.MAIL_HOST)
	rlog.Debug(fmt.Sprintf("host_string, ```%v```", host_string))
	host_port_string := fmt.Sprintf("%v:%v", settings.MAIL_HOST, settings.MAIL_PORT)
	rlog.Debug(fmt.Sprintf("host_port_string, ```%v```", host_port_string))
	password_string := ""                                                          // no password needed, but smtp.SendMail(), called below, requires the auth module
	auth := smtp.PlainAuth("", settings.MAIL_SENDER, password_string, host_string) // settings.MAIL_SENDER used for smtp.PlainAuth() and smtp.SendMail() commands
	/// sender stuff
	rlog.Debug(fmt.Sprintf("actual-sender from settings, ```%v```", settings.MAIL_SENDER))
	perceived_sender_string := "Brown Library automated site-checker"
	rlog.Debug(fmt.Sprintf("perceived_sender_string, ```%v```", perceived_sender_string))
	/// recipent stuff
	var recipent_entry string = site.email_addresses
	rlog.Debug(fmt.Sprintf("recipent_entry, ```%v```", recipent_entry))
	var recipients []string = strings.Split(recipent_entry, ",")
	rlog.Debug(fmt.Sprintf("recipients, ```%v```", recipients))
	/// body stuff
	var body string = make_failure_body(site)
	/// assemble pieces
	msg := []byte(
		fmt.Sprintf("To: %v\r\n", recipients) +
			fmt.Sprintf("From: %v\r\n", perceived_sender_string) +
			fmt.Sprintf("Subject: Service-Status alert: \"%v\" problem\r\n", site.name) +
			"\r\n" +
			body +
			"\r\n",
	)
	/// send
	err := smtp.SendMail(host_port_string, auth, settings.MAIL_SENDER, recipients, msg)
	if err != nil {
		log.Fatal(err)
	}
} // end func send_failure_email()

func make_failure_body(site Site) string {
	var frequency_unit string = site.check_frequency_unit
	if site.check_frequency_number > 1 {
		frequency_unit = frequency_unit + "s"
	}
	var body string = ""
	body += fmt.Sprintf("The service: \"%v\" appears to be down.\r\n", site.name)
	body += "\r\n"
	body += fmt.Sprintf(
		"The \"%v\" service failed two consecutive automated checks a few minutes apart. Checks will continue every few minutes while the failures persist, but you will only be emailed again when the automated check succeeds. Once the automated check succeeds, the check-frequency will return to the specified value of every-%v-%v.\r\n",
		site.name, site.check_frequency_number, frequency_unit)
	body += "\r\n"
	body += fmt.Sprintf(
		"- Url checked: \"%v\"\r\n",
		site.url)
	body += fmt.Sprintf(
		"- Text expected: \"%v\"\r\n",
		site.text_expected)
	body += fmt.Sprintf(
		"- Specified failure message: \"%v\"\r\n",
		site.email_message)
	body += "\r\n"
	body += "You can view the current status of all services set up for automated checking at:\r\n"
	body += "<http://library.brown.edu/services/site_checker/status/>\r\n"
	body += "\r\n"

	rlog.Debug(fmt.Sprintf("body, ```%v```", body))
	return body
}
