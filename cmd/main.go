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
	"net/http"
	"net/mail"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/wjase/recontact-mailer/internal/recontact"
)

var logger = log.Default()

func init() {

}

func main() {
	appEnv := recontact.NewAppEnv()

	logger.Println("Starting recontact-mailer server...")
	recaptcha.Init(appEnv.PrivateKey)

	// send happy email
	err := recontact.SendMail(fmt.Sprintf("%s:%s", appEnv.EmailHost, appEnv.EmailPort), (&mail.Address{Name: "App Admin", Address: appEnv.AdminEmail}).String(), "Email Subject", "Recapture started successfully", []string{(&mail.Address{Name: "admin", Address: appEnv.AdminEmail}).String()})
	if err != nil {
		logger.Fatal("failed to send mail", err)
	}

	fmt.Printf("endpoint is %s", appEnv.Endpoint)

	http.HandleFunc(appEnv.Endpoint, recontact.BuildHandleContactFormFn(recontact.SendMail, recaptcha.Confirm, appEnv))
	http.HandleFunc("/ping", pingHandler)

	logger.Printf("About To start server on port %s\n", appEnv.AppPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", appEnv.AppPort), nil); err != nil {
		logger.Fatal("failed to start server", err)
	}
}

func pingHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("ping ok"))
}
