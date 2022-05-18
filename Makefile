CMD=handler

clean:
	rm -v bin/handler_*
build:
	go build -tags=jsoniter -o bin/$(CMD)_mac ./*.go
	GOOS=linux GOARCH=amd64 go build -o bin/$(CMD)_linux ./*.go
