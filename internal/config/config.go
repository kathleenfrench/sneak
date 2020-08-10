package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Settings represent the sneak configs
type Settings struct {
	viper           *viper.Viper
	OpenVPNFilepath string `yaml:"openvpn_filepath"`
	HTBUsername     string `yaml:"htb_username"`
	DefaultEditor   string `yaml:"default_editor"`
	DBPath          string `yaml:"data"`
	Home            string
	WebShortcuts    map[string]string `yaml:"webshort"`
	HTBNetworkIP    string            `yaml:"htb_network_ip"` // https://www.hackthebox.eu/home/htb/access
}

// Parse unmarshals the viper configs into the sneak settings struct
func Parse(v *viper.Viper) (*Settings, error) {
	cfg.viper = v
	dco := decodeWithYaml("yaml")
	err := v.Unmarshal(cfg, dco)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

// GetSettings returns the settings
func GetSettings() *Settings {
	return cfg
}

// ParseAndUpdate parses the viper settings as a sneak settings struct and updates the config file
func ParseAndUpdate(v *viper.Viper) error {
	s, err := Parse(v)
	if err != nil {
		return err
	}

	return s.UpdateSettings()
}

func decodeWithYaml(tagName string) viper.DecoderConfigOption {
	return func(c *mapstructure.DecoderConfig) {
		c.TagName = tagName
	}
}

// Get fetches a config value by key
func Get(key string) interface{} {
	return viper.Get(key)
}

// Set sets a config key and value and saves it to the config file
func Set(key string, value string) {
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		panic(err)
	}
}
