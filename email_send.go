package main

import (
	"fmt"
	"github.com/romana/rlog"
	"log"
	"net/smtp"
	"strings"
)

var settings Settings = load_settings() // settings.go

func send_success_email(site Site) {
	rlog.Debug("success email will be sent here")
}

func send_failure_email(site Site) {

	// var recipent_entry string = settings.TEST_MAIL_RECIPIENT
	// rlog.Debug(fmt.Sprintf("recipent_entry, ```%v```", recipent_entry))

	var recipent_entry string = site.email_addresses
	rlog.Debug(fmt.Sprintf("recipent_entry, ```%v```", recipent_entry))

	var recipients []string = strings.Split(recipent_entry, ",")
	rlog.Debug(fmt.Sprintf("recipients, ```%v```", recipients))

	host_string := fmt.Sprintf("%v", settings.MAIL_HOST)
	rlog.Debug(fmt.Sprintf("host_string, ```%v```", host_string))

	host_port_string := fmt.Sprintf("%v:%v", settings.MAIL_HOST, settings.MAIL_PORT)
	rlog.Debug(fmt.Sprintf("host_port_string, ```%v```", host_port_string))

	actual_sender_string := settings.MAIL_SENDER
	rlog.Debug(fmt.Sprintf("actual_sender_string, ```%v```", actual_sender_string))

	perceived_sender_string := "Brown Library automated site-checker"
	rlog.Debug(fmt.Sprintf("perceived_sender_string, ```%v```", perceived_sender_string))

	/// no password needed, but smtp.SendMail requires the auth module
	password_string := ""
	auth := smtp.PlainAuth("", actual_sender_string, password_string, host_string)

	// msg := []byte(
	// 	fmt.Sprintf("To: %v\r\n", recipients) +
	// 		fmt.Sprintf("From: %v\r\n", perceived_sender_string) +
	// 		fmt.Sprintf("Subject: %v\r\n", site.name) +
	// 		"\r\n" +
	// 		"This is the email body test.\r\n",
	// )

	var body string = make_failure_body(site)
	msg := []byte(
		fmt.Sprintf("To: %v\r\n", recipients) +
			fmt.Sprintf("From: %v\r\n", perceived_sender_string) +
			// fmt.Sprintf("Subject: %v\r\n", site.name) +
			fmt.Sprintf("Subject: Service-Status alert: \"%v\" problem\r\n", site.name) +
			"\r\n" +
			// "This is the email body test.\r\n",
			body +
			"\r\n",
	)

	// rlog.Debug(fmt.Sprintf("msg, ```%v```", msg))

	err := smtp.SendMail(host_port_string, auth, actual_sender_string, recipients, msg)
	if err != nil {
		log.Fatal(err)
	}
	make_failure_body(site)

} // end func send_failure_email()

func make_failure_body(site Site) string {
	var frequency_unit string = site.check_frequency_unit
	if site.check_frequency_number > 1 {
		frequency_unit = frequency_unit + "s"
	}
	var body string = ""
	body += fmt.Sprintf("The service: \"%v\"" appears to be down.", site.name)
	body += "\r\n"
	body += "\r\n"
	body += fmt.Sprintf(
		"The \"%v\" service failed two consecutive automated checks a few minutes apart. Checks will continue every few minutes while the failures persist, but you will only be emailed again when the automated check succeeds. Once the automated check succeeds, the check-frequency will return to the specified values of every-%v-%v.",
		site.name, site.check_frequency_number, frequency_unit)

	rlog.Debug(fmt.Sprintf("body, ```%v```", body))
	return body
}

//     message = '''The service "%s" appears to be down.

// The "%s" service failed two consecutive automated checks a few minutes apart. Checks will continue every few minutes while the failures persist, but you will only be emailed again when the automated check succeeds. Once the automated check succeeds, the check-frequency will return to the specified values of every-%s-%s(s).

// - Url checked: "%s"
// - Text expected: "%s"
// - Specified failure message: "%s"

// You can view the current status of all services set up for automated checking at:
// <http://library.brown.edu/services/site_checker/status/>

// If authorized, you can edit service automated checking at:
// <%s>

// [end]
// ''' % ( site.name,
//         site.name,
//         site.check_frequency_number,
//         site.check_frequency_unit,
//         site.url,
//         site.text_expected,
//         site.email_message,
//         settings_app.EMAIL_ADMIN_URL )
