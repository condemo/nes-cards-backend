binary-name=cards-backend

build:
	@GOOS=windows GOARCH=amd64 go build -o ./bin/${binary-name}-win.exe ./cmd/main.go
	@GOOS=linux GOARCH=amd64 go build -o ./bin/${binary-name}-linux ./cmd/main.go
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/${binary-name}-darwin ./cmd/main.go

run: build
	@./bin/${binary-name}-linux

arm-build:
	@GOOS=linux GOARCH=arm64 go build -o ./bin/${binary-name}-arm64 ./cmd/main.go

arm-run: arm-build
	@./bin/${binary-name}-arm64

kill-services:
	@lsof -t -i:3000 | xargs -r kill

test:
	@go test ./cmd/main.go

clean:
	@rm -rf ./bin/*
	@go clean


