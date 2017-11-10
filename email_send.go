package main

/* all credit: <https://hackernoon.com/golang-sendmail-sending-mail-through-net-smtp-package-5cadbe2670e0> */

import (
	// "bytes"
	// "crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	// "strings"

	"github.com/romana/rlog"
)

var settings Settings = load_settings() // settings.go

func main_send(site Site) {

	recipient_string := settings.TEST_MAIL_RECIPIENT
	rlog.Debug(fmt.Sprintf("recipient_string, ```%v```", recipient_string))

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
	auth := smtp.PlainAuth("", sender_string, password_string, host_string)

	to := []string{recipient_string}

	msg := []byte(
		fmt.Sprintf(
			"To: %v\r\n", recipient_string) +
			"From: automated-site-checker\r\n" +
			"Subject: discount Gophers!\r\n" +
			"\r\n" +
			"This is the email body test.\r\n",
	)

	rlog.Debug(fmt.Sprintf("msg, ```%v```", msg))

	err := smtp.SendMail(host_port_string, auth, sender_string, to, msg)
	if err != nil {
		log.Fatal(err)
	}

}
