package gosap

type Item struct {
	ItemCode          string
	ItemName          string
	PurchaseUnitWidth *float64
}

type DocumentLine struct {
	LineNum          int
	ItemCode         string  `json:",omitempty"`
	ItemDescription  string  `json:",omitempty"`
	Quantity         float64 `json:",omitempty"`
	SelectedQuantity float64 `json:"U_SelectedQuantity,omitempty"`
	ShipDate         string  `json:",omitempty"`
	Price            float64 `json:",omitempty"`
}

type Document struct {
	DocNum        int    `json:"DocNum,omitempty"`
	DocEntry      int    `json:"DocEntry,omitempty"`
	DocType       string `json:"DocType,omitempty"`
	CardCode      string `json:",omitempty"`
	Status        string `json:"DocumentStatus,omitempty"`
	DocumentLines []DocumentLine
}

type (
	DeliveryNote             = Document
	DeliveryNoteLine         = DocumentLine
	PurchaseOrder            = Document
	PurchaseOrderLine        = DocumentLine
	PurchaseDeliveryNote     = Document
	PurchaseDeliveryNoteLine = DocumentLine
)

func (dn *DeliveryNote) IsOpen() bool {
	return dn.Status == "bost_Open"
}

func (dn *DeliveryNote) IsClosed() bool {
	return dn.Status == "bost_Close"
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

type PurchaseOrders struct {
	Metadata string          `json:"odata.metadata"` //nolint:tagliatelle
	Value    []PurchaseOrder `json:"value"`
	NextLink *string         `json:"odata.nextLink"` //nolint:tagliatelle
}

type PurchaseDeliveryNotes struct {
	Metadata string                 `json:"odata.metadata"` //nolint:tagliatelle
	Value    []PurchaseDeliveryNote `json:"value"`
	NextLink *string                `json:"odata.nextLink"` //nolint:tagliatelle
}
