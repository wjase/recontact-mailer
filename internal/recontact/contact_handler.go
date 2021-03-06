package recontact

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type confirmFn func(remoteip, response string) (result bool, err error)
type sendFn func(addr, from, subject, body string, to []string) error

func BuildHandleContactFormFn(sendFn sendFn, confirmFn confirmFn, env AppEnv) func(writer http.ResponseWriter, request *http.Request) {
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
			toList, err := toList(env.ToEmail)

			if err != nil {
				fmt.Printf("Bad email to address %s", err.Error())
				fmt.Fprintf(writer, "false")
				return
			}
			m := mailArgs{
				Addr: env.EmailHost + ":" + env.EmailPort,
				From: contactRequest.Email,
				To:   toList,
			}
			sendFn(m.Addr, m.From, contactRequest.Subject, contactRequest.Message, m.To)
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
