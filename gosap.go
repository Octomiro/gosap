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

	statusOK := resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
	if !statusOK {
		return nil, fmt.Errorf("request to SAP API (%s) was not successful due to %s", cfg.LoginEndpoint(), resp.Status)
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
func (s *Session) Do(req *http.Request) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	s.setSessionCookies(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	statusOK := resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
	if !statusOK {
		return nil, fmt.Errorf("request to SAP API (%s) was not successful due to %s", req.URL, resp.Status)
	}

	return resp, nil
}

func (s *Session) GetItem(cfg Config, id string) (*Item, error) {
	req, err := http.NewRequest(http.MethodGet, cfg.GetItemEndpoint(id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body content due to %s", err)
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

	resp, err := s.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body content due to %s", err)
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

	resp, err := s.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body content due to %s", err)
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
	req, err := http.NewRequest(http.MethodGet, cfg.GetClientsEndpoint(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body content due to %s", err)
	}

	var clients Clients
	if err := json.Unmarshal(content, &clients); err != nil {
		return nil, fmt.Errorf("could not load json response due to %s", err)
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

	resp, err := s.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body content due to %s", err)
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

	resp, err := s.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body content due to %s", err)
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

	_, err = s.Do(req)

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

func (s *Session) CreatePurchaseDeliveryNote(cfg Config, note PurchaseDeliveryNotes) (bool, error) {
	payload, err := json.Marshal(note)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(http.MethodPost,
		cfg.CreatePurchaseDeliveryNoteEndpoint(), strings.NewReader(string(payload)))
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
