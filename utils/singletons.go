package utils

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"os"
)

func FirebaseAppProvider() *firebase.App {
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
