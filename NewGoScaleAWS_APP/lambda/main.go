package main

import (
	"fmt"
	"lambda/app"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

func HandleRequest(e MyEvent) (string, error) {

	if e.Username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	return fmt.Sprintf("succeesfuly called by - %s!", e.Username), nil
}

func main() {

	myApp := app.NewApp()

	lambda.Start(myApp.AppHandler.RegisterUser)
}
