package nodeping_api_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CustomerDoesNotExist struct {
	customerId string
}

func (err *CustomerDoesNotExist) Error() string {
	return fmt.Sprintf("Customer '%s' does not exist.", err.customerId)
}

func (client *Client) GetCustomer(ctx context.Context, Id string) (*Customer, error) {
	body, err := client.doRequest(ctx, http.MethodGet, fmt.Sprintf("%s/accounts/%s", client.HostURL, Id), nil)
	if err != nil {
		return nil, err
	}

	if string(body) == "{}" {
		e := CustomerDoesNotExist{Id}
		return nil, &e
	}

	customer := Customer{}
	err = json.Unmarshal(body, &customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (client *Client) CreateCustomer(ctx context.Context, customer *Customer) (*Customer, error) {
	rb, err := customer.MarshalJSONForCreate()

	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(ctx, http.MethodPost, fmt.Sprintf("%s/accounts", client.HostURL), rb)
	if err != nil {
		return nil, err
	}

	newCustomer := Customer{}
	err = json.Unmarshal(body, &newCustomer)
	if err != nil {
		return nil, err
	}

	return &newCustomer, nil
}

func (client *Client) UpdateCustomer(ctx context.Context, customer *Customer) (*Customer, error) {
	rb, err := customer.MarshalJSONForCreate()
	if err != nil {
		return nil, err
	}

	// although json already contains contact "_id", the API seems to require
	// "id" this time, so it's easier to simply add id to url.
	body, err := client.doRequest(ctx, "PUT",
		fmt.Sprintf("%s/accounts/%s", client.HostURL, customer.ID), rb)
	if err != nil {
		return nil, err
	}

	newCustomer := Customer{}
	err = json.Unmarshal(body, &newCustomer)
	if err != nil {
		return nil, err
	}

	return &newCustomer, nil
}

func (client *Client) DeleteCustomer(ctx context.Context, Id string) error {
	_, err := client.doRequest(ctx, "DELETE", fmt.Sprintf("%s/accounts?customerid=%s", client.HostURL, Id), nil)
	return err
}
