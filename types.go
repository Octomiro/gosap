package gosap

type Item struct {
	ItemCode          string
	ItemName          string
	PurchaseUnitWidth *float64
}

type DeliveryNoteLine struct {
	LineNum          int
	ItemCode         string  `json:",omitempty"`
	ItemDescription  string  `json:",omitempty"`
	Quantity         float64 `json:",omitempty"`
	SelectedQuantity float64 `json:"U_SelectedQuantity,omitempty"`
	ShipDate         string  `json:",omitempty"`
	Price            float64 `json:",omitempty"`
}

type DeliveryNote struct {
	DocNum        int    `json:"DocNum,omitempty"`
	DocEntry      int    `json:"DocEntry,omitempty"`
	DocType       string `json:"DocType,omitempty"`
	CardCode      string `json:",omitempty"`
	DocumentLines []DeliveryNoteLine
}

type DeliveryNotes struct {
	Metadata string         `json:"odata.metadata"` //nolint:tagliatelle
	Value    []DeliveryNote `json:"value"`
	NextLink *string        `json:"odata.nextLink"` //nolint:tagliatelle
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
