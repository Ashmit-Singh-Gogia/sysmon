# Variables (makes it easy to change the name later if you want)
APP_NAME=sysmon

# Formats all your Go code to look clean and standard
fmt:
	go fmt ./...

# Runs the application without compiling a final file (great for quick testing)
run:
	go run cmd/sysmon/main.go

# Compiles the application into a finished, runnable executable
build:
	go build -o $(APP_NAME) cmd/sysmon/main.go

# Runs any automated tests you write later
test:
	go test ./... -v