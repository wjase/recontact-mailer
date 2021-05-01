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
	"net/smtp"
	"os"
	"time"

	"github.com/dpapathanasiou/go-recaptcha"
)

type AppEnv struct {
	private_key   string
	toEmail       string
	emailUsername string
	emailPassword string
	emailHost     string
	emailPort     string
	appPort       string
}

func NewAppEnv() AppEnv {
	return AppEnv{
		private_key:   os.Getenv("RECAPTCHA_PRIVATE_KEY"),
		toEmail:       os.Getenv("TO_MAIL"),
		emailUsername: os.Getenv("EMAIL_USERNAME"),
		emailPassword: os.Getenv("EMAIL_PASSWORD"),
		emailHost:     os.Getenv("EMAIL_HOST"),
		emailPort:     os.Getenv("EMAIL_PORT"),
		appPort:       os.Getenv("APP_PORT"),
	}
}

var appEnv = NewAppEnv()

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	recaptcha.Init(appEnv.private_key)

	http.HandleFunc("/contactform", buildHandleContactFormFn(smtp.SendMail, recaptcha.Confirm, appEnv))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", appEnv.appPort), nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}
