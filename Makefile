.PHONY: all
all: zk

zk: main.go command.go go.mod go.sum map.go set.go doc.go
	go build .

.PHONY: install
install: zk
	@install -vm0755 zk /usr/bin/zk

uninstall:
	@rm -vrf /usr/bin/zk
