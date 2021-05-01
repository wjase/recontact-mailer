package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/smtp"
	"strings"
	"testing"

	"github.com/corbym/gocrest"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

var testAppEnv = AppEnv{
	toEmail:       "to@some.com",
	emailHost:     "somehost",
	emailPort:     "somePort",
	emailUsername: "emailUser",
	emailPassword: "emailPass",
}

func sendMatcherFn(matcher *gocrest.Matcher, t *testing.T) sendFn {
	return func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		then.AssertThat(t, mailArgs{Addr: addr,
			Auth: a,
			From: from,
			To:   to,
			Msg:  msg,
		}, matcher)
		return nil
	}
}

func confirmOK(remoteip, response string) (result bool, err error) {
	return true, nil
}

func confirmFailed(remoteip, response string) (result bool, err error) {
	return false, nil
}

func confirmFailedError(remoteip, response string) (result bool, err error) {
	return false, fmt.Errorf("this is an error")
}

func TestContactHandler(t *testing.T) {
	happyCasePayloadJson := `{
		"g-recaptcha-response":"bob",
		"email":"bob@bob.com",
		"subject": "a thing",
		"message": "a message"
		}`

	badEnv := testAppEnv
	badEnv.toEmail = "badbad"

	testCases := []struct {
		testName    string
		sendFn      sendFn
		confirmFn   confirmFn
		requestBody string
		env         AppEnv
	}{
		{
			testName: "happy case valid captcha",
			sendFn: sendMatcherFn(is.EqualTo(mailArgs{
				Addr: "somehost:somePort",
				Auth: smtp.PlainAuth("", "emailUser", "emailPass", "somehost"),
				From: "bob@bob.com",
				To:   []string{"to@some.com"},
				Msg: []byte(`Subject: a thing

a message`),
			}), t),
			confirmFn:   confirmOK,
			requestBody: happyCasePayloadJson,
			env:         testAppEnv,
		},
		{
			testName:    "happy case invalid captcha",
			sendFn:      sendMatcherFn(is.Nil(), t),
			confirmFn:   confirmFailed,
			requestBody: happyCasePayloadJson,
			env:         testAppEnv,
		},
		{
			testName:    "happy case invalid captcha error",
			sendFn:      sendMatcherFn(is.Nil(), t),
			confirmFn:   confirmFailedError,
			requestBody: happyCasePayloadJson,
			env:         testAppEnv,
		},
		{
			testName:    "invalid body error",
			sendFn:      sendMatcherFn(is.Nil(), t),
			requestBody: "{malformed json",
			env:         testAppEnv,
		},
		{
			testName:    "bad to email address",
			sendFn:      sendMatcherFn(is.Nil(), t),
			requestBody: happyCasePayloadJson,
			env:         badEnv,
			confirmFn:   confirmOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.testName, func(t *testing.T) {
			handler := buildHandleContactFormFn(tC.sendFn, tC.confirmFn, tC.env)
			request := http.Request{
				Host:   "AHost",
				Body:   ioutil.NopCloser(bytes.NewReader([]byte(tC.requestBody))),
				Header: http.Header{},
			}
			response := httptest.NewRecorder()
			handler(response, &request)
		})
	}
}

//EqualTo checks if two values are equal. Uses DeepEqual (could be slow).
//Like DeepEquals, if the types are not the same the matcher returns false.
//Returns a matcher that will return true if two values are equal.
func MatchesJson(expected interface{}) *gocrest.Matcher {
	match := new(gocrest.Matcher)
	match.Describe = fmt.Sprintf("value equal to <%v>", expected)
	match.Matches = func(actual interface{}) bool {
		return strings.Compare(strings.TrimSpace(asJson(actual)), strings.TrimSpace(asJson(expected))) == 0
	}
	return match
}

func asJson(o interface{}) string {
	bj, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
		return "couldn't convert:" + err.Error()
	}
	fmt.Println(string(bj))
	return string(bj)
}
