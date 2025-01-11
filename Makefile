GOARCH=arm
GOARM=7
GOOS=linux
KINDLE_SSH_ADDRESS=root@kindle

.PHONY: all dev install run run-background

all: kindle-http-server

dev:
	echo kindle-http-server.go | entr -r make -f ./Makefile all install run

install: kindle-http-server
	ssh $(KINDLE_SSH_ADDRESS) "killall -9 kindle-http-server" || true
	scp kindle-http-server $(KINDLE_SSH_ADDRESS):/tmp/

run:
	ssh $(KINDLE_SSH_ADDRESS) "killall -9 kindle-http-server; /tmp/kindle-http-server"

run-background:
	ssh $(KINDLE_SSH_ADDRESS) "killall -9 kindle-http-server; nohup /tmp/kindle-http-server < /dev/null > /dev/null 2>&1 &"

kindle-http-server: kindle-http-server.go
	GOARCH=$(GOARCH) GOARM=$(GOARM) GOOS=$(GOOS) go build -o kindle-http-server kindle-http-server.go
