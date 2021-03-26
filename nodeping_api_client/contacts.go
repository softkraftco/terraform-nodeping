package nodeping_api_client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ContactNotExists struct {
	contactId string
}

func (err *ContactNotExists) Error() string {
	return fmt.Sprintf("Contact '%s' does not exist.", err.contactId)
}

func (client *Client) GetContact(Id string) (*Contact, error) {
	/*
		Returns a single contact.
	*/

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/contacts/%s", client.HostURL, Id), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(req)
	if err != nil {
		return nil, err
	}

	if string(body) == "{}" {
		e := ContactNotExists{Id}
		return nil, &e
	}

	contact := Contact{}
	err = json.Unmarshal(body, &contact)
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (client *Client) CreateContact(contact *Contact) (*Contact, error) {
	/*
		Creates a new contact, along with all needed addresses
	*/
	rb, err := contact.MarshalJSONForCreate()
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

func (client *Client) UpdateContact(contact *Contact) (*Contact, error) {
	/*
		Updates an existing contact.

		Note about addresses from nodeping documentation:
		> When updating existing addresses, the entire list is required.
		> Entries missing from the object are removed from the contact [...].
		> Adding non-existing address IDs to the list will generate an error.
	*/
	rb, err := json.Marshal(contact)
	if err != nil {
		return nil, err
	}

	// although json already contains contact "_id", the API seems to require
	// "id" this time, so it's easier to simply add id to url.
	req, err := http.NewRequest("PUT",
		fmt.Sprintf("%s/contacts/%s", client.HostURL, contact.ID),
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
	/*
		Deletes an existing contact
	*/
	req, err := http.NewRequest("DELETE",
		fmt.Sprintf("%s/contacts/%s", client.HostURL, Id), nil)
	_, err = client.doRequest(req)
	return err
}
