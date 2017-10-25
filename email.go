package main

// import (
// 	"github.com/romana/rlog"
// )

// func run_email_check(site Site) {
// 	/*	Determines whether email should be sent.
// 		Called as go-routine by check_sites_with_goroutines()  */
// 	send, type_send := assess_email_need(site)
// 	if send == true {
// 		send_email(site, type_send)
// 	}
// 	rlog.Info("end of run_email_check()")
// }

// func assess_email_need(site Site) (bool, string) {
// 	  Determines if email needs to be sent, and, if so,whether it should be a `service-back-up` or a `service-down` email.
// 		Called by run_email_check()
// 	send := true
// 	send_type := "send_success_email"
// 	rlog.Debug(fmt.Sprintf("send, `%v`; send_type, `%v`", send, send_type))
// 	return send, send_type
// }

// func send_email(site Site, type_send string) {
// 	  Sends email if called.
//         Called by run_email_check()
// 	rlog.Debug("email sent")
// 	return
// }
