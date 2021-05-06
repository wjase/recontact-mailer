package recontact

import (
	"fmt"
	"io"
	"net/smtp"
	"os"
	"os/exec"
	"regexp"

	"gopkg.in/gomail.v2"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

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

//ex: SendMail("127.0.0.1:25", (&mail.Address{"from name", "from@example.com"}).String(), "Email Subject", "message body", []string{(&mail.Address{"to name", "to@example.com"}).String()})
func SendMail(addr, from, subject, body string, to []string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	s := gomail.SendFunc(submitMail)

	return gomail.Send(s, m)

}

const sendmail = "/usr/sbin/sendmail"

func submitMail(from string, to []string, m io.WriterTo) (err error) {
	cmd := exec.Command(sendmail, "-t")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pw, err := cmd.StdinPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	var errs [3]error
	_, errs[0] = m.WriteTo(pw)
	errs[1] = pw.Close()
	errs[2] = cmd.Wait()
	for _, err = range errs {
		if err != nil {
			return
		}
	}
	return
}
