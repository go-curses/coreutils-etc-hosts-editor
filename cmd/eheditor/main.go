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

	clcli "github.com/go-corelibs/cli"
	"github.com/go-curses/cdk/log"

	"github.com/go-curses/coreutils-etc-hosts-editor/ui"
)

var (
	BuildVersion = "0.7.1"
	BuildRelease = "trunk"
)

func init() {
	cli.FlagStringer = clcli.NewFlagStringer().
		PruneDefaultBools(true).
		Make()
}

func main() {
	ehe := ui.NewUI(
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
		Name:    "read-only",
		Usage:   "do not write any changes to the etc hosts file",
		Aliases: []string{"r"},
	})
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Usage:   "display the version",
		Aliases: []string{"v"},
	}
	clcli.ClearEmptyCategories(appCLI.Flags)
	if err := ehe.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
