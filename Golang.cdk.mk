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

GOLANG_MAKEFILE_KEYS += CDK

GOLANG_CDK_MK_VERSION := v0.1.1

CUSTOM_HELP_SECTIONS += CDK_HELP

CDK_HELP_NAME := "go-curses"
CDK_HELP_KEYS := DRUN PCPU PMEM

CDK_HELP_DRUN_TARGET := debug-run
CDK_HELP_DRUN_USAGE  := run the debug build (and sanely handle crashes)

CDK_HELP_PCPU_TARGET := profile.cpu
CDK_HELP_PCPU_USAGE  := run the dev build and profile CPU

CDK_HELP_PMEM_TARGET := profile.mem
CDK_HELP_PMEM_USAGE  := run the dev build and profile MEM

debug-run: export GO_CDK_LOG_FILE=./${BUILD_NAME}.cdk.log
debug-run: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
debug-run: export GO_CDK_LOG_FULL_PATHS=true
debug-run: debug
	@if [ -f ${BUILD_NAME} ]; \
	then \
		echo "# running: ${BUILD_NAME} ${RUN_ARGS}"; \
		( ./${BUILD_NAME} ${RUN_ARGS} ) 2>> ${GO_CDK_LOG_FILE}; \
		if [ $$? -ne 0 ]; \
		then \
			stty sane; echo ""; \
			echo "# ${BUILD_NAME} crashed, see: ./${BUILD_NAME}.cdk.log"; \
			read -p "# Press <Enter> to reset terminal, <Ctrl+C> to cancel" RESP; \
			reset; \
			echo "# ${BUILD_NAME} crashed, terminal reset, see: ./${BUILD_NAME}.cdk.log"; \
		else \
			echo "# ${BUILD_NAME} exited normally."; \
		fi; \
	else \
		echo "# ${BUILD_NAME} not found"; \
		false; \
	fi

debug-dlv: export GO_CDK_LOG_FILE=./${BUILD_NAME}.cdk.log
debug-dlv: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
debug-dlv: export GO_CDK_LOG_FULL_PATHS=true
debug-dlv: debug
	@if [ -f ${BUILD_NAME} ]; \
	then \
		echo "# running: ${BUILD_NAME} ${RUN_ARGS}"; \
		( dlv.sh ./${BUILD_NAME} ${RUN_ARGS} ) 2>> ${GO_CDK_LOG_FILE}; \
		if [ $$? -ne 0 ]; \
		then \
			stty sane; echo ""; \
			echo "# ${BUILD_NAME} crashed, see: ./${BUILD_NAME}.cdk.log"; \
			read -p "# Press <Enter> to reset terminal, <Ctrl+C> to cancel" RESP; \
			reset; \
			echo "# ${BUILD_NAME} crashed, terminal reset, see: ./${BUILD_NAME}.cdk.log"; \
		else \
			echo "# ${BUILD_NAME} exited normally."; \
		fi; \
	else \
		echo "# ${BUILD_NAME} not found"; \
		false; \
	fi

profile.cpu: export GO_CDK_LOG_FILE=./${BUILD_NAME}.cdk.log
profile.cpu: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
profile.cpu: export GO_CDK_LOG_FULL_PATHS=true
profile.cpu: export GO_CDK_PROFILE_PATH=/tmp/${BUILD_NAME}.cdk.pprof
profile.cpu: export GO_CDK_PROFILE=cpu
profile.cpu: debug
	@mkdir -v /tmp/${BUILD_NAME}.cdk.pprof 2>/dev/null || true
	@if [ -f ${BUILD_NAME} ]; \
		then \
			./${BUILD_NAME} && \
			if [ -f /tmp/${BUILD_NAME}.cdk.pprof/cpu.pprof ]; \
			then \
				read -p "# Press enter to open a pprof instance" JUNK \
				&& go tool pprof -http=:8080 /tmp/${BUILD_NAME}.cdk.pprof/cpu.pprof ; \
			else \
				echo "# missing /tmp/${BUILD_NAME}.cdk.pprof/cpu.pprof"; \
			fi ; \
	else \
		echo "# ${BUILD_NAME} not found"; \
		false; \
	fi

profile.mem: export GO_CDK_LOG_FILE=./${BUILD_NAME}.log
profile.mem: export GO_CDK_LOG_LEVEL=${LOG_LEVEL}
profile.mem: export GO_CDK_LOG_FULL_PATHS=true
profile.mem: export GO_CDK_PROFILE_PATH=/tmp/${BUILD_NAME}.cdk.pprof
profile.mem: export GO_CDK_PROFILE=mem
profile.mem: debug
	@mkdir -v /tmp/${BUILD_NAME}.cdk.pprof 2>/dev/null || true
	@if [ -f ${BUILD_NAME} ]; \
		then \
			./${BUILD_NAME} && \
			if [ -f /tmp/${BUILD_NAME}.cdk.pprof/mem.pprof ]; \
			then \
				read -p "# Press enter to open a pprof instance" JUNK \
				&& go tool pprof -http=:8080 /tmp/${BUILD_NAME}.cdk.pprof/mem.pprof; \
			else \
				echo "# missing /tmp/${BUILD_NAME}.cdk.pprof/mem.pprof"; \
			fi ; \
	else \
		echo "# ${BUILD_NAME} not found"; \
		false; \
	fi
