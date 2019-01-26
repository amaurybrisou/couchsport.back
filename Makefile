ifndef GOROOT 
$(error "GOROOT is not set")
endif

ifndef GOBIN
$(error "GOBIN is not set")
endif

ifndef ENV
$(error "ENV is not set")
endif

PROJECTNAME=$(shell basename "$(PWD)")
PUBLIC=$(CURDIR)/public

GOCMD=$(GOROOT)/bin/go
NPM=$(shell which npm)
GIT=$(shell which git)

.DEFAULT_GOAL := dev

#make execute dev by default
dev: config.dev.json build_back server.PID client.PID
stop: stop-server stop-client

stop-client: client.PID
	kill `cat $<` && rm $<

stop-server: server.PID
	kill `cat $<` && rm $<

server.PID:
		cd $(CURDIR) && { $(GOBIN)/$(PROJECTNAME) --env=$(ENV) & echo $$! > $@; }

client.PID: 
		cd $(PUBLIC) && $(NPM) run dev


clean:
	rm $(GOBIN)/$(PROJECTNAME)

start: config.prod.json build server.PID

build: pull build_front build_back

pull:
		$(GIT) pull

build_front:
		cd $(PUBLIC) && $(NPM) run build && cd $(CURDIR)

build_back:
		$(GOCMD) build -o $(GOBIN)/$(PROJECTNAME)

.PHONY: dev stop build_back build_front pull clean stop-prod start-prod