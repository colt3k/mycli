package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/colt3k/utils/version"

	log "github.com/colt3k/nglog/ng"
	"github.com/colt3k/utils/file"
	"github.com/colt3k/utils/lock"

	"github.com/colt3k/mycli"
	"github.com/colt3k/mycli/custom"
)

const (
	appName     = "mycli"
	title       = "myCLI"
	description = "Demo application for my CLI"
	author      = "FirstName LastName (example@gmail.com)"
	copyright   = "(c) 2018 Example Inc.,"
	companyName = "colt3k"
)

var (
	logDir  = file.HomeFolder()
	logfile                                       = filepath.Join(logDir, appName+".log")
	t                                             bool
	capture, protocol, path, url, applicationName string
	port, t2, t3, t4                              int64
	countStringList                               mycli.StringList
	c                                             *mycli.CLI
	clients                                       custom.Clients
	l                                             *lock.Lock
	locked                                        bool
)

func init() {
	l = lock.New(appName)
	if l.Try() {
		locked = true
	}
	fa, err := log.NewFileAppender("*", logfile, "", 0)
	if err != nil {
		log.Logf(log.FATAL, "issue creating file appender\n%+v", err)
	}
	ca := log.NewConsoleAppender("*")
	log.Modify(log.LogLevel(log.INFO), log.ColorsOn(), log.Appenders(ca, fa))
}
func setLogger() error {

	if logDir != file.HomeFolder() || mycli.Debug {
		logfile = filepath.Join(logDir, companyName, appName+".log")
		// update our logger
		fa, err := log.NewFileAppender("*", logfile, "", 0)
		if err != nil {
			return err
		}
		ca := log.NewConsoleAppender("*")
		if mycli.Debug {
			switch mycli.DebugLevel {
			case 2:
				log.Modify(log.LogLevel(log.DBGL2), log.ColorsOn(), log.Appenders(ca, fa))
			case 3:
				log.Modify(log.LogLevel(log.DBGL3), log.ColorsOn(), log.Appenders(ca, fa))
			default:
				log.Modify(log.LogLevel(log.DEBUG), log.ColorsOn(), log.Appenders(ca, fa))
			}
			log.Logf(log.DEBUG, "> debug mode on level %d", mycli.DebugLevel)
		} else {
			log.Modify(log.LogLevel(log.INFO), log.ColorsOn(), log.Appenders(ca, fa))
		}
	}
	return nil
}
func main() {

	setupFlags()
	if locked {
		l.Unlock()
		log.Logln(log.DEBUG, "unlocked")
	}
	os.Exit(0)
}
func setupFlags() {

	opts := []string{"gc"}
	c = mycli.NewCli(nil, nil)
	c.Title = title
	c.Description = description
	c.Version = version.VERSION
	c.BuildDate = version.BUILDDATE
	c.GitCommit = version.GITCOMMIT
	c.GoVersion = version.GOVERSION
	c.Author = author
	c.Copyright = copyright
	c.MainAction = func() { fmt.Println("Main Action Ran") }
	c.PostGlblAction = func() error { return setLogger() }
	//c.EnvPrefix = "T"
	c.Flgs = []mycli.CLIFlag{
		&mycli.StringFlg{Variable: &capture, Name: "capture", ShortName: "cap", Usage: "Used to test string", Options: []string{"hello", "bye"}},
		&mycli.StringFlg{Variable: &path, Name: "path", Usage: "Used to test path with slash"},
		&mycli.StringFlg{Variable: &url, Name: "url", Usage: "Used to test url with slashes"},
		&custom.TomlFlg{Variable: &clients, Name: "clients", Usage: "Set name to toml table type"},
	}

	c.Cmds = []*mycli.CLICommand{
		{
			Name:           "server",
			ShortName:      "s",
			Usage:          "use as a server",
			Value:          nil,
			Action:         func() { runAsServer() },
			PreAction:      func() { checkDebug("cmd") },
			PostAction:     nil,
			Flags: []mycli.CLIFlag{
				&mycli.StringFlg{Variable: &protocol, Name: "protocol", ShortName: "proto", Usage: "Set Protocol http(s)", Value: "http"},
				// if value is set and required passed in value has to be different or it will think it wasn't set
				&mycli.Int64Flg{Variable: &t2, Name: "port", ShortName: "p", Usage: "Change server port", Value: 8080},
			},
		},
		{
			Name:           "client",
			ShortName:      "c",
			Usage:          "use as a client",
			Value:          nil,
			Action:         func() { runAsClient() },
			PreAction:      func() { checkDebug("cmd") },
			PostAction:     nil,
			Flags: []mycli.CLIFlag{
				&mycli.Int64Flg{Variable: &t3, Name: "port", ShortName: "p", Usage: "Change client port", Value: 8080, Required: true},
			},
		},
		{
			Name:     "clients",
			Usage:    "use as a test to hide command but capture value",
			Variable: &clients,
			Hidden:   true,
		},
		{
			Name:           "weserve",
			Usage:          "use as a client",
			SubCommands: []*mycli.CLICommand{
				{
					Name:      "config",
					ShortName: "c",
					Usage:     "use config file",
					Action: func() {
						log.Println("ran clients config")
					},
					Flags: []mycli.CLIFlag{
						&mycli.Int64Flg{Variable: &t4, Name: "port", ShortName: "p", Usage: "Set Port", Value: 9111, Required: false},
						&mycli.StringFlg{Variable: &applicationName, Name: "application", Usage: "Select application name", Required: true, Options: opts},
					},
				},
				{
					Name:      "cmdln",
					ShortName: "cl",
					Usage:     "use command line",
					Action: func() {
						log.Println("ran clients cmdline")
					},
					Flags: []mycli.CLIFlag{
						&mycli.Int64Flg{Variable: &t4, Name: "port", ShortName: "p", Usage: "Set Port", Value: 9111, Required: false},
						&mycli.StringFlg{Variable: &applicationName, Name: "application", ShortName: "a", Usage: "Select application name", Required: true, Options: opts},
					},
				},
			},
		},
	}

	err := c.Parse()
	if err != nil {
		if locked {
			l.Unlock()
			log.Logln(log.DEBUG, "unlocked")
		}
		log.Logf(log.FATAL, "error(s)\n%+v", err)
	}
}

func checkDebug(txt string) {
	showStuff()
}
func runAsServer() {
	fmt.Println("\n** Running as server **")
	fmt.Println()
}
func runAsClient() {
	fmt.Println("\n** Running as client **")
	fmt.Println()
}

func showStuff() {
	fmt.Println("Path: ", path)
	fmt.Println("Debug Flag:", mycli.Debug)
	fmt.Println("Test Flag:", t)
	fmt.Println("Capture Flag:", capture)
	fmt.Println("Server Port:", t2)
	fmt.Println("Client Port:", t3)
	fmt.Println("Subcommand weserve config Port:", t4)
	fmt.Println("Protocol:", protocol)
	fmt.Println("Help Flag:", c.Help())
	fmt.Println("SubstringList Flag:", countStringList)
	fmt.Println("clients:", clients)
	fmt.Println()
}
