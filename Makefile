build-osx-intel:
	GOOS=darwin GOARCH=amd64 go build -o builds/gql-osx-intel main.go

build-osx-arm:
	GOOS=darwin GOARCH=arm64 go build -o builds/gql-osx-arm main.go

build-linux-intel:
	GOOS=linux GOARCH=amd64 go build -o builds/gql-linux-intel main.go

build-linux-arm:
	GOOS=linux GOARCH=arm64 go build -o builds/gql-linux-arm main.go

build-windows-intel:
	GOOS=windows GOARCH=amd64 go build -o builds/gql-windows-intel main.go

build: build-osx-intel build-osx-arm build-linux-intel build-linux-arm build-windows-intel
