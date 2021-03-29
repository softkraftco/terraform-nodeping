package nodeping_api_client

import "encoding/json"

type Address struct {
	ID            string `json:"id,omitempty"`
	Address       string `json:"address"`
	Type          string `json:"type"`
	Suppressup    bool   `json:"suppressup"`
	Suppressdown  bool   `json:"suppressdown"`
	Suppressfirst bool   `json:"suppressfirst"`
	Suppressall   bool   `json:"suppressall"`
	// webhook attribures
	Action       string            `json:"action,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	Querystrings map[string]string `json:"querystrings,omitempty"`
	// pushover attributes
	Priority int `json:"priority"`
}

type Check struct {
	ID            string                    `json:"_id,omitempty"`
	Rev           string                    `json:"_rev,omitempty"`
	Label         string                    `json:"label,omitempty"`
	Type          string                    `json:"type,omitempty"`
	CustomerId    string                    `json:"customer_id,omitempty"`
	Interval      int                       `json:"interval,omitempty"`
	Status        string                    `json:"status,omitempty"`
	Enable        string                    `json:"enable,omitempty"`
	Public        bool                      `json:"public,omitempty"`
	Notifications []map[string]Notification `json:"notifications,omitempty"`
	Parameters    map[string]interface{}    `json:"parameters,omitempty"`
	Runlocations  []string                  `json:"runlocations,omitempty"`
	Created       int                       `json:"created,omitempty"`
	Modified      int                       `json:"modified,omitempty"`
	Queue         string                    `json:"queue,omitempty"`
	Uuid          string                    `json:"uuid,omitempty"`
	State         interface{}               `json:"state,omitempty"`
	Firstdown     int                       `json:"firstdown,omitempty"`
}

type CheckUpdate struct { // used for PUT and POST requests.
	/*
		Since checks require a different structure for PUT and POST request,
		compared	to the one received from GET requests, this is a separate struct
		for creating and updating checks.
	*/
	ID            string                    `json:"_id,omitempty"`
	Label         string                    `json:"label,omitempty"`
	CustomerId    string                    `json:"customer_id,omitempty"`
	Type          string                    `json:"type,omitempty"`
	Target        string                    `json:"target,omitempty"`
	Interval      int                       `json:"interval,omitempty"`
	Enable        string                    `json:"enabled,omitempty"` // Note this is called `enable` on GET responses
	Public        bool                      `json:"public,omitempty"`
	RunLocations  []string                  `json:"runlocations,omitempty"`
	HomeLoc       string                    `json:"homeloc,omitempty"`
	Notifications []map[string]Notification `json:"notifications,omitempty"`
	Threshold     int                       `json:"threshold,omitempty"`
	Sens          int                       `json:"sens,omitempty"`
	Dep           string                    `json:"dep,omitempty"`
	Description   string                    `json:"description,omitempty"`
	// the following are only relevant for certain types
	CheckToken     string                 `json:"checktoken,omitempty"`
	ClientCert     string                 `json:"clientcert,omitempty"`
	ContentString  string                 `json:"contentstring,omitempty"`
	Dohdot         string                 `json:"dohdot,omitempty"`
	DnsType        string                 `json:"dnstype,omitempty"`
	DnsToResolve   string                 `json:"dnstoresolve,omitempty"`
	Dnsrd          bool                   `json:"dnsrd,omitempty"`
	Transport      string                 `json:"transport,omitempty"`
	Follow         bool                   `json:"follow,omitempty"`
	Email          string                 `json:"email,omitempty"`
	Port           int                    `json:"port,omitempty"`
	Username       string                 `json:"username,omitempty"`
	Password       string                 `json:"password,omitempty"`
	Secure         string                 `json:"secure,omitempty"`
	Verify         string                 `json:"verify,omitempty"`
	Ignore         string                 `json:"ignore,omitempty"`
	Invert         string                 `json:"invert,omitempty"`
	WarningDays    int                    `json:"warningdays,omitempty"`
	Fields         map[string]CheckField  `json:"fields,omitempty"`
	Postdata       string                 `json:"postdata,omitempty"`
	Data           map[string]interface{} `json:"data,omitempty"`
	ReceiveHeaders map[string]interface{} `json:"receiveheaders,omitempty"`
	SendHeaders    map[string]interface{} `json:"sendheaders,omitempty"`
	Edns           map[string]interface{} `json:"edns,omitempty"`
	Method         string                 `json:"method,omitempty"`
	Statuscode     int                    `json:"statuscode,omitempty"`
	Ipv6           bool                   `json:"ipv6,omitempty"`
	Regex          bool                   `json:"regex,omitempty"`
	ServerName     string                 `json:"servername,omitempty"`
	Snmpv          string                 `json:"snmpv,omitempty"`
	Snmpcom        string                 `json:"snmpcom,omitempty"`
	VerifyVolume   bool                   `json:"verifyvolume,omitempty"`
	VolumeMin      int                    `json:"volumemin,omitempty"`
	WhoisServer    string                 `json:"whoisserver,omitempty"`
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
		When calling API to create a new contract, passed json object is not
		allowed to have "addresses" field, and doesn't need the "id" field.
	*/
	return json.Marshal(struct {
		CustomerId   string    `json:"customer_id,omitempty"`
		Name         string    `json:"name,omitempty"`
		Custrole     string    `json:"custrole,omitempty"`
		NewAddresses []Address `json:"newaddresses,omitempty"`
	}{c.CustomerId, c.Name, c.Custrole, c.NewAddresses})
}

type Notification struct {
	Delay    int    `json:"delay,omitempty"`
	Schedule string `json:"schedule,omitempty"`
}

type CheckField struct {
	Name string `json:"name"`
	Min  int    `json:"min"`
	Max  int    `json:"max"`
}
