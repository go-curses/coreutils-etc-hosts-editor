#!/usr/bin/make --no-print-directory --jobs=1 --environment-overrides -f

# Copyright (c) 2023  The Go-Curses Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#: uncomment to echo instead of execute
#CMD=echo

-include .env
#export

.PHONY: run

BIN_NAME ?= eheditor
UNTAGGED_VERSION ?= v0.1.0
UNTAGGED_COMMIT ?= 0000000000

SHELL := /bin/bash
RUN_ARGS ?= ./example.etc-hosts
LOG_LEVEL := debug

GO_ENJIN_PKG := nil
BE_LOCAL_PATH := nil

GOPKG_KEYS ?= CDK CTK

CDK_GO_PACKAGE ?= github.com/go-curses/cdk
CDK_LOCAL_PATH ?= ../cdk

CTK_GO_PACKAGE ?= github.com/go-curses/ctk
CTK_LOCAL_PATH ?= ../ctk

CLEAN_FILES     ?= ${BIN_NAME} ${BIN_NAME}.*.* coverage.out pprof.*
DISTCLEAN_FILES ?=
REALCLEAN_FILES ?=

BUILD_VERSION_VAR := main.BuildVersion
BUILD_RELEASE_VAR := main.BuildRelease

SRC_CMD_PATH := ./cmd/eheditor

define help_custom_targets
	@echo "  run         - run the dev build (sanely handle crashes)"
	@echo "  profile.cpu - run the dev build and profile CPU"
	@echo "  profile.mem - run the dev build and profile memory"
endef

include Golang.cmd.mk
include Golang.def.mk

run: export GO_CDK_LOG_FILE=./${BUILD_NAME}.cdk.log
run: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
run: export GO_CDK_LOG_FULL_PATHS=true
run:
	@if [ -f ${BUILD_NAME} ]; \
	then \
		echo "# running: ${BUILD_NAME} ${RUN_ARGS}"; \
		( ./${BUILD_NAME} ${RUN_ARGS} ) 2>> ${GO_CDK_LOG_FILE}; \
		if [ $$? -ne 0 ]; \
		then \
			stty sane; echo ""; \
			echo "# ${BIN_NAME} crashed, see: ./${BIN_NAME}.cdk.log"; \
			read -p "# Press <Enter> to reset terminal, <Ctrl+C> to cancel" RESP; \
			reset; \
			echo "# ${BIN_NAME} crashed, terminal reset, see: ./${BIN_NAME}.cdk.log"; \
		else \
			echo "# ${BIN_NAME} exited normally."; \
		fi; \
	fi

profile.cpu: export GO_CDK_LOG_FILE=./${BIN_NAME}.cdk.log
profile.cpu: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
profile.cpu: export GO_CDK_LOG_FULL_PATHS=true
profile.cpu: export GO_CDK_PROFILE_PATH=/tmp/${BIN_NAME}.cdk.pprof
profile.cpu: export GO_CDK_PROFILE=cpu
profile.cpu: debug
	@mkdir -v /tmp/${BIN_NAME}.cdk.pprof 2>/dev/null || true
	@if [ -f ${BIN_NAME} ]; \
		then \
			./${BIN_NAME} && \
			if [ -f /tmp/${BIN_NAME}.cdk.pprof/cpu.pprof ]; \
			then \
				read -p "# Press enter to open a pprof instance" JUNK \
				&& go tool pprof -http=:8080 /tmp/${BIN_NAME}.cdk.pprof/cpu.pprof ; \
			else \
				echo "# missing /tmp/${BIN_NAME}.cdk.pprof/cpu.pprof"; \
			fi ; \
		fi

profile.mem: export GO_CDK_LOG_FILE=./${BIN_NAME}.log
profile.mem: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
profile.mem: export GO_CDK_LOG_FULL_PATHS=true
profile.mem: export GO_CDK_PROFILE_PATH=/tmp/${BIN_NAME}.cdk.pprof
profile.mem: export GO_CDK_PROFILE=mem
profile.mem: debug
	@mkdir -v /tmp/${BIN_NAME}.cdk.pprof 2>/dev/null || true
	@if [ -f ${BIN_NAME} ]; \
		then \
			./${BIN_NAME} && \
			if [ -f /tmp/${BIN_NAME}.cdk.pprof/mem.pprof ]; \
			then \
				read -p "# Press enter to open a pprof instance" JUNK \
				&& go tool pprof -http=:8080 /tmp/${BIN_NAME}.cdk.pprof/mem.pprof; \
			else \
				echo "# missing /tmp/${BIN_NAME}.cdk.pprof/mem.pprof"; \
			fi ; \
		fi
