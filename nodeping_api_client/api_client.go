package nodeping_api_client

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"time"
)

// Default HostURL
const HostURL string = "https://api.nodeping.com/api/1"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

func NewClient(token string) *Client {
	client := Client{
		HostURL:    HostURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		Token:      token,
	}

	return &client
}

func (client *Client) doRequest(req *http.Request) ([]byte, error) {
	// set request authentication.
	// from API docs: The password will be ignored so you can either leave it
	// blank or pass a random string
	authStr := b64.StdEncoding.EncodeToString([]byte(client.Token + ":whatever"))
	req.Header.Set("Authorization", "Basic "+authStr)

	// send
	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	// cleanup
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (client *Client) doRequest2(ctx context.Context, method string, url string, data []byte) ([]byte, error) {
	request, err := client.prepareRequest2(ctx, method, url, data)
	if err != nil {
		return nil, err
	}
	body, err := client.sendRequest(request)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (client *Client) prepareRequest(method string, url string, data []byte) (*http.Request, error) {
	// prepare data to be sent
	var dataReader *strings.Reader
	if data != nil {
		dataReader = strings.NewReader(string(data))
	}

	// initialize request
	reqest, err := http.NewRequest(method, url, dataReader)
	if err != nil {
		return nil, err
	}

	// set request authentication.
	// from API docs: The password will be ignored so you can either leave it
	// blank or pass a random string
	authStr := b64.StdEncoding.EncodeToString([]byte(client.Token + ":whatever"))
	reqest.Header.Set("Authorization", "Basic "+authStr)

	// set json content type for request that have data to be sent
	if data != nil {
		reqest.Header.Set("Content-Type", "application/json")
	}

	return reqest, nil
}

func (client *Client) prepareRequest2(ctx context.Context, method string, url string, data []byte) (*http.Request, error) {
	// prepare data to be sent
	var body *strings.Reader
	if data != nil {
		body = strings.NewReader(string(data))
	} else {
		body = strings.NewReader("")
	}

	// initialize request
	reqest, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// set request authentication.
	// from API docs: The password will be ignored so you can either leave it
	// blank or pass a random string
	authStr := b64.StdEncoding.EncodeToString([]byte(client.Token + ":whatever"))
	reqest.Header.Set("Authorization", "Basic "+authStr)

	// set json content type for request that have data to be sent
	if data != nil {
		reqest.Header.Set("Content-Type", "application/json")
	}

	return reqest, nil
}

func (client *Client) sendRequest(request *http.Request) ([]byte, error) {
	// send
	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}

	// cleanup
	defer response.Body.Close()

	// read response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// error handling
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", response.StatusCode, body)
	}

	return body, err
}
