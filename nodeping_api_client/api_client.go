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

func (client *Client) doRequest(ctx context.Context, method string, url string, data []byte) ([]byte, error) {
	request, err := client.prepareRequest(ctx, method, url, data)
	if err != nil {
		return nil, err
	}
	body, err := client.sendRequest(request)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (client *Client) prepareRequest(ctx context.Context, method string, url string, data []byte) (*http.Request, error) {
	// prepare data to be sent
	var body *strings.Reader
	if data != nil {
		body = strings.NewReader(string(data))
	} else {
		body = strings.NewReader("")
	}

	// initialize request
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// set request authentication.
	// from API docs: The password will be ignored so you can either leave it
	// blank or pass a random string
	authStr := b64.StdEncoding.EncodeToString([]byte(client.Token + ":whatever"))
	request.Header.Set("Authorization", "Basic "+authStr)

	// set json content type for request that have data to be sent
	if data != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	return request, nil
}

func (client *Client) sendRequest(request *http.Request) ([]byte, error) {
	// send
	println(">>> 0", request.URL.String(), request.Method)
	cpy, _ := request.GetBody()
	data, _ := ioutil.ReadAll(cpy)
	println(string(data))
	response, err := client.HTTPClient.Do(request)
	if err != nil {
		println(">>> 1x", err.Error())
		return nil, err
	}

	println(">>> 2", response.Status)
	// cleanup
	defer response.Body.Close()
	println(">>> 3")
	// read response
	body, err := ioutil.ReadAll(response.Body)
	println(">>> 4")
	if err != nil {
		return nil, err
	}
	println(">>> 5")
	println(string(body))
	println("↑↑↑↑")
	// error handling
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", response.StatusCode, body)
	} else if strings.HasPrefix(string(body), "{\"error\"") {
		// It happens that nodeping API resonses with status 200, but an error
		// written into response body (ie. `{"error":"A target is required."}`).
		// Raise an error in such a case.
		return nil, fmt.Errorf("response content indicates an error: `%s`", body)
	}

	return body, err
}
