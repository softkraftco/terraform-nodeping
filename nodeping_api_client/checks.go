package nodeping_api_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *Client) GetCheck(ctx context.Context, customerID, Id string) (*Check, error) {
	/*
		Returns a check.
	*/
	body, err := client.doRequest(ctx, http.MethodGet, fmt.Sprintf("%s/checks/?id=%s&customerid=%s", client.HostURL, Id, customerID), nil)
	if err != nil {
		return nil, err
	}

	check := Check{}
	err = json.Unmarshal(body, &check)
	if err != nil {
		return nil, err
	}

	return &check, nil
}

func (client *Client) CreateCheck(ctx context.Context, checkUpdate *CheckUpdate) (*Check, error) {
	/*
		Creates a new check.
		Returns a new chec object based on API response.
	*/
	requestBody, err := json.Marshal(checkUpdate)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(ctx, http.MethodPost, fmt.Sprintf("%s/checks", client.HostURL), requestBody)
	if err != nil {
		return nil, err
	}

	newCheck := Check{}
	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	return &newCheck, nil
}

func (client *Client) UpdateCheck(ctx context.Context, check *CheckUpdate) (*Check, error) {
	/*
		Updates an existing check.
		Returns an updated version, as given in API response.
	*/
	requestBody, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(ctx, "PUT", fmt.Sprintf("%s/checks/%s", client.HostURL, check.ID), requestBody)
	if err != nil {
		return nil, err
	}

	newCheck := Check{}
	err = json.Unmarshal(body, &newCheck)
	if err != nil {
		return nil, err
	}

	return &newCheck, nil
}

func (client *Client) DeleteCheck(ctx context.Context, customerId, Id string) error {
	/*
		Deletes an existing check.
	*/
	_, err := client.doRequest(ctx, "DELETE", fmt.Sprintf("%s/checks/?id=%s&customerid=%s", client.HostURL, Id, customerId), nil)
	return err
}
