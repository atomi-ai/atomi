package utils

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

type FirebaseAppWrapper interface {
	Auth(ctx context.Context) (*auth.Client, error)
}

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
