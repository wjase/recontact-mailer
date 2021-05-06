package recontact

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var whitespaceRemover = strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

type mailArgs struct {
	Addr string
	Auth smtp.Auth
	From string
	To   []string
	Msg  string
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

func buildBody(to []string, from, subject, body string) []byte {
	return []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" + base64.StdEncoding.EncodeToString([]byte(body)))
}

//ex: SendMail("127.0.0.1:25", (&mail.Address{"from name", "from@example.com"}).String(), "Email Subject", "message body", []string{(&mail.Address{"to name", "to@example.com"}).String()})
func SendMail(addr, from, subject, body string, to []string) error {

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(whitespaceRemover.Replace(from)); err != nil {
		return err
	}
	for i := range to {
		to[i] = whitespaceRemover.Replace(to[i])
		if err = c.Rcpt(to[i]); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(buildBody(to, from, subject, body))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
