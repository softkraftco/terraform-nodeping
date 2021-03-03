package nodeping_api_client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *Client) GetContact(Id string) (*Contact, error) {
	/*
		Returns a list of all contacts
	*/
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/contacts/%s", client.HostURL, Id), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	contact := Contact{}
	err = json.Unmarshal(body, &contact)
	if err != nil {
		return nil, err
	}

	return &contact, nil
}
