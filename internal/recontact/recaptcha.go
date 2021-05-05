package recontact

import (
	"encoding/json"
	"net/http"
)

type ContactRequest struct {
	Recaptcha string `json:"g-recaptcha-response"`
	RemoteIP  string
	Subject   string
	Email     string
	Message   string
}

func NewContactRequest(request *http.Request) (*ContactRequest, error) {
	var contactRequest ContactRequest
	err := json.NewDecoder(request.Body).Decode(&contactRequest)

	if err != nil {
		return nil, err
	}

	return &contactRequest, nil
}
