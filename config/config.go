package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path"
)

type ConvertParams struct {
	Delimiter    string `mapstructure:"delimiter"`
	Url          string `mapstructure:"url"`
	CheckUrl     bool   `mapstructure:"check_url"`
	Delivery     bool   `mapstructure:"delivery"`
	DeliveryTime string `mapstructure:"delivery_time"`
	DeliveryCost int    `mapstructure:"delivery_cost"`
	DeliveryOnly bool   `mapstructure:"delivery_only"`
	NeedOrder    bool   `mapstructure:"need_order"`
	Currency     string `mapstructure:"currency"`
	Warranty     bool   `mapstructure:"warranty"`
	WixUrl       string `mapstructure:"wix_url"`
}

type Config struct {
	ConvertParams ConvertParams `mapstructure:"convert-params"`
}

var Cfg = Config{
	ConvertParams: ConvertParams{
		Delimiter:    ";",
		Url:          "https://магазин",
		CheckUrl:     true,
		WixUrl:       "https://static.wixstatic.com/media/",
		Delivery:     false,
		DeliveryTime: "",
		DeliveryCost: 0,
		DeliveryOnly: false,
		NeedOrder:    false,
		Currency:     "",
		Warranty:     false,
	},
}

func InitConfig(configFile string) {
	viper.SetConfigName(path.Base(configFile))
	viper.SetConfigType("toml")
	viper.AddConfigPath(path.Dir(configFile))

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Ошибка при инициализации настроек: %s \n", err))
	}

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		panic(fmt.Errorf("ошибка при чтении настроек: %s \n", err))
	}
}
