package recontact

import (
	"fmt"
	"log"
	"os"
)

var logger = log.Default()

type AppEnv struct {
	PrivateKey string
	ToEmail    string
	AdminEmail string
	EmailHost  string
	EmailPort  string
	AppPort    string
	Endpoint   string
}

func ensureEnvNotBlank(name string) string {
	if val, ok := os.LookupEnv(name); !ok {
		logger.Printf("Unexpected blank property %s\n", name)
		panic(fmt.Sprintf("Unexpected blank property %s\n", name))
	} else {
		return val
	}
}

// NewAppEnv cerates a new env.
func NewAppEnv() AppEnv {

	return AppEnv{
		PrivateKey: ensureEnvNotBlank("RECAPTCHA_PRIVATE_KEY"),
		ToEmail:    ensureEnvNotBlank("TO_MAIL"),
		AdminEmail: ensureEnvNotBlank("ADMIN_MAIL"),
		EmailHost:  ensureEnvNotBlank("EMAIL_HOST"),
		EmailPort:  ensureEnvNotBlank("EMAIL_PORT"),
		AppPort:    ensureEnvNotBlank("APP_PORT"),
		Endpoint:   ensureEnvNotBlank("ENDPOINT"),
	}
}
