package main

/* all credit: <https://hackernoon.com/golang-sendmail-sending-mail-through-net-smtp-package-5cadbe2670e0> */

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/romana/rlog"
)

var settings Settings = load_settings() // settings.go

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

// func main_send() {
func main_send(site Site) {
	mail := Mail{}
	mail.senderId = "from_user@x.com"
	mail.toIds = []string{"to_user_a@y.com", "to_user_b@z.com"}
	mail.subject = "This is the email subject"
	mail.body = "This is the\n\nemail body."

	messageBody := mail.BuildMessage()

	// smtpServer := SmtpServer{host: "smtp.something.com", port: "the_port"}
	smtpServer := SmtpServer{host: settings.MAIL_HOST, port: "the_port"}

	// log.Println(smtpServer.host)
	rlog.Debug(fmt.Sprintf("smtpServer.host, `%v`", smtpServer.host))
	//build an auth
	auth := smtp.PlainAuth("", mail.senderId, "password", smtpServer.host)

	// Gmail will reject connection if it's not secure
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName: smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if err != nil {
		// log.Panic(err)
		rlog.Error(fmt.Sprintf("error on tls.Dial(), ```%v```", err))
		panic(err)
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		// log.Panic(err)
		rlog.Error(fmt.Sprintf("error on smtp.NewClient(), ```%v```", err))
		panic(err)
	}

	// step 1: Use Auth
	if err = client.Auth(auth); err != nil {
		// log.Panic(err)
		rlog.Error(fmt.Sprintf("error on client.Auth(), ```%v```", err))
		panic(err)
	}

	// step 2: add all from and to
	if err = client.Mail(mail.senderId); err != nil {
		log.Panic(err)
	}
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			// log.Panic(err)
			rlog.Error(fmt.Sprintf("error iterating through mail.toIds, ```%v```", err))
			panic(err)
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		// log.Panic(err)
		rlog.Error(fmt.Sprintf("error accessing client.Data(), ```%v```", err))
		panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		// log.Panic(err)
		rlog.Error(fmt.Sprintf("error on w.Write(), ```%v```", err))
		panic(err)
	}

	err = w.Close()
	if err != nil {
		// log.Panic(err)
		rlog.Error(fmt.Sprintf("error on w.Close(), ```%v```", err))
		panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")

}
