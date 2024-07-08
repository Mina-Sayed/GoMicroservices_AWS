build:
	@echo "Building the Go Lambda function..."
	@GOOS=linux GOARCH=amd64 go build -o bootstrap
	@zip function.zip bootstrap

clean:
	@echo "Cleaning up..."
	rm -f bootstrap function.zip
