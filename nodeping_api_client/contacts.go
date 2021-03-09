package nodeping_api_client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func (client *Client) CreateContact(contact *NewContact) (*Contact, error) {
	/*
		Creates a new contact, along with all needed addresses
	*/

	rb, err := json.Marshal(contact)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/contacts", client.HostURL),
		strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	newContact := Contact{}
	err = json.Unmarshal(body, &newContact)
	if err != nil {
		return nil, err
	}

	return &newContact, nil
}

func (client *Client) DeleteContact(Id string) error {
	req, err := http.NewRequest("DELETE",
		fmt.Sprintf("%s/contacts/%s", client.HostURL, Id), nil)
	_, err = client.doRequest(req)
	return err
}
