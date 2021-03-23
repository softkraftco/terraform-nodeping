package nodeping_api_client

import "encoding/json"

type Address struct {
	ID            string `json:"id,omitempty"`
	Address       string `json:"address"`
	Type          string `json:"type"`
	Suppressup    bool   `json:"accountsuppressup"`
	Suppressdown  bool   `json:"accountsuppressdown"`
	Suppressfirst bool   `json:"accountsuppressfirst"`
	Suppressall   bool   `json:"accountsuppressall"`
	// webhook attribures
	Action       string            `json:"action,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	Querystrings map[string]string `json:"querystrings,omitempty"`
	// pushover attributes
	Priority int `json:"priority"`
}

type Contact struct {
	/*
		Note that "addresses" can't be omitted from json, even if it's empty, as
		an empty "addresses" map might mean that some addresses should be
		removed.
	*/
	ID           string             `json:"_id,omitempty"`
	Type         string             `json:"type,omitempty"`
	CustomerId   string             `json:"customer_id,omitempty"`
	Name         string             `json:"name,omitempty"`
	Custrole     string             `json:"custrole,omitempty"`
	Addresses    map[string]Address `json:"addresses"`
	NewAddresses []Address          `json:"newaddresses,omitempty"`
}

func (c *Contact) MarshalJSONForCreate() ([]byte, error) {
	/*
		When calling API to create a new contract, passed json coject is not
		allowed to have "addresses" field, and doesn't need the "id" field.
	*/
	return json.Marshal(struct {
		CustomerId   string    `json:"customer_id,omitempty"`
		Name         string    `json:"name,omitempty"`
		Custrole     string    `json:"custrole,omitempty"`
		NewAddresses []Address `json:"newaddresses,omitempty"`
	}{c.CustomerId, c.Name, c.Custrole, c.NewAddresses})
}
