package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Settings represent the sneak configs
type Settings struct {
	viper           *viper.Viper
	OpenVPNFilepath string            `yaml:"openvpn_filepath"`
	HTBUsername     string            `yaml:"htb_username"`
	BoxIPs          map[string]string `yaml:"box_ips"`
	DefaultEditor   string            `yaml:"default_editor"`
}

// Parse unmarshals the viper configs into the sneak settings struct
func Parse(v *viper.Viper) (Settings, error) {
	s := Settings{
		viper: v,
	}

	dco := decodeWithYaml("yaml")
	err := v.Unmarshal(&s, dco)
	if err != nil {
		return s, err
	}

	return s, nil
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
