package cfg

import (
	"bytes"
	"errors"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Name    string `toml:"Name"`
	Cluster string `toml:"ClusterURL"`
	ApiKey  string `toml:"-"`
	Default bool   `toml:"Default"`
}

type Profiles struct {
	Version string   `toml:"Version"`
	Configs []Config `toml:"Profile"`
}

func (p *Profiles) GetDefaultProfile() (*Config, error) {
	for _, c := range p.Configs {
		if c.Default {
			CONFIG = &c
			return &c, nil
		}
	}
	return nil, errors.New("no default profile")
}

func (c *Config) String() string {
	buf := new(bytes.Buffer)
	toml.NewEncoder(buf).Encode(c)
	return buf.String()
}

func (p *Profiles) String() string {
	buf := new(bytes.Buffer)
	toml.NewEncoder(buf).Encode(p)
	return buf.String()
}

func GetSampleProfile() *Profiles {
	return &Profiles{
		Version: "V1",
		Configs: []Config{
			{
				Name:    "<Enter your profile name here>",
				Cluster: "<Enter your cluster URL here>",
				ApiKey:  "<Enter your api key here>",
				Default: true,
			},
		},
	}
}
