package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/warete/alice_kuzya_proxy/pkg/aqara"
)

type MainConfig struct {
	Port      string            `mapstructure:"port"`
	SecretKey string            `mapstructure:"secret_key"`
	Aqara     aqara.AqaraConfig `mapstructure:"aqara"`
}

func KeyAuthMiddleware(trueKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := c.Query("secret")
		if len(secretKey) == 0 || secretKey != trueKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
			c.Abort()
			return
		}
		c.Next()
	}
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

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(KeyAuthMiddleware(cfg.SecretKey))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	aq.AddRoutes(r)

	r.Run("0.0.0.0:" + cfg.Port)
}
