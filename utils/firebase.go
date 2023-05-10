package utils

import (
	"context"
	"fmt"
	"os"

	"firebase.google.com/go/v4/auth"

	firebase "firebase.google.com/go/v4"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

type AuthAppWrapper interface {
	AuthAndDecode(ctx context.Context, token string) (*auth.Token, error)
}

type FirebaseAppWrapper struct {
	FirebaseApp *firebase.App
}

func NewFirebaseAppWrapper(firebaseApp *firebase.App) AuthAppWrapper {
	return &FirebaseAppWrapper{
		FirebaseApp: firebaseApp,
	}
}

func (w *FirebaseAppWrapper) AuthAndDecode(ctx context.Context, token string) (*auth.Token, error) {
	client, err := w.FirebaseApp.Auth(ctx)
	if err != nil {
		return nil, err
	}

	decodedToken, err := client.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return decodedToken, nil
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
