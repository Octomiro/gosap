package gosap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/spf13/viper"
)

const B1DeaultPort = 50000

type Config struct {
	IP        string `mapstructure:"IP"`
	Port      uint16 `mapstructure:"PORT"`
	CompanyDB string `mapstructure:"COMPANY_DB"`
	Username  string `mapstructure:"DB_USERNAME"`
	Password  string `mapstructure:"DB_PASSWORD"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("gosap")
	viper.SetConfigType("env")

	viper.SetDefault("IP", "")
	viper.SetDefault("PORT", B1DeaultPort)
	viper.SetDefault("COMPANY_DB", "")
	viper.SetDefault("USERNAME", "")
	viper.SetDefault("PASSWORD", "")

	viper.AutomaticEnv()

	config := Config{}

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, err
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil
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

func (c *Config) GetItemEndpoint(id string) string {
	return fmt.Sprintf("https://%s/b1s/v1/Items('%s')?$select=ItemCode,ItemName,PurchaseUnitWidth",
		c.hostPort(), id)
}

func (c *Config) GetSuppliersEndpoint() string {
	return fmt.Sprintf("https://%s/b1s/v1/BusinessPartners?$select=CardCode,CardName&$filter=CardType eq 'S'",
		c.hostPort())
}

func (c *Config) GetClientsEndpoint() string {
	return fmt.Sprintf("https://%s/b1s/v1/BusinessPartners?$select=CardCode,CardName&$filter=CardType eq 'C'",
		c.hostPort())
}

func (c *Config) GetDeliveryNoteEndpoint(id string) string {
	return fmt.Sprintf("https://%s/b1s/v1/DeliveryNotes(%s)", c.hostPort(), id)
}

func (c *Config) GetDeliveryNotesEndpoint() string {
	return fmt.Sprintf("https://%s/b1s/v1/DeliveryNotes", c.hostPort())
}

func (c *Config) CreatePurchaseDeliveryNoteEndpoint() string {
	return fmt.Sprintf("https://%s/b1s/v1/PurchaseDeliveryNotes", c.hostPort())
}

func (c *Config) hostPort() string {
	return net.JoinHostPort(c.IP, strconv.Itoa(int(c.Port)))
}
