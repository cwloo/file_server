SHELL := /usr/bin/bash

MAKEFLAGS += -rR --include-dir=$(CURDIR)/..

export COLOR_SUFFIX="\033[0m"
export BLACK_PREFIX="\033[30m"
export RED_PREFIX="\033[31m"
export GREEN_PREFIX="\033[32m"
export YELLOW_PREFIX="\033[33m"
export BLUE_PREFIX="\033[34m"
export PURPLE_PREFIX="\033[35m"
export SKY_BLUE_PREFIX="\033[36m"

export MAKE=make
export BUILD=build
export INSTALL=install

#include $(PWD)/src/loader/Makefile
#include $(PWD)/src/http_gate/Makefile
#include $(PWD)/src/file_server/Makefile
#include $(PWD)/src/file_client/Makefile

.PHONY: subsystem help

.shellmkdir:
	$(shell if [ ! -d $(BIN_DIR) ]; then mkdir -p $(BIN_DIR); fi)

.bashmkdir:
	@bash -c "if [ ! -d $(BIN_DIR) ]; then mkdir -p $(BIN_DIR); fi"

.mkdir:
	@if [ ! -d $(BIN_DIR) ]; then mkdir -p $(BIN_DIR); fi

.mkdirtip:
	@echo -e ${BLUE_PREFIX}"mkdir -p"${COLOR_SUFFIX} ${SKY_BLUE_PREFIX}${BIN_DIR}${COLOR_SUFFIX}

subsystem: .build

.cdtip:
	@echo -e ${BLUE_PREFIX}"cd"${COLOR_SUFFIX} ${PURPLE_PREFIX}$(MAKE_DIR)${COLOR_SUFFIX}

.buildtip:
	@echo -e ${BLUE_PREFIX}$(MAKE)${COLOR_SUFFIX} ${PURPLE_PREFIX}$(INSTALL)${COLOR_SUFFIX}

.build: .cdtip
	@bash -c "cd $(MAKE_DIR) && $(MAKE) $(INSTALL)"