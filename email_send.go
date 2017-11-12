package main

import (
	"fmt"
	"github.com/romana/rlog"
	"log"
	"net/smtp"
	"reflect"
	"strings"
)

var settings Settings = load_settings() // settings.go

func main_send(site Site) {

	recipient_string := settings.TEST_MAIL_RECIPIENT
	rlog.Debug(fmt.Sprintf("recipient_string, ```%v```", recipient_string))

	////////////////

	/// split testing
	somevar := strings.Split(recipient_string, ",")
	lg_msgA := fmt.Sprintf("`somevar` has TypeOf, ```%v```", reflect.TypeOf(somevar))
	rlog.Debug(lg_msgA)
	somevar_kind := reflect.ValueOf(somevar)
	lg_msgA = fmt.Sprintf("`somevar_kind` has Kind, ```%v```", somevar_kind.Kind())
	rlog.Debug(lg_msgA)
	rlog.Debug(fmt.Sprintf("somevar, ```%v```", somevar))

	/////////////////////

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

	to := []string{recipient_string}
	lg_msg := fmt.Sprintf("`to` has TypeOf, ```%v```", reflect.TypeOf(to))
	rlog.Debug(lg_msg)
	to_kind := reflect.ValueOf(to)
	lg_msg = fmt.Sprintf("`to` has Kind, ```%v```", to_kind.Kind())
	rlog.Debug(lg_msg)

	msg := []byte(
		fmt.Sprintf("To: %v\r\n", recipient_string) +
			fmt.Sprintf("From: %v\r\n", perceived_sender_string) +
			"Subject: discount Gophers!\r\n" +
			"\r\n" +
			"This is the email body test.\r\n",
	)

	rlog.Debug(fmt.Sprintf("msg, ```%v```", msg))

	err := smtp.SendMail(host_port_string, auth, actual_sender_string, to, msg)
	if err != nil {
		log.Fatal(err)
	}

}
