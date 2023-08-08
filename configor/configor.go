package configor

import (
	"bytes"
	"github.com/spf13/viper"
)

type Configurator interface {
	LoadContents(any, []byte) error
}

func newViperWrapper() *viperWrapper {
	return &viperWrapper{}
}

type viperWrapper struct {
	*viper.Viper
}

func (v *viperWrapper) initWithJson(data []byte) error {
	v.Viper = viper.New()
	v.SetConfigType("json")
	return v.ReadConfig(bytes.NewBuffer(data))
}

func (v *viperWrapper) initWithToml(data []byte) error {
	v.Viper = viper.New()
	v.SetConfigType("toml")
	return v.ReadConfig(bytes.NewBuffer(data))
}
