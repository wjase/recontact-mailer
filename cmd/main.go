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

type AppEnv struct {
	privateKey    string
	toEmail       string
	adminEmail    string
	emailUsername string
	emailPassword string
	emailHost     string
	emailPort     string
	appPort       string
	endpoint      string
}

func ensureNotBlank(name string, s string) bool {
	if len(s) == 0 {
		fmt.Printf("Unexpected blank property %s\n", s)
		return false
	}
	return true
}

func (a AppEnv) validate() bool {
	return ensureNotBlank("privateKey", a.privateKey) &&
		ensureNotBlank("toEmail", a.toEmail) &&
		ensureNotBlank("adminEmail", a.adminEmail) &&
		ensureNotBlank("emailUsername", a.emailUsername) &&
		ensureNotBlank("emailPassword", a.emailPassword) &&
		ensureNotBlank("emailHost", a.emailHost) &&
		ensureNotBlank("emailPort", a.emailPort) &&
		ensureNotBlank("appPort", a.appPort) &&
		ensureNotBlank("endpoint", a.endpoint)
}

// NewAppEnv cerates a new env.
func NewAppEnv() AppEnv {
	return AppEnv{
		privateKey:    os.Getenv("RECAPTCHA_PRIVATE_KEY"),
		toEmail:       os.Getenv("TO_MAIL"),
		adminEmail:    os.Getenv("ADMIN_MAIL"),
		emailUsername: os.Getenv("EMAIL_USERNAME"),
		emailPassword: os.Getenv("EMAIL_PASSWORD"),
		emailHost:     os.Getenv("EMAIL_HOST"),
		emailPort:     os.Getenv("EMAIL_PORT"),
		appPort:       os.Getenv("APP_PORT"),
		endpoint:      os.Getenv("ENDPOINT"),
	}
}

var appEnv = NewAppEnv()

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	recaptcha.Init(appEnv.privateKey)

	if !appEnv.validate() {
		fmt.Println("Env validate failed. Stopping")
	}

	// send happy email
	SendMail("127.0.0.1:25", (&mail.Address{"App Admin", appEnv.adminEmail}).String(), "Email Subject", "Recapture started successfully", []string{(&mail.Address{Name: "to name", Address: appEnv.adminEmail}).String()})

	http.HandleFunc("/contactform", buildHandleContactFormFn(SendMail, recaptcha.Confirm, appEnv))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", appEnv.appPort), nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}
