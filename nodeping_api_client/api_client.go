package nodeping_api_client

import (
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

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
	authStr := b64.StdEncoding.EncodeToString([]byte(client.Token + ":whatever"))
	req.Header.Set("Authorization", "Basic "+authStr)

	res, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
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
