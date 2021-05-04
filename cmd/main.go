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
	"os"
	"time"

	"github.com/dpapathanasiou/go-recaptcha"
)

var logger = log.Default()

type AppEnv struct {
	privateKey string
	toEmail    string
	adminEmail string
	emailHost  string
	emailPort  string
	appPort    string
	endpoint   string
}

func ensureEnvNotBlank(name string) string {
	if val, ok := os.LookupEnv(name); !ok {
		logger.Printf("Unexpected blank property %s\n", name)
		panic(fmt.Sprintf("Unexpected blank property %s\n", name))
	} else {
		return val
	}
}

// NewAppEnv cerates a new env.
func NewAppEnv() AppEnv {

	return AppEnv{
		privateKey: ensureEnvNotBlank("RECAPTCHA_PRIVATE_KEY"),
		toEmail:    ensureEnvNotBlank("TO_MAIL"),
		adminEmail: ensureEnvNotBlank("ADMIN_MAIL"),
		emailHost:  ensureEnvNotBlank("EMAIL_HOST"),
		emailPort:  ensureEnvNotBlank("EMAIL_PORT"),
		appPort:    ensureEnvNotBlank("APP_PORT"),
		endpoint:   ensureEnvNotBlank("ENDPOINT"),
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	appEnv := NewAppEnv()

	logger.Println("Starting recontact-mailer...")
	recaptcha.Init(appEnv.privateKey)

	// send happy email
	err := SendMail("127.0.0.1:25", (&mail.Address{Name: "App Admin", Address: appEnv.adminEmail}).String(), "Email Subject", "Recapture started successfully", []string{(&mail.Address{Name: "admin", Address: appEnv.adminEmail}).String()})
	if err != nil {
		logger.Fatal("failed to send mail", err)
	}

	http.HandleFunc("/contactform", buildHandleContactFormFn(SendMail, recaptcha.Confirm, appEnv))
	logger.Println("About To start server")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", appEnv.appPort), nil); err != nil {
		logger.Fatal("failed to start server", err)
	}
}
