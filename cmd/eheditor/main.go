package main

import (
	_ "embed"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/go-curses/ctk"

	"github.com/go-curses/cdk"
	cstrings "github.com/go-curses/cdk/lib/strings"
	"github.com/go-curses/cdk/log"
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
	BuildVersion = "0.1.0"
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

//go:embed eheditor.accelmap
var eheditorAccelMap string

//go:embed eheditor.styles
var eheditorStyles string

func main() {
	app := ctk.NewApplication(
		"eheditor",
		"etc hosts editor",
		"command line utility for managing the OS /etc/hosts file",
		BuildVersion+" ("+BuildRelease+")",
		"eheditor",
		"/etc/hosts editor",
		"/dev/tty",
	)
	app.CLI().UsageText = "eheditor [options] [/etc/hosts]"
	app.CLI().HideHelpCommand = true
	app.CLI().EnableBashCompletion = true
	app.CLI().UseShortOptionHandling = true
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Usage:   "display the version",
		Aliases: []string{"v"},
	}
	app.AddFlag(&cli.BoolFlag{
		Name:    "read-only",
		Usage:   "do not write any changes to the etc hosts file",
		Aliases: []string{"r"},
	})
	app.Connect(cdk.SignalStartup, "eheditor-startup-handler", startup)
	// app.Connect(cdk.SignalStartupComplete, "eheditor-startup-complete-handler", startupComplete)
	app.Connect(cdk.SignalShutdown, "eheditor-quit-handler", shutdown)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
