package nodeping_api_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ScheduleDoesNotExist struct {
	scheduleName string
}

func (err *ScheduleDoesNotExist) Error() string {
	return fmt.Sprintf("Schedule '%s' does not exist.", err.scheduleName)
}

func (client *Client) GetSchedule(ctx context.Context, customerId, name string) (*Schedule, error) {
	/*
		Returns a schedule.
	*/

	body, err := client.doRequest(ctx, http.MethodGet, fmt.Sprintf("%s/schedules/?id=%s&customerid=%s", client.HostURL, name, customerId), nil)
	if err != nil {
		return nil, err
	}

	if string(body) == "\"\"" {
		e := ScheduleDoesNotExist{name}
		return nil, &e
	}

	schedule := Schedule{}
	err = json.Unmarshal(body, &schedule.Data)
	if err != nil {
		return nil, err
	}

	// since there is no name in API response, set name manually.
	schedule.Name = name
	schedule.CustomerId = customerId

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
	requestBody, err := schedule.MarshalJSONForCreate()

	if err != nil {
		return "", err
	}

	body, err := client.doRequest(ctx, "PUT", fmt.Sprintf("%s/schedules/%s", client.HostURL, schedule.Name), requestBody)
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

func (client *Client) DeleteSchedule(ctx context.Context, customerId, name string) error {
	/*
		Deletes an existing schedule.
	*/
	_, err := client.doRequest(ctx, "DELETE", fmt.Sprintf("%s/schedules/?id=%s&customerId=%s", client.HostURL, name, customerId), nil)
	return err
}
