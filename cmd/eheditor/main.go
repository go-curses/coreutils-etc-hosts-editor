// Copyright (c) 2023  The Go-Curses Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/go-curses/cdk"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/log"

	"github.com/go-curses/coreutils-etc-hosts-editor/ui/eheditor"
)

// Build Configuration Flags
// setting these will enable command line flags and their corresponding features
// use `go build -v -ldflags="-X 'main.IncludeLogFullPaths=false'"`
var (
	IncludeProfiling          = "false"
	IncludeLogFile            = "false"
	IncludeLogFormat          = "false"
	IncludeLogFullPaths       = "false"
	IncludeLogLevel           = "false"
	IncludeLogLevels          = "false"
	IncludeLogTimestamps      = "false"
	IncludeLogTimestampFormat = "false"
	IncludeLogOutput          = "false"
)

var (
	BuildVersion = "0.1.1"
	BuildRelease = "trunk"
)

func init() {
	cdk.Build.Profiling = cstrings.IsTrue(IncludeProfiling)
	cdk.Build.LogFile = cstrings.IsTrue(IncludeLogFile)
	cdk.Build.LogFormat = cstrings.IsTrue(IncludeLogFormat)
	cdk.Build.LogFullPaths = cstrings.IsTrue(IncludeLogFullPaths)
	cdk.Build.LogLevel = cstrings.IsTrue(IncludeLogLevel)
	cdk.Build.LogLevels = cstrings.IsTrue(IncludeLogLevels)
	cdk.Build.LogTimestamps = cstrings.IsTrue(IncludeLogTimestamps)
	cdk.Build.LogTimestampFormat = cstrings.IsTrue(IncludeLogTimestampFormat)
	cdk.Build.LogOutput = cstrings.IsTrue(IncludeLogOutput)
}

func main() {
	ehe := eheditor.NewEheditor(
		"eheditor",
		"etc hosts editor",
		"command line utility for managing the OS /etc/hosts file",
		BuildVersion+" ("+BuildRelease+")",
		"eheditor",
		"/etc/hosts editor",
		"/dev/tty",
	)
	appCLI := ehe.App.CLI()
	appCLI.UsageText = "eheditor [options] [/etc/hosts]"
	appCLI.HideHelpCommand = true
	appCLI.EnableBashCompletion = true
	appCLI.UseShortOptionHandling = true
	ehe.App.AddFlag(&cli.BoolFlag{
		Name:   "with-old-version",
		Usage:  "include the original version user-interface",
		Hidden: true,
	})
	ehe.App.AddFlag(&cli.BoolFlag{
		Name:    "read-only",
		Usage:   "do not write any changes to the etc hosts file",
		Aliases: []string{"r"},
	})
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Usage:   "display the version",
		Aliases: []string{"v"},
	}
	if err := ehe.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}