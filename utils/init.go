package utils

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v74"
	"google.golang.org/api/option"
	"log"
	"os"
)

func InitStripe(key string) {
	stripe.Key = key
}

func LoadConfig() {
	configFile := os.Getenv("CONFIG_FILE")
	log.Printf("Load config from file: %v", configFile)
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
}

func InitFirebase() *firebase.App {
	// Initialize Firebase app, set your Firebase local emulator URL for testing.
	if viper.GetBool("firebaseEnableEmulator") {
		os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", viper.GetString("firebaseAuthEmulatorHost"))
	}
	opt := option.WithCredentialsFile(viper.GetString("firebaseCredentialsFile"))
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println("error initializing firebase app:", err)
		os.Exit(1)
	}
	return firebaseApp
}
