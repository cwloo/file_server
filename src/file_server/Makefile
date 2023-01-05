SHELL := /usr/bin/bash

BINARY_NAME=$(shell basename $(PWD))

.PHONY: install run clean help

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"

run:
	$(shell cd $(BIN_DIR); ./$(BINARY_NAME))

gotool:
	go fmt ./
#	go vet ./

install: clean gotool build
	@echo -e ${BLUE_PREFIX}$(INSTALL)${COLOR_SUFFIX} ${GREEN_PREFIX}$(BIN_DIR)/$(BINARY_NAME)${COLOR_SUFFIX}
	$(MAKE) $(BUILD)
	@echo -e ${BLUE_PREFIX}"mv"${COLOR_SUFFIX} ${RED_PREFIX}${BINARY_NAME}${COLOR_SUFFIX} ${SKY_BLUE_PREFIX}">>>"${COLOR_SUFFIX} ${PURPLE_PREFIX}${BIN_DIR}${COLOR_SUFFIX}
	@mv ${BINARY_NAME} ${BIN_DIR}

clean:
	@echo -e ${BLUE_PREFIX}"clean"${COLOR_SUFFIX} ${GREEN_PREFIX}$(BIN_DIR)/$(BINARY_NAME)${COLOR_SUFFIX}
	@bash -c "cd $(BIN_DIR); if [ -f $(BINARY_NAME) ] ; then rm $(BINARY_NAME) ; fi"