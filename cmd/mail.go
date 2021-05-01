package main

import (
	"fmt"
	"net/smtp"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type mailArgs struct {
	Addr string
	Auth smtp.Auth
	From string
	To   []string
	Msg  []byte
}

func (m mailArgs) String() string {
	return fmt.Sprintf(`Addr:%s Auth:%v From:%s To:%s Msg:%s`, m.Addr, m.Auth, m.From, m.To, string(m.Msg))
}

func toList(addr string) ([]string, error) {
	if emailRegex.MatchString(addr) {
		return []string{addr}, nil
	}
	return []string{}, fmt.Errorf("malformed email")
}

func buildBody(subject, body string) []byte {
	return []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))
}
