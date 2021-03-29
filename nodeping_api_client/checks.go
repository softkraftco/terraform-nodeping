package nodeping_api_client

import (
	"context"
	"encoding/json"
	"fmt"
)

// TODO: does not exist

func (client *Client) GetCheck(ctx context.Context, Id string) (*Check, error) {
	/*
		Returns a check.
	*/

	body, err := client.doRequest2(ctx, "GET", fmt.Sprintf("%s/checks/%s", client.HostURL, Id), nil)
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
	requestBody, err := json.Marshal(check)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest2(ctx, "POST", fmt.Sprintf("%s/checks", client.HostURL), requestBody)
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

	body, err := client.doRequest2(ctx, "PUT", fmt.Sprintf("%s/checks/%s", client.HostURL, check.ID), requestBody)
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

func (client *Client) DeleteCheck(ctx context.Context, Id string) error {
	/*
		Deletes an existing check.
	*/
	_, err := client.doRequest2(ctx, "DELETE", fmt.Sprintf("%s/checks/%s", client.HostURL, Id), nil)
	return err
}
