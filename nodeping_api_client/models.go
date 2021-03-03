package nodeping_api_client

type Address struct {
	ID            string
	Address       string `json:"address"`
	Type          string `json:"type"`
	Suppressup    bool   `json:"accountsuppressup"`
	Suppressdown  bool   `json:"accountsuppressdown"`
	Suppressfirst bool   `json:"accountsuppressfirst"`
	Suppressall   bool   `json:"accountsuppressall"`
}

type Contact struct {
	ID         string             `json:"_id"`
	Type       string             `json:"type"`
	CustomerId string             `json:"customer_id"`
	Name       string             `json:"name"`
	Custrole   string             `json:"custrole"`
	Addresses  map[string]Address `json:"addresses"`
}
