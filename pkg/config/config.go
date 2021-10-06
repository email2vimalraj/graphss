package config

import "github.com/BurntSushi/toml"

// Cfg is the configuration object
type Cfg struct {
	// HTTP Server configs
	HTTP struct {
		Addr     string `toml:"addr"`
		Port     string `toml:"port"`
		CertFile string `toml:"certfile"`
		KeyFile  string `toml:"keyfile"`
		Protocol string `toml:"protocol"`
	} `toml:"http"`

	// Paths
	Paths struct {
		LogsDir string `toml:"logs-dir"`
	} `toml:"paths"`
}

func NewCfg(configFilePath string) (*Cfg, error) {
	cfg := &Cfg{}
	if err := cfg.loadConfiguration(configFilePath); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Cfg) loadConfiguration(configFilePath string) error {
	_, err := toml.DecodeFile(configFilePath, &c)
	if err != nil {
		return err
	}
	return nil
}
