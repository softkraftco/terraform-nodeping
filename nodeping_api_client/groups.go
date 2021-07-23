package nodeping_api_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type GroupDoesNotExist struct {
	groupId string
}

func (err *GroupDoesNotExist) Error() string {
	return fmt.Sprintf("Group '%s' does not exist.", err.groupId)
}

func (client *Client) GetGroup(ctx context.Context, Id string) (*Group, error) {
	body, err := client.doRequest(ctx, http.MethodGet, fmt.Sprintf("%s/contactgroups/%s", client.HostURL, Id), nil)
	if err != nil {
		return nil, err
	}

	if string(body) == "{}" {
		e := GroupDoesNotExist{Id}
		return nil, &e
	}

	group := Group{}
	err = json.Unmarshal(body, &group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (client *Client) CreateGroup(ctx context.Context, group *Group) (*Group, error) {
	rb, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(ctx, http.MethodPost, fmt.Sprintf("%s/contactgroups", client.HostURL), rb)
	if err != nil {
		return nil, err
	}

	newGroup := Group{}
	err = json.Unmarshal(body, &newGroup)
	if err != nil {
		return nil, err
	}

	return &newGroup, nil
}

func (client *Client) UpdateGroup(ctx context.Context, group *Group) (*Group, error) {
	rb, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	// although json already contains contact "_id", the API seems to require
	// "id" this time, so it's easier to simply add id to url.
	body, err := client.doRequest(ctx, "PUT",
		fmt.Sprintf("%s/contactgroups/%s", client.HostURL, group.ID), rb)
	if err != nil {
		return nil, err
	}

	newGroup := Group{}
	err = json.Unmarshal(body, &newGroup)
	if err != nil {
		return nil, err
	}

	return &newGroup, nil
}

func (client *Client) DeleteGroup(ctx context.Context, Id string) error {
	_, err := client.doRequest(ctx, "DELETE", fmt.Sprintf("%s/contactgroups/%s", client.HostURL, Id), nil)
	return err
}
