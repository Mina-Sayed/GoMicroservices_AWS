package app

import (
	"lambda/api"
	"lambda/database"
)

type App struct {
	AppHandler api.ApiHandler
}

func NewApp() App {

	// we actually  initialize our database store
	// gets passed DOWN Into the ApiHandler

	db := database.NewDynamoDBClient()
	apiHandler := api.NewApiHandler(db)

	return App{
		AppHandler: apiHandler,
	}
}
