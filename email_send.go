package main

import (
	"fmt"
	"github.com/romana/rlog"
	"log"
	"net/smtp"
	"strings"
)

var settings Settings = load_settings() // settings.go

func main_send(site Site) {

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
	// 		"Subject: discount Gophers!\r\n" +
	// 		"\r\n" +
	// 		"This is the email body test.\r\n",
	// )

	msg := []byte(
		fmt.Sprintf("To: %v\r\n", recipients) +
			fmt.Sprintf("From: %v\r\n", perceived_sender_string) +
			fmt.Sprintf("Subject: %v\r\n", site.name) +
			"\r\n" +
			"This is the email body test.\r\n",
	)

	rlog.Debug(fmt.Sprintf("msg, ```%v```", msg))

	err := smtp.SendMail(host_port_string, auth, actual_sender_string, recipients, msg)
	if err != nil {
		log.Fatal(err)
	}

}
