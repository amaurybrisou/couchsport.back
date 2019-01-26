PROJECTNAME=$(shell basename "$(PWD)")
PUBLIC=$(CURDIR)/public

GOCMD=$(GOROOT)/bin/go
NPM=$(shell which npm)
GIT=$(shell which git)

.DEFAULT_GOAL := dev

#make execute dev by default
dev: config.dev.json server.PID client.PID  build_back

stop: server.PID client.PID
	kill `cat $+` && rm $+

server.PID:
		cd $(CURDIR) && { $(GOBIN)/$(PROJECTNAME) & echo $$! > $@; }

client.PID: 
		cd $(PUBLIC) && $(NPM) run dev


clean:
	rm $(GOBIN)/$(PROJECTNAME)

start-prod: config.prod.json build
	cd $(CURDIR) && { $(GOBIN)/$(PROJECTNAME) --env=prod & echo $! > server.PID; }

stop-prod: server.PID
	kill `cat $<` && rm $<

build: pull build_front build_back

pull:
		$(GIT) pull

build_front:
		cd $(PUBLIC) && $(NPM) run build && cd $(CURDIR)

build_back:
		$(GOCMD) build -o $(GOBIN)/$(PROJECTNAME)

.PHONY: dev stop build_back build_front pull clean stop-prod start-prod