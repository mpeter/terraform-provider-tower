package tower

import (
	"github.com/mpeter/go-towerapi/towerapi"
)

type Config struct {
	Endpoint string `mapstructure:"endpoint"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (c *Config) NewClient() (*towerapi.Client, error) {

	config := towerapi.DefaultConfig()

	if c.Endpoint != "" {
		config.Endpoint = c.Endpoint
	}
	if c.Username != "" {
		config.Password = c.Password
	}
	if c.Username != "" {
		config.Username = c.Username
	}

	if err := config.LoadAndValidate() ; err != nil {
		return nil, err
	}

	return towerapi.NewClient(config)
}
