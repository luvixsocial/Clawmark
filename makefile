all:
	CGO_ENABLED=0 go build -v
start:
	./clawmark
clean:
	go fmt ./...