package nodeping_api_client

import (
	"context"
	"encoding/json"
	"fmt"
)

type ContactDoesNotExist struct {
	contactId string
}

func (err *ContactDoesNotExist) Error() string {
	return fmt.Sprintf("Contact '%s' does not exist.", err.contactId)
}

func (client *Client) GetContact(ctx context.Context, Id string) (*Contact, error) {
	/*
		Returns a single contact.
	*/

	body, err := client.doRequest(ctx, "GET", fmt.Sprintf("%s/contacts/%s", client.HostURL, Id), nil)
	if err != nil {
		return nil, err
	}

	if string(body) == "{}" {
		e := ContactDoesNotExist{Id}
		return nil, &e
	}

	contact := Contact{}
	err = json.Unmarshal(body, &contact)
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (client *Client) CreateContact(ctx context.Context, contact *Contact) (*Contact, error) {
	/*
		Creates a new contact, along with all needed addresses
	*/
	rb, err := contact.MarshalJSONForCreate()
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(ctx, "POST", fmt.Sprintf("%s/contacts", client.HostURL), rb)
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

func (client *Client) UpdateContact(ctx context.Context, contact *Contact) (*Contact, error) {
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
	body, err := client.doRequest(ctx, "PUT",
		fmt.Sprintf("%s/contacts/%s", client.HostURL, contact.ID), rb)
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

func (client *Client) DeleteContact(ctx context.Context, Id string) error {
	/*
		Deletes an existing contact
	*/
	_, err := client.doRequest(ctx, "DELETE", fmt.Sprintf("%s/contacts/%s", client.HostURL, Id), nil)
	return err
}
