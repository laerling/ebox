.PHONY: goimports gopath install run clean uninstall purge test

EXE=ebox
ARGS=


# building and installing stuff

$(GOPATH)/bin/$(EXE): $(EXE)
	go install

$(EXE): goimports $(wildcard *.go)
	go build -o "$(EXE)"

install: "/bin/$(EXE)"

"/bin/$(EXE)":
	sudo -E install -m 0755 "$(GOPATH)/bin/$(EXE)" $@


# running stuff

run: $(GOPATH)/bin/$(EXE)
	$(GOPATH)/bin/$(EXE) $(ARGS)

test: goimports
	go test


# removing stuff

clean:
	rm -f "$(EXE)"

uninstall:
	sudo -E rm -f "/bin/$(EXE)"
	# gopath is no dep of this target because empty gopath is effectively the same as /bin/$(EXE)
	sudo -E rm -f "$(GOPATH)/bin/$(EXE)"

purge: uninstall clean


# depending on stuff

gopath:
	if [ -z "$(GOPATH)" ]; then exit 1; fi


# other stuff

goimports: $(wildcard *.go)
	goimports -w $+
