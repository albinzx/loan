package viper

import (
	"encoding/base64"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
}

// NewConfig creates new Config backed by viper file configuration
func New(path, file, ext string) (*Config, error) {
	v := &Config{}

	if err := v.init(path, file, ext); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Config) init(path, file, ext string) error {
	v.viper = viper.New()
	v.viper.SetConfigName(file)
	v.viper.SetConfigType(ext)
	v.viper.AddConfigPath(path)

	if err := v.viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func (v *Config) GetString(key string) string {
	return v.viper.GetString(key)
}

func (v *Config) GetBool(key string) bool {
	return v.viper.GetBool(key)
}

func (v *Config) GetInt(key string) int64 {
	return v.viper.GetInt64(key)
}

func (v *Config) GetFloat(key string) float64 {
	return v.viper.GetFloat64(key)
}

func (v *Config) GetBinary(key string) []byte {
	value := v.viper.GetString(key)
	bytes, err := base64.StdEncoding.DecodeString(value)
	if err == nil {
		return bytes
	}
	return nil
}

func (v *Config) GetArray(key string) []string {
	str := v.viper.GetString(key)
	return strings.Split(str, ",")
}
