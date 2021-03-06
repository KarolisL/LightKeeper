PROGRAM := lightkeeper
BIN_DIR := bin
IPK_DIR := build/ipkbuild/lightkeeper

opkg-omnia: build-linux-armv7l copy-omnia
	./scripts/package-opkg.sh

copy-omnia:
	mkdir -p $(IPK_DIR)/data/usr/bin/
	mkdir -p $(IPK_DIR)/data/etc/init.d
	cp init/procd/lightkeeper $(IPK_DIR)/data/etc/init.d
	chmod +x $(IPK_DIR)/data/etc/init.d/lightkeeper
	cp $(BIN_DIR)/$(PROGRAM)-linux-armv7l $(IPK_DIR)/data/usr/bin/$(PROGRAM)

build-linux-armv7l:
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARM=7 GOARCH=arm go build -o $(BIN_DIR)/$(PROGRAM)-linux-armv7l ./cmd/$(PROGRAM)

