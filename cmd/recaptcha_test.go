package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
)

func TestCreateContact(t *testing.T) {
	request := http.Request{
		Host: "AHost",
		Body: ioutil.NopCloser(bytes.NewReader([]byte(`{
			"g-recaptcha-response":"bob",
			"email":"bob@bob.com",
			"subject": "a thing",
			"message": "a message"
			}`))),
	}
	contact, err := NewContactRequest(&request)
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, contact.Recaptcha, is.EqualTo("bob"))
	then.AssertThat(t, contact.Subject, is.EqualTo("a thing"))
	then.AssertThat(t, contact.Email, is.EqualTo("bob@bob.com"))
	then.AssertThat(t, contact.Message, is.EqualTo("a message"))
}

func TestCreateContactError(t *testing.T) {
	request := http.Request{
		Host: "AHost",
		Body: ioutil.NopCloser(bytes.NewReader([]byte(`{
			`))),
	}
	_, err := NewContactRequest(&request)
	then.AssertThat(t, err, is.Not(is.Nil()))
}

func TestCreateContactProxied(t *testing.T) {
	request := http.Request{
		Host:   "AHost",
		Header: http.Header{},
		Body:   ioutil.NopCloser(bytes.NewReader([]byte(``))),
	}
	request.Header.Add("X-FORWARDED-FOR", "someIP")
	ip := GetIP(&request)
	then.AssertThat(t, ip, is.EqualTo("someIP"))
}
