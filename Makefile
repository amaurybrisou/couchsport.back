ifndef GOROOT 
$(error "GOROOT is not set")
endif

ifndef GOBIN
$(error "GOBIN is not set")
endif

PROJECTNAME=$(shell basename "$(PWD)")
PUBLIC=$(CURDIR)/public
RELEASE_PATH=$(CURDIR)/release

GOCMD=$(GOROOT)/bin/go
NPM=$(shell which npm)
GIT=$(shell which git)

.DEFAULT_GOAL := dev

#make execute dev by default
dev: config.dev.json server.PID client.PID
stop: stop-server stop-client

stop-client: client.PID
		kill `cat $<` && rm $<

stop-server: server.PID
		kill `cat $<` && rm $<

server.PID:
		cd $(CURDIR) && { $(GOBIN)/$(PROJECTNAME) --env=dev & echo $$! > $@; }

client.PID: 
		cd $(PUBLIC) && $(NPM) run dev

build: build_front build_back

release: build
	[ -d $(RELEASE_PATH) ] && rm -rf $(RELEASE_PATH) && \
	[ -d $(RELEASE_PATH) ] || mkdir $(RELEASE_PATH) && \
	rm -rf .git/worktree/deploy && \
	git worktree prune && \
	git worktree add -B deploy $(RELEASE_PATH) puzzledge/deploy && \
	rm -rf $(RELEASE_PATH)/* 2> /dev/null && \
	rm -rf $(RELEASE_PATH)/.gitignore && \
	cp -rf $(GOBIN)/$(PROJECTNAME) $(RELEASE_PATH) && \
	cp -rf $(PUBLIC)/dist $(RELEASE_PATH) && \
	cp -rf $(CURDIR)/config.dev.json.default $(RELEASE_PATH) && \
	mkdir $(RELEASE_PATH)/localizer && \
	cp -rf $(CURDIR)/localizer/*.json $(RELEASE_PATH)/localizer/ && \
	cd $(RELEASE_PATH) && \
	git add --all -f 2> /dev/null && \
	git commit -m "new release" -q 2> /dev/null && \
	git push puzzledge deploy 2> /dev/null && \
	rm -rf $(RELEASE_PATH) 2> /dev/null

build_front:
		cd $(PUBLIC) && $(NPM) run build && cd $(CURDIR)

build_back:
		GO111MODULE=on GOARCH=amd64 $(GOCMD) build -o $(GOBIN)/$(PROJECTNAME)

.PHONY: dev stop build_back build_front client.PID config.dev.json server.PID



