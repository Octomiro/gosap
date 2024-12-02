package gosap

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Session struct {
	B1Session string
	RouteID   string
}

func Authenticate(cfg Config) (*Session, error) {
	loginPayload, err := cfg.LoginPayload()
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Post(cfg.LoginEndpoint(), "application/json", strings.NewReader(loginPayload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body due to %s", err)
	}

	statusOK := resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
	if !statusOK {
		return nil, fmt.Errorf("request to SAP API (%s) was not successful due to %s - %s", cfg.LoginEndpoint(), resp.Status, string(content))
	}

	cookies := resp.Cookies()
	session := Session{}

	for _, cookie := range cookies {
		if cookie.Name == "B1SESSION" {
			session.B1Session = cookie.Value
		}

		if cookie.Name == "ROUTEID" {
			session.RouteID = cookie.Value
		}
	}

	return &session, nil
}

func (s *Session) setSessionCookies(req *http.Request) {
	req.AddCookie(&http.Cookie{Name: "B1SESSION", Value: s.B1Session})
	req.AddCookie(&http.Cookie{Name: "ROUTEID", Value: s.RouteID})
}

// Do sends the request and returns the response.
// Caller should close Body of response after reading it.
func (s *Session) Do(req *http.Request) (*http.Response, []byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	s.setSessionCookies(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, []byte{}, err
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, []byte{}, fmt.Errorf("could not read body of response due to %s", err)
	}

	statusOK := resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
	if !statusOK {
		return nil, content, fmt.Errorf("request to SAP API (%s) was not successful due to %s - %s", req.URL, resp.Status, string(content))
	}

	return resp, content, nil
}

func (s *Session) GetItem(cfg Config, id string) (*Item, error) {
	req, err := http.NewRequest(http.MethodGet, cfg.GetItemEndpoint(id), nil)
	if err != nil {
		return nil, err
	}

	_, content, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	var item Item
	if err := json.Unmarshal(content, &item); err != nil {
		return nil, fmt.Errorf("could not load json response due to %s", err)
	}

	return &item, nil
}

func (s *Session) getItems(cfg Config, endpoint string) (*Items, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	_, content, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	var items Items
	if err := json.Unmarshal(content, &items); err != nil {
		return nil, fmt.Errorf("could not load json response due to %s", err)
	}

	if items.NextLink != nil && *items.NextLink != "" {
		next, err := s.getItems(cfg, cfg.BuildEndpoint(*items.NextLink))
		if err != nil {
			return &items, err
		}

		items.Value = append(items.Value, next.Value...)
	}

	return &items, nil
}

func (s *Session) GetItems(cfg Config) (*Items, error) {
	return s.getItems(cfg, cfg.GetItemsEndpoint())
}

func (s *Session) GetSuppliers(cfg Config) (*Suppliers, error) {
	return s.getSuppliers(cfg, cfg.GetSuppliersEndpoint())
}

func (s *Session) getSuppliers(cfg Config, endpoint string) (*Suppliers, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	_, content, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	var suppliers Suppliers
	if err := json.Unmarshal(content, &suppliers); err != nil {
		return nil, fmt.Errorf("could not load json response due to %s", err)
	}

	if suppliers.NextLink != nil && *suppliers.NextLink != "" {
		next, err := s.getSuppliers(cfg, cfg.BuildEndpoint(*suppliers.NextLink))
		if err != nil {
			return &suppliers, err
		}

		suppliers.Value = append(suppliers.Value, next.Value...)
	}

	return &suppliers, nil
}

func (s *Session) GetClients(cfg Config) (*Clients, error) {
	return s.getClients(cfg, cfg.GetClientsEndpoint())
}

func (s *Session) getClients(cfg Config, endpoint string) (*Clients, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	_, content, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	var clients Clients
	if err := json.Unmarshal(content, &clients); err != nil {
		return nil, fmt.Errorf("could not load json response due to %s", err)
	}

	if clients.NextLink != nil && *clients.NextLink != "" {
		next, err := s.getSuppliers(cfg, cfg.BuildEndpoint(*clients.NextLink))
		if err != nil {
			return &clients, err
		}

		clients.Value = append(clients.Value, next.Value...)
	}

	return &clients, nil
}

func (s *Session) GetDeliveryNotes(cfg Config) (*DeliveryNotes, error) {
	return s.getDeliveryNotes(cfg, cfg.GetDeliveryNotesEndpoint())
}

func (s *Session) getDeliveryNotes(cfg Config, endpoint string) (*DeliveryNotes, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	_, content, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	var notes DeliveryNotes
	if err := json.Unmarshal(content, &notes); err != nil {
		return nil, fmt.Errorf("could not load json response due to %s", err)
	}

	if notes.NextLink != nil && *notes.NextLink != "" {
		next, err := s.getDeliveryNotes(cfg, cfg.BuildEndpoint(*notes.NextLink))
		if err != nil {
			return &notes, err
		}

		notes.Value = append(notes.Value, next.Value...)
	}

	return &notes, nil
}

func (s *Session) GetDeliveryNote(cfg Config, id string) (*DeliveryNote, error) {
	req, err := http.NewRequest(http.MethodGet, cfg.GetDeliveryNoteEndpoint(id), nil)
	if err != nil {
		return nil, err
	}

	_, content, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	var note DeliveryNote
	if err := json.Unmarshal(content, &note); err != nil {
		return nil, fmt.Errorf("could not load json response due to %s", err)
	}

	return &note, nil
}

func (s *Session) changeDeliveryNote(endpoint string) error {
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}

	_, _, err = s.Do(req)

	return err
}

func (s *Session) RopenDeliveryNote(cfg Config, id string) error {
	return s.changeDeliveryNote(cfg.ReopenDeliveryNoteEndpoint(id))
}

func (s *Session) CloseDeliveryNote(cfg Config, id string) error {
	return s.changeDeliveryNote(cfg.CloseDeliveryNoteEndpoint(id))
}

func (s *Session) CancelDeliveryNote(cfg Config, id string) error {
	return s.changeDeliveryNote(cfg.CancelDeliveryNoteEndpoint(id))
}

func (s *Session) GetPurchaseOrders(cfg Config) (*PurchaseOrders, error) {
	return s.getPurchaseOrders(cfg, cfg.GetPurchaseOrdersEndpoint())
}

func (s *Session) getPurchaseOrders(cfg Config, endpoint string) (*PurchaseOrders, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	_, content, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	var notes PurchaseOrders
	if err := json.Unmarshal(content, &notes); err != nil {
		return nil, fmt.Errorf("could not load json response due to %s", err)
	}

	if notes.NextLink != nil && *notes.NextLink != "" {
		next, err := s.getPurchaseOrders(cfg, cfg.BuildEndpoint(*notes.NextLink))
		if err != nil {
			return &notes, err
		}

		notes.Value = append(notes.Value, next.Value...)
	}

	return &notes, nil
}

func (s *Session) GetPurchaseOrder(cfg Config, id string) (*PurchaseOrder, error) {
	req, err := http.NewRequest(http.MethodGet, cfg.GetPurchaseOrderEndpoint(id), nil)
	if err != nil {
		return nil, err
	}

	_, content, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	var note PurchaseOrder
	if err := json.Unmarshal(content, &note); err != nil {
		return nil, fmt.Errorf("could not load json response due to %s", err)
	}

	return &note, nil
}

func (s *Session) ReopenPurchaseOrder(cfg Config, id string) error {
	return s.changeDeliveryNote(cfg.ReopenPurchaseOrderEndpoint(id))
}

func (s *Session) ClosePurchaseOrder(cfg Config, id string) error {
	return s.changeDeliveryNote(cfg.ClosePurchaseOrderEndpoint(id))
}

func (s *Session) CancelPurchaseOrder(cfg Config, id string) error {
	return s.changeDeliveryNote(cfg.CancelPurchaseOrderEndpoint(id))
}

func (s *Session) GetPurchaseDeliveryNotes(cfg Config) (*PurchaseDeliveryNotes, error) {
	return s.getPurchaseDeliveryNotes(cfg, cfg.GetPurchaseDeliveryNotesEndpoint())
}

func (s *Session) getPurchaseDeliveryNotes(cfg Config, endpoint string) (*PurchaseDeliveryNotes, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	_, content, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	var notes PurchaseDeliveryNotes
	if err := json.Unmarshal(content, &notes); err != nil {
		return nil, fmt.Errorf("could not load json response due to %s", err)
	}

	if notes.NextLink != nil && *notes.NextLink != "" {
		next, err := s.getPurchaseDeliveryNotes(cfg, cfg.BuildEndpoint(*notes.NextLink))
		if err != nil {
			return &notes, err
		}

		notes.Value = append(notes.Value, next.Value...)
	}

	return &notes, nil
}

func (s *Session) GetPurchaseDeliveryNote(cfg Config, id string) (*PurchaseDeliveryNote, error) {
	return retrieveDocument[PurchaseDeliveryNote](s, cfg.GetPurchaseDeliveryNoteEndpoint(id))
}

func (s *Session) ReopenPurchaseDeliveryNote(cfg Config, id string) error {
	return s.changeDeliveryNote(cfg.ReopenPurchaseDeliveryNoteEndpoint(id))
}

func (s *Session) ClosePurchaseDeliveryNote(cfg Config, id string) error {
	return s.changeDeliveryNote(cfg.ClosePurchaseDeliveryNoteEndpoint(id))
}

func (s *Session) CancelPurchaseDeliveryNote(cfg Config, id string) error {
	return s.changeDeliveryNote(cfg.CancelPurchaseDeliveryNoteEndpoint(id))
}

func (s *Session) CreatePurchaseDeliveryNote(cfg Config, note PurchaseDeliveryNote) (bool, error) {
	payload, err := json.Marshal(note)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodPost,
		cfg.GetPurchaseDeliveryNotesEndpoint(), strings.NewReader(string(payload)))
	if err != nil {
		return false, err
	}

	_, _, err = s.Do(req)
	if err != nil {
		return false, err
	}

	return true, nil
}

// retrieveDocument pulls a document type from an SAP endpoint. The type of Unmarshal needs to
// be specified when calling the function.
func retrieveDocument[T any](s *Session, endpoint string) (*T, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	_, content, err := s.Do(req)
	if err != nil {
		return nil, err
	}

	var doc T
	if err := json.Unmarshal(content, &doc); err != nil {
		return nil, fmt.Errorf("could not load json response due to %s", err)
	}

	return &doc, nil
}

func (s *Session) GetInventoryCounting(cfg Config, id int) (*InventoryCounting, error) {
	url := cfg.GetInventoryCountingEndpoint(id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var counting InventoryCounting
	if err := json.NewDecoder(resp.Body).Decode(&counting); err != nil {
		return nil, fmt.Errorf("could not decode response body due to %s", err)
	}

	return &counting, nil
}

func (s *Session) GetInventoryCountings(cfg Config, filter string) ([]InventoryCounting, error) {
	url := cfg.GetInventoryCountingsEndpoint(filter)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var countings []InventoryCounting
	if err := json.NewDecoder(resp.Body).Decode(&countings); err != nil {
		return nil, fmt.Errorf("could not decode response body due to %s", err)
	}

	return countings, nil
}

func (s *Session) CreateInventoryCounting(cfg Config, counting InventoryCounting) (bool, error) {
	payload, err := json.Marshal(counting)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodPost, cfg.CreateInventoryCountingEndpoint(), strings.NewReader(string(payload)))
	if err != nil {
		return false, err
	}

	resp, err := s.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("could not read response body content due to %s", err)
	}

	return true, nil
}

func (s *Session) UpdateInventoryCounting(cfg Config, id int, updates InventoryCounting) (bool, error) {
	payload, err := json.Marshal(updates)
	if err != nil {
		return false, err
	}

	url := cfg.GetInventoryCountingEndpoint(id)
	req, err := http.NewRequest(http.MethodPatch, url, strings.NewReader(string(payload)))
	if err != nil {
		return false, err
	}

	resp, err := s.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("could not read response body content due to %s", err)
	}

	return true, nil
}

func (s *Session) DeleteInventoryCounting(cfg Config, id int) (bool, error) {
	url := cfg.GetInventoryCountingEndpoint(id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := s.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("could not read response body content due to %s", err)
	}

	return true, nil
}

func (s *Session) CloseInventoryCounting(cfg Config, id int) (bool, error) {
	url := cfg.CloseInventoryCountingEndpoint(id)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := s.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("could not read response body content due to %s", err)
	}

	return true, nil
}

func (s *Session) AddLinesToInventoryCounting(cfg Config, id int, lines []InventoryCountingLine) (bool, error) {
	// Retrieve the existing inventory counting
	existingCounting, err := s.GetInventoryCounting(cfg, id)
	if err != nil {
		return false, err
	}

	// Append new lines to the existing lines
	existingCounting.InventoryCountingLines = append(existingCounting.InventoryCountingLines, lines...)

	// Prepare the payload for the update
	payload, err := json.Marshal(existingCounting)
	if err != nil {
		return false, err
	}

	// Send the PATCH request to update the inventory counting
	url := cfg.GetInventoryCountingEndpoint(id)
	req, err := http.NewRequest(http.MethodPatch, url, strings.NewReader(string(payload)))
	if err != nil {
		return false, err
	}

	resp, err := s.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("could not read response body content due to %s", err)
	}

	return true, nil
}
