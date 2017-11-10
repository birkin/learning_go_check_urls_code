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

	sender_string := settings.MAIL_SENDER
	rlog.Debug(fmt.Sprintf("sender_string, ```%v```", sender_string))

	password_string := ""

	/// Set up authentication information.
	// auth := smtp.PlainAuth("", sender_string, password_string, "mail.example.com")
	auth := smtp.PlainAuth("", sender_string, password_string, host_string)

	/// Connect to the server, authenticate, set the sender and recipient,
	/// and send the email all in one step.

	to := []string{recipient_string}

	// msg := []byte(
	// 	fmt.Sprintf("To: %v\r\n", recipient_string) +
	// 		fmt.Sprintf("From: %v\r\n", sender_string) +
	// 		"Subject: discount Gophers!\r\n" +
	// 		"\r\n" +
	// 		"This is the email body test.\r\n")

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

/// ---------- ///

// func main_send_b(site Site) {
// 	/// Connect to the remote SMTP server.
// 	dial_string := fmt.Sprintf("%v:%v", settings.MAIL_HOST, settings.MAIL_PORT)
// 	rlog.Debug(fmt.Sprintf("dial_string, ```%v```", dial_string))
// 	c, err := smtp.Dial(dial_string)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer c.Close()
// 	/// Set the sender and recipient.
// 	// rlog.Debug("sender, %v", settings.MAIL_SENDER)
// 	rlog.Debug(fmt.Sprintf("sender, %v", settings.MAIL_SENDER))
// 	c.Mail(settings.MAIL_SENDER)
// 	c.Rcpt(settings.TEST_MAIL_RECIPIENT)
// 	/// Send the email body.
// 	wc, err := c.Data()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer wc.Close()
// 	buf := bytes.NewBufferString("Email body -- test email-send using golang.")
// 	if _, err = buf.WriteTo(wc); err != nil {
// 		log.Fatal(err)
// 	}
// }

/// ---------- ///

// type Mail struct {
// 	senderId string
// 	toIds    []string
// 	subject  string
// 	body     string
// }

// type SmtpServer struct {
// 	host string
// 	port string
// }

// func (s *SmtpServer) ServerName() string {
// 	return s.host + ":" + s.port
// }

// func (mail *Mail) BuildMessage() string {
// 	message := ""
// 	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
// 	if len(mail.toIds) > 0 {
// 		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
// 	}

// 	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
// 	message += "\r\n" + mail.body

// 	return message
// }

// // func main_send() {
// func main_send(site Site) {
// 	mail := Mail{}
// 	mail.senderId = settings.MAIL_SENDER
// 	mail.toIds = []string{"to_user_a@y.com", "to_user_b@z.com"}
// 	mail.subject = "This is the email subject"
// 	mail.body = "This is the\n\nemail body."

// 	messageBody := mail.BuildMessage()

// 	// smtpServer := SmtpServer{host: "smtp.something.com", port: "the_port"}
// 	smtpServer := SmtpServer{host: settings.MAIL_HOST, port: settings.MAIL_PORT}
// 	rlog.Debug(fmt.Sprintf("smtpServer.host, `%v`", smtpServer.host))
// 	rlog.Debug(fmt.Sprintf("smtpServer.port, `%v`", smtpServer.port))

// 	//build an auth
// 	// auth := smtp.PlainAuth("", mail.senderId, "password", smtpServer.host)
// 	auth := smtp.PlainAuth("", mail.senderId, "", smtpServer.host)
// 	rlog.Debug(fmt.Sprintf("auth, `%v`", auth))

// 	// Gmail will reject connection if it's not secure
// 	// TLS config
// 	// tlsconfig := &tls.Config{
// 	// 	InsecureSkipVerify: true,
// 	// 	ServerName:         smtpServer.host,
// 	// }
// 	tlsconfig := &tls.Config{
// 		InsecureSkipVerify: false,
// 		ServerName:         smtpServer.host,
// 	}
// 	rlog.Debug(fmt.Sprintf("tlsconfig, `%v`", tlsconfig))

// 	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
// 	if err != nil {
// 		// log.Panic(err)
// 		rlog.Error(fmt.Sprintf("error on tls.Dial(), ```%v```", err))
// 		panic(err)
// 	}

// 	client, err := smtp.NewClient(conn, smtpServer.host)
// 	if err != nil {
// 		// log.Panic(err)
// 		rlog.Error(fmt.Sprintf("error on smtp.NewClient(), ```%v```", err))
// 		panic(err)
// 	}

// 	// step 1: Use Auth
// 	if err = client.Auth(auth); err != nil {
// 		// log.Panic(err)
// 		rlog.Error(fmt.Sprintf("error on client.Auth(), ```%v```", err))
// 		panic(err)
// 	}

// 	// step 2: add all from and to
// 	if err = client.Mail(mail.senderId); err != nil {
// 		log.Panic(err)
// 	}
// 	for _, k := range mail.toIds {
// 		if err = client.Rcpt(k); err != nil {
// 			// log.Panic(err)
// 			rlog.Error(fmt.Sprintf("error iterating through mail.toIds, ```%v```", err))
// 			panic(err)
// 		}
// 	}

// 	// Data
// 	w, err := client.Data()
// 	if err != nil {
// 		// log.Panic(err)
// 		rlog.Error(fmt.Sprintf("error accessing client.Data(), ```%v```", err))
// 		panic(err)
// 	}

// 	_, err = w.Write([]byte(messageBody))
// 	if err != nil {
// 		// log.Panic(err)
// 		rlog.Error(fmt.Sprintf("error on w.Write(), ```%v```", err))
// 		panic(err)
// 	}

// 	err = w.Close()
// 	if err != nil {
// 		// log.Panic(err)
// 		rlog.Error(fmt.Sprintf("error on w.Close(), ```%v```", err))
// 		panic(err)
// 	}

// 	client.Quit()

// 	log.Println("Mail sent successfully")

// }
