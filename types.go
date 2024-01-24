package gosap

type Item struct {
	ItemCode          string
	ItemName          string
	PurchaseUnitWidth *float64
}

type Items struct {
	Metadata string  `json:"odata.metadata"` //nolint:tagliatelle
	Value    []Item  `json:"value"`
	NextLink *string `json:"odata.nextLink"` //nolint:tagliatelle
}

type BusinessPartner struct {
	CardCode string
	CardName string
}

type (
	Supplier = BusinessPartner
	Client   = BusinessPartner
)

type BusinessPartners struct {
	Metadata string            `json:"odata.metadata"` //nolint:tagliatelle
	Value    []BusinessPartner `json:"value"`
	NextLink *string           `json:"odata.nextLink"` //nolint:tagliatelle
}

type (
	Suppliers = BusinessPartners
	Clients   = BusinessPartners
)

type PurchaseDeliveryNoteLine struct {
	ItemCode  string
	Quantity  string
	TaxCode   string
	UnitPrice *string
}

type PurchaseDeliveryNotes struct {
	CardCode      string
	DocumentLines []PurchaseDeliveryNoteLine
}
