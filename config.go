package gosap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type Config struct {
	IP        string
	Port      uint16
	CompanyDB string
	Username  string
	Password  string
}

func (c *Config) LoginEndpoint() string {
	return fmt.Sprintf("https://%s/b1s/v1/Login", net.JoinHostPort(c.IP, strconv.Itoa(int(c.Port))))
}

func (c *Config) LoginPayload() (string, error) {
	res, err := json.Marshal(map[string]string{
		"CompanyDB": c.CompanyDB,
		"Password":  c.Password,
		"UserName":  c.Username,
	})
	if err != nil {
		return "", errors.New("Could not create payload")
	}

	return string(res), nil
}

func (c *Config) GetItemsEndpoint() string {
	return fmt.Sprintf("https://%s/b1s/v1/Items?$select=ItemCode,ItemName,PurchaseUnitWidth",
		c.hostPort())
}

func (c *Config) GetSuppliersEndpoint() string {
	return fmt.Sprintf("https://%s/b1s/v1/BusinessPartners?$select=CardCode,CardName&$filter=CardType eq 'S'",
		c.hostPort())
}

func (c *Config) GetClientsEndpoint() string {
	return fmt.Sprintf("https://%s/b1s/v1/BusinessPartners?$select=CardCode,CardName&$filter=CardType eq 'C'",
		c.hostPort())
}

func (c *Config) CreatePurchaseDeliveryNoteEndpoint() string {
	return fmt.Sprintf("https://%s/b1s/v1/PurchaseDeliveryNotes", c.hostPort())
}

func (c *Config) hostPort() string {
	return net.JoinHostPort(c.IP, strconv.Itoa(int(c.Port)))
}
