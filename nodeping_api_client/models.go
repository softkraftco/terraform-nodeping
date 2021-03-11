package nodeping_api_client

type Address struct {
	ID            string
	Address       string            `json:"address"`
	Type          string            `json:"type"`
	Suppressup    bool              `json:"accountsuppressup"`
	Suppressdown  bool              `json:"accountsuppressdown"`
	Suppressfirst bool              `json:"accountsuppressfirst"`
	Suppressall   bool              `json:"accountsuppressall"`
	Action        string            `json:"action"`
	Data          map[string]string `json:"data"`
	Headers       map[string]string `json:"headers"`
	Querystrings  map[string]string `json:"querystrings"`
}

type Contact struct {
	ID         string             `json:"_id"`
	Type       string             `json:"type"`
	CustomerId string             `json:"customer_id"`
	Name       string             `json:"name"`
	Custrole   string             `json:"custrole"`
	Addresses  map[string]Address `json:"addresses"`
}

type NewContact struct {
	/*
		New contacts don't have IDs, and use "newaddresses" array instead of
		"addresses" map.
	*/
	Type       string    `json:"type"`
	CustomerId string    `json:"customer_id"`
	Name       string    `json:"name"`
	Custrole   string    `json:"custrole"`
	Addresses  []Address `json:"newaddresses"`
}
