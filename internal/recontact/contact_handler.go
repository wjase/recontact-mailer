package recontact

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var rnd = rand.New(rand.NewSource(time.Now().Unix()))

type confirmFn func(remoteip, response string) (result bool, err error)
type sendFn func(addr, from, subject, body string, to []string) error

func BuildHandleContactFormFn(sendFn sendFn, confirmFn confirmFn, env AppEnv) func(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Creating http handler")
	return func(writer http.ResponseWriter, request *http.Request) {
		contactRequest, err := NewContactRequest(request)
		if err != nil {
			fmt.Fprintf(writer, "false")
			return
		}

		bytes, _ := json.Marshal(contactRequest)
		fmt.Printf("contact: <%s>", string(bytes))

		result, err := confirmFn(GetIP(request), contactRequest.Recaptcha)
		if err != nil {
			fmt.Fprintf(writer, "false")
			fmt.Println("recaptcha failed", string(bytes))
			writer.WriteHeader(500)
			return
		}
		if result {
			// toList is list of email address that email is to be sent.
			toList, err := toList(env.ToEmail)

			if err != nil {
				fmt.Printf("Bad email to address %s", err.Error())
				fmt.Fprintf(writer, "false")
				writer.WriteHeader(400)
				return
			}
			m := mailArgs{
				Addr: env.EmailHost + ":" + env.EmailPort,
				From: contactRequest.Email,
				To:   toList,
			}

			err = sendFn(m.Addr, m.From, contactRequest.Subject, contactRequest.Message, m.To)
			if err != nil {
				fmt.Printf("Couldn't send message %s", contactRequest.Message)
				writer.WriteHeader(500)
			}
		} else {
			time.Sleep(time.Duration(rnd.Intn(8)) * time.Second)
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
