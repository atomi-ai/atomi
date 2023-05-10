package tests

import "github.com/atomi-ai/atomi/app"

func Setup(dbName string) (*app.Application, error) {
	app, err := app.InitializeTestingApplication(dbName)
	if err != nil {
		return nil, err
	}
	return app, nil
}
