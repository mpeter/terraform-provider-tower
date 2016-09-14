package tower

import (
	"net/http"

	"github.com/mpeter/go-towerapi/towerapi"
)

type Config struct {
	Endpoint string
	Username string
	Password string
}

func (c *Config) NewClient() (*towerapi.Client, error) {
	config := new(towerapi.ClientConfig)
	config.Endpoint = c.Endpoint
	config.Password = c.Password
	config.Username = c.Username

	return towerapi.NewClient(http.DefaultClient, config)
}
