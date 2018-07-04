.PHONY: gopath install run clean uninstall purge

EXE=ebox
ARGS=

$(EXE): $(wildcard *.go)
	if [ -x /usr/bin/goimports ]; then goimports -w "$<"; fi
	go build -o "$(EXE)"
	go install

run: $(EXE)
	$(EXE) $(ARGS)

install: "/bin/$(EXE)"
"/bin/$(EXE)":
	sudo -E install -m 0755 "$(GOPATH)/bin/$(EXE)" $@

clean:
	rm -f "$(EXE)"

uninstall:
	sudo -E rm -f "/bin/$(EXE)"
	# gopath is no dep of this target because empty gopath is effectively the same as /bin/$(EXE)
	sudo -E rm -f "$(GOPATH)/bin/$(EXE)"

purge: uninstall clean

gopath:
	if [ -z "$(GOPATH)" ]; then exit 1; fi
