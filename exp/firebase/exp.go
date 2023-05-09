package main

import (
	"context"
	"fmt"
	"os"

	"firebase.google.com/go/v4/auth"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// TODO(lamuguo): Remove the file in the future.
func main() {
	// Set your Firebase local emulator URL.
	os.Setenv("FIREBASE_AUTH_EMULATOR_HOST", "localhost:9099")

	ctx := context.Background()
	sa := option.WithCredentialsFile("testing-firebase-secret.json")

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
		return
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		fmt.Printf("error initializing Auth client: %v\n", err)
		return
	}

	//// Disable the use of the token source while using the Auth Emulator
	//authClient.Opts.TokenSource = nil
	//
	params := (&auth.UserToCreate{}).
		Email("user@atomi.ai").
		Password("YourPassword123!")

	user, err := authClient.CreateUser(ctx, params)
	if err != nil {
		fmt.Printf("error creating user: %v\n", err)
		return
	}

	fmt.Printf("Successfully created user: %v\n", user)
}
