package nodeping_api_client

type Address struct {
	ID            string `json:"id,omitempty"`
	Address       string `json:"address"`
	Type          string `json:"type"`
	Suppressup    bool   `json:"accountsuppressup"`
	Suppressdown  bool   `json:"accountsuppressdown"`
	Suppressfirst bool   `json:"accountsuppressfirst"`
	Suppressall   bool   `json:"accountsuppressall"`
	// webhook attribures
	Action       string            `json:"action,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	Querystrings map[string]string `json:"querystrings,omitempty"`
	// pushover attributes
	Priority int `json:"priority"`
}

type Contact struct {
	ID           string             `json:"_id,omitempty"`
	Type         string             `json:"type,omitempty"`
	CustomerId   string             `json:"customer_id,omitempty"`
	Name         string             `json:"name,omitempty"`
	Custrole     string             `json:"custrole,omitempty"`
	Addresses    map[string]Address `json:"addresses,omitempty"`
	NewAddresses []Address          `json:"newaddresses,omitempty"`
}
