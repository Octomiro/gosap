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

type PurchaseDeliveryNote struct {
	DocNum        int    `json:"DocNum,omitempty"`
	DocEntry      int    `json:"DocEntry,omitempty"`
	DocType       string `json:"DocType,omitempty"`
	CardCode      string `json:",omitempty"`
	Status        string `json:"DocumentStatus,omitempty"`
	DocumentLines []PurchaseDeliveryNoteLine
}

type PurchaseDeliveryNoteLine struct {
	LineNum         int     `json:",omitempty"`
	ItemCode        string  `json:",omitempty"`
	ItemDescription string  `json:",omitempty"`
	Quantity        float64 `json:",omitempty"`
	ShipDate        string  `json:",omitempty"`
	Price           float64 `json:",omitempty"`
	BaseType        int     `json:",omitempty"`
	BaseEntry       int     `json:",omitempty"`
	BaseLine        int
}

type (
	DeliveryNote      = Document
	DeliveryNoteLine  = DocumentLine
	PurchaseOrder     = Document
	PurchaseOrderLine = DocumentLine
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

type InventoryCountingLine struct {
	ItemCode        string  `json:"ItemCode,omitempty"`
	WarehouseCode   string  `json:"WarehouseCode,omitempty"`
	CountedQuantity float64 `json:"CountedQuantity,omitempty"`
	LineNum         int     `json:"LineNumber,omitempty"`
	ItemDescription string  `json:"ItemDescription,omitempty"`
	BinEntry        int     `json:"BinEntry,omitempty"`
}

type InventoryCounting struct {
	DocumentEntry          int                     `json:"DocumentEntry,omitempty"`
	DocumentNumber         int                     `json:"DocumentNumber,omitempty"`
	Series                 int                     `json:"Series,omitempty"`
	CountingType           string                  `json:"CountingType,omitempty"`
	DocumentStatus         string                  `json:"DocumentStatus,omitempty"`
	InventoryCountingLines []InventoryCountingLine `json:"InventoryCountingLines,omitempty"`
}

type InventoryCountingResponse struct {
	ODataMetadata string              `json:"@odata.metadata"`
	Value         []InventoryCounting `json:"value"`
	NextLink      *string             `json:"odata.nextLink"`
}

type BinLocation struct {
	AbsEntry    int     `json:"AbsEntry,omitempty"`
	Warehouse   string  `json:"Warehouse,omitempty"`
	BinCode     string  `json:"BinCode,omitempty"`
	Description *string `json:"Description,omitempty"`
	MinimumQty  float64 `json:"MinimumQty,omitempty"`
	MaximumQty  float64 `json:"MaximumQty,omitempty"`
}

type BinLocationsResponse struct {
	Metadata string        `json:"odata.metadata,omitempty"` //nolint:tagliatelle
	Value    []BinLocation `json:"value,omitempty"`
	NextLink *string       `json:"odata.nextLink,omitempty"` //nolint:tagliatelle
}
