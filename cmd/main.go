// example.go
//
// A simple HTTP server which presents a reCaptcha input form and evaulates the result,
// using the github.com/dpapathanasiou/go-recaptcha package.
//
// See the main() function for usage.
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/mail"
	"time"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/wjase/recontact-mailer/internal/recontact"
)

var logger = log.Default()

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	appEnv := recontact.NewAppEnv()

	logger.Println("Starting recontact-mailer...")
	recaptcha.Init(appEnv.PrivateKey)

	// send happy email
	err := recontact.SendMail("127.0.0.1:25", (&mail.Address{Name: "App Admin", Address: appEnv.AdminEmail}).String(), "Email Subject", "Recapture started successfully", []string{(&mail.Address{Name: "admin", Address: appEnv.AdminEmail}).String()})
	if err != nil {
		logger.Fatal("failed to send mail", err)
	}

	http.HandleFunc("/contactform", recontact.BuildHandleContactFormFn(recontact.SendMail, recaptcha.Confirm, appEnv))
	logger.Printf("About To start server on port %s\n", appEnv.AppPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", appEnv.AppPort), nil); err != nil {
		logger.Fatal("failed to start server", err)
	}
}
