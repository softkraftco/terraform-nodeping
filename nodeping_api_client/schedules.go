package nodeping_api_client

import (
	"context"
	"encoding/json"
	"fmt"
)

// TODO: does not exist

func (client *Client) GetSchedule(ctx context.Context, Name string) (*Schedule, error) {
	/*
		Returns a schedule.
	*/

	body, err := client.doRequest2(ctx, "GET", fmt.Sprintf("%s/schedules/%s", client.HostURL, Name), nil)
	if err != nil {
		return nil, err
	}

	schedule := Schedule{}
	err = json.Unmarshal(body, &schedule.Data)
	if err != nil {
		return nil, err
	}

	// since there is no name in API response, set name manually.
	schedule.Name = Name

	return &schedule, nil
}

func (client *Client) CreateSchedule(ctx context.Context, schedule *Schedule) (string, error) {
	/*
		Creates a new schedule.
		Returns a new chec object based on API response.

		API has no separate endpoint for create, so this actually just calls update.
	*/
	return client.UpdateSchedule(ctx, schedule)

}

func (client *Client) UpdateSchedule(ctx context.Context, schedule *Schedule) (string, error) {
	/*
		Updates an existing schedule.
		Returns an updated version, as given in API response.
	*/
	requestBody, err := json.Marshal(schedule)
	if err != nil {
		return "", err
	}

	body, err := client.doRequest2(ctx, "PUT", fmt.Sprintf("%s/schedules/%s", client.HostURL, schedule.Name), requestBody)
	if err != nil {
		return "", err
	}

	responseContent := make(map[string]interface{})
	err = json.Unmarshal(body, &responseContent)
	if err != nil {
		return "", err
	}

	return responseContent["id"].(string), nil
}

func (client *Client) DeleteSchedule(ctx context.Context, Name string) error {
	/*
		Deletes an existing schedule.
	*/
	_, err := client.doRequest2(ctx, "DELETE", fmt.Sprintf("%s/schedules/%s", client.HostURL, Name), nil)
	return err
}
