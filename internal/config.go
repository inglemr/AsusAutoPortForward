package internal

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	RouterAddress        string `json:"router_address"`
	Username             string `json:"username"`
	Password             string `json:"password"`
	DefaultTargetAddress string `json:"default_target_address"`
}

func GetConfig() Config {
	config := Config{}
	log.SetLevel(log.DebugLevel)
	config.RouterAddress = os.Getenv("ROUTER_ADDRESS")
	config.Username = os.Getenv("ROUTER_USERNAME")
	config.Password = os.Getenv("ROUTER_PASSWORD")
	config.DefaultTargetAddress = os.Getenv("DEFAULT_TARGET_ADDRESS")
	return config
}
