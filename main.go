package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/warete/alice_kuzya_proxy/pkg/aqara"
)

type MainConfig struct {
	Aqara aqara.AqaraConfig `mapstructure:"aqara"`
}

func main() {

	viper.SetConfigFile("config.yml")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	var cfg MainConfig
	viper.Unmarshal(&cfg)

	aq, err := aqara.NewAqara(cfg.Aqara)
	if err != nil {
		log.Fatal(err.Error())
	}

	//TODO: make auth by secret key

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	aq.AddRoutes(r)

	r.Run("0.0.0.0:8080")
}
