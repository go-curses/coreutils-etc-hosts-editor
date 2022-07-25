#!/usr/bin/make -f

SHELL := /bin/bash
BUILD_CMD := eheditor
CDK_PATH := ../../cdk
CTK_PATH := ../../ctk
ETC_HOSTS_PATH := ./example.etc-hosts
LOG_LEVEL := debug

.PHONY: all build clean clean-logs dev fmt help profile.cpu profile.mem run tidy

all: help

help:
	@echo "usage: make [target]"
	@echo
	@echo "qa targets:"
	@echo "  vet         - run go vet command"
	@echo "  test        - perform all available tests"
	@echo "  cover       - perform all available tests with coverage report"
	@echo
	@echo "cleanup targets:"
	@echo "  clean       - cleans package and built files"
	@echo "  clean-logs  - cleans *.log from the project"
	@echo
	@echo "go.mod helpers:"
	@echo "  local       - add go.mod local CDK package replacements"
	@echo "  unlocal     - remove go.mod local CDK package replacements"
	@echo
	@echo "build targets:"
	@echo "  deps        - install stringer and bitmasker tools"
	@echo "  generate    - run go generate"
	@echo "  build       - build the ${BUILD_CMD} command"
	@echo "  dev         - build the ${BUILD_CMD} with profiling"
	@echo
	@echo "run targets:"
	@echo "  run         - run the dev build (sanely handle crashes)"
	@echo "  profile.cpu - run the dev build and profile CPU"
	@echo "  profile.mem - run the dev build and profile memory"

vet:
	@echo -n "# vetting etc-hosts-editor ..."
	@go vet && echo " done"

test: vet
	@echo "# testing etc-hosts-editor ..."
	@go test -v ./...

cover:
	@echo "# testing etc-hosts-editor (with coverage) ..."
	@go test -cover -coverprofile=coverage.out ./...
	@echo "# test coverage ..."
	@go tool cover -html=coverage.out

clean-build-logs:
	@echo "# cleaning *.build.log files"
	@rm -fv *.build.log || true

clean-logs:
	@echo "# cleaning *.log files"
	@rm -fv *.log || true
	@echo "# cleaning *.out files"
	@rm -fv *.out || true
	@echo "# cleaning pprof files"
	@rm -rfv /tmp/*.cdk.pprof || true

clean-cmd: clean-build-logs
	@echo "# cleaning built commands"
	@for tgt in `ls cmd`; do \
		if [ -f $$tgt ]; then rm -fv $$tgt; fi; \
	done

clean: clean-logs clean-cmd
	@echo "# cleaning goland builds"
	@rm -rfv go_* || true

define _gocache_dir =
$(shell go env | grep ^GOCACHE | perl -pe 's#^GOCACHE=(.)(.+?)(\1)#\2#')
endef

clean-go-cache: GOLANG_CACHE_DIR=$(call _gocache_dir)
clean-go-cache:
	@echo "# cleaning go caches (${GOLANG_CACHE_DIR})"
	@echo "## before: $(shell du -hs ${CACHE_DIR} | awk '{print $$1}')"
	@go clean -r github.com/go-curses/coreutils/etc-hosts-editor
	@echo "## after: $(shell du -hs ${CACHE_DIR} | awk '{print $$1}')"

build:
	@echo -n "# building command ${BUILD_CMD}... "
	@cd cmd/${BUILD_CMD}; \
		( go build -v \
				-trimpath \
				-ldflags="\
-X 'main.IncludeProfiling=false' \
-X 'main.IncludeLogFile=false'   \
-X 'main.IncludeLogLevel=false'  \
" \
				-o ../../${BUILD_CMD} \
			2>&1 ) > ../../${BUILD_CMD}.build.log; \
		rv="$$?"; \
		cd - > /dev/null; \
		if [ $$rv = "0" -a -f ${BUILD_CMD} ]; then \
			echo "done."; \
		else \
			echo "failed.\n>\tsee ./${BUILD_CMD}.build.log for errors"; \
			false; \
		fi

deps:
	@echo "# installing dependencies..."

generate:
	@echo "# generate go sources..."
	@go generate -v ./...

depends-on-cdk-path:
	@if [ ! -d ${CDK_PATH} ]; then \
			echo "Error: $(MAKECMDGOALS) depends upon a valid CDK_PATH."; \
			echo "Default: ../cdk"; \
			echo ""; \
			echo "Specify the path to an existing CDK checkout with the"; \
			echo "CDK_PATH variable as follows:"; \
			echo ""; \
			echo " make CDK_PATH=../path/to/cdk $(MAKECMDGOALS)"; \
			echo ""; \
			false; \
		fi

local: depends-on-cdk-path
	@echo "# adding go.mod local CTK package replacements..."
	@go mod edit -replace=github.com/go-curses/ctk=${CTK_PATH}
	@echo "# adding go.mod local CDK package replacements..."
	@go mod edit -replace=github.com/go-curses/cdk=${CDK_PATH}
	@for tgt in charset encoding env log memphis; do \
		if [ -f ${CDK_PATH}/$$tgt/go.mod ]; then \
			echo "#\t$$tgt"; \
			go mod edit -replace=github.com/go-curses/cdk/$$tgt=${CDK_PATH}/$$tgt ; \
		fi; \
	done
	@for tgt in `ls ${CDK_PATH}/lib`; do \
		if [ -f ${CDK_PATH}/lib/$$tgt/go.mod ]; then \
			echo "#\tlib/$$tgt"; \
			go mod edit -replace=github.com/go-curses/cdk/lib/$$tgt=${CDK_PATH}/lib/$$tgt ; \
		fi; \
	done
	@echo "# running go mod tidy"
	@go mod tidy

unlocal: depends-on-cdk-path
	@echo "# removing go.mod local CTK package replacements..."
	@go mod edit -dropreplace=github.com/go-curses/ctk
	@echo "# removing go.mod local CDK package replacements..."
	@go mod edit -dropreplace=github.com/go-curses/cdk
	@for tgt in charset encoding env log memphis; do \
		if [ -f ${CDK_PATH}/$$tgt/go.mod ]; then \
			echo "#\t$$tgt"; \
			go mod edit -dropreplace=github.com/go-curses/cdk/$$tgt ; \
		fi; \
	done
	@for tgt in `ls ${CDK_PATH}/lib`; do \
		if [ -f ${CDK_PATH}/lib/$$tgt/go.mod ]; then \
			echo "#\tlib/$$tgt"; \
			go mod edit -dropreplace=github.com/go-curses/cdk/lib/$$tgt ; \
		fi; \
	done
	@echo "# running go mod tidy"
	@go mod tidy

dev:
	@if [ -d cmd/${BUILD_CMD} ]; \
	then \
		echo -n "# building: ${BUILD_CMD} [dev]... "; \
		cd cmd/${BUILD_CMD}; \
		( go build -v \
				-gcflags=all="-N -l" \
				-ldflags="\
-X 'main.IncludeProfiling=true' \
-X 'main.IncludeLogFile=true'   \
-X 'main.IncludeLogLevel=true'  \
" \
				-o ../../${BUILD_CMD} \
			2>&1 ) > ../../${BUILD_CMD}.build.log; \
		rv="$$?"; \
		cd - > /dev/null; \
		if [ $$rv = "0" -a -f ${BUILD_CMD} ]; then \
			echo "done."; \
		else \
			echo "failed.\n>\tsee ./${BUILD_CMD}.build.log errors below:"; \
			cat ./${BUILD_CMD}.build.log; \
			false; \
		fi; \
	else \
		echo "# build cmd not found: ${BUILD_CMD}"; \
	fi

run: export GO_CDK_LOG_FILE=./${BUILD_CMD}.cdk.log
run: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
run: export GO_CDK_LOG_FULL_PATHS=true
run:
	@if [ -f ${BUILD_CMD} ]; \
	then \
		echo "# running: ${BUILD_CMD} ${ETC_HOSTS_PATH}"; \
		( ./${BUILD_CMD} ${ETC_HOSTS_PATH} ) 2>> ${GO_CDK_LOG_FILE}; \
		if [ $$? -ne 0 ]; \
		then \
			stty sane; echo ""; \
			echo "# ${BUILD_CMD} crashed, see: ./${BUILD_CMD}.cdk.log"; \
			read -p "# Press <Enter> to reset terminal, <Ctrl+C> to cancel" RESP; \
			reset; \
			echo "# ${BUILD_CMD} crashed, terminal reset, see: ./${BUILD_CMD}.cdk.log"; \
		else \
			echo "# ${BUILD_CMD} exited normally."; \
		fi; \
	fi

profile.cpu: export GO_CDK_LOG_FILE=./${BUILD_CMD}.cdk.log
profile.cpu: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
profile.cpu: export GO_CDK_LOG_FULL_PATHS=true
profile.cpu: export GO_CDK_PROFILE_PATH=/tmp/${BUILD_CMD}.cdk.pprof
profile.cpu: export GO_CDK_PROFILE=cpu
profile.cpu: dev
	@mkdir -v /tmp/${BUILD_CMD}.cdk.pprof 2>/dev/null || true
	@if [ -f ${BUILD_CMD} ]; \
		then \
			./${BUILD_CMD} && \
			if [ -f /tmp/${BUILD_CMD}.cdk.pprof/cpu.pprof ]; \
			then \
				read -p "# Press enter to open a pprof instance" JUNK \
				&& go tool pprof -http=:8080 /tmp/${BUILD_CMD}.cdk.pprof/cpu.pprof ; \
			else \
				echo "# missing /tmp/${BUILD_CMD}.cdk.pprof/cpu.pprof"; \
			fi ; \
		fi

profile.mem: export GO_CDK_LOG_FILE=./${BUILD_CMD}.log
profile.mem: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
profile.mem: export GO_CDK_LOG_FULL_PATHS=true
profile.mem: export GO_CDK_PROFILE_PATH=/tmp/${BUILD_CMD}.cdk.pprof
profile.mem: export GO_CDK_PROFILE=mem
profile.mem: dev
	@mkdir -v /tmp/${BUILD_CMD}.cdk.pprof 2>/dev/null || true
	@if [ -f ${BUILD_CMD} ]; \
		then \
			./${BUILD_CMD} && \
			if [ -f /tmp/${BUILD_CMD}.cdk.pprof/mem.pprof ]; \
			then \
				read -p "# Press enter to open a pprof instance" JUNK \
				&& go tool pprof -http=:8080 /tmp/${BUILD_CMD}.cdk.pprof/mem.pprof; \
			else \
				echo "# missing /tmp/${BUILD_CMD}.cdk.pprof/mem.pprof"; \
			fi ; \
		fi

#
# Cross-Compilation Targets
#

dev-linux-amd64: export GOOS=linux
dev-linux-amd64: export GOARCH=amd64
dev-linux-amd64: dev

dev-linux-mips64: export GOOS=linux
dev-linux-mips64: export GOARCH=mips64
dev-linux-mips64: dev

dev-darwin-amd64: export GOOS=darwin
dev-darwin-amd64: export GOARCH=amd64
dev-darwin-amd64: dev

#
# Experiments
#

tail-logs:
	@tail -F ${BUILD_CMD}.cdk.log ${BUILD_CMD}.build.log

