package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v74"
)

func InitStripe(key string) {
	stripe.Key = key
}

func logAllSettings() {
	fmt.Println("All Viper settings:")
	for key, value := range viper.AllSettings() {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func LoadConfig() {
	configFile := os.Getenv("CONFIG_FILE")
	log.Printf("Load config from file: %v", configFile)
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	// The code below is debug only
	if viper.GetBool("logAllSettings") {
		logAllSettings()
	}
}
