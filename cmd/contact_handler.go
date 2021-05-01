package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"
)

type confirmFn func(remoteip, response string) (result bool, err error)
type sendFn func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

func buildHandleContactFormFn(sendFn sendFn, confirmFn confirmFn, env AppEnv) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		contactRequest, err := NewContactRequest(request)
		if err != nil {
			fmt.Fprintf(writer, "false")
			return
		}
		result, err := confirmFn(GetIP(request), contactRequest.Recaptcha)
		if err != nil {
			fmt.Fprintf(writer, "false")
			return
		}
		if result {
			// toList is list of email address that email is to be sent.
			toList, err := toList(env.toEmail)

			if err != nil {
				fmt.Printf("Bad email to address %s", err.Error())
				fmt.Fprintf(writer, "false")
				return
			}
			m := mailArgs{
				Addr: env.emailHost + ":" + env.emailPort,
				Auth: smtp.PlainAuth("", env.emailUsername, env.emailPassword, env.emailHost),
				From: contactRequest.Email,
				To:   toList,
				Msg:  buildBody(contactRequest.Subject, contactRequest.Message),
			}
			sendFn(m.Addr, m.Auth, m.From, m.To, m.Msg)
		} else {
			time.Sleep(time.Duration(rand.Intn(8)) * time.Second)
		}
	}
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
