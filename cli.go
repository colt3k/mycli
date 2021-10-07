package mycli

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"

	"github.com/colt3k/nglog/ng"
)

const (
	UseHTTPProxy  = "Sets http_proxy for network connections"
	UseHTTPSProxy = "Sets https_proxy for network connections"
	UseNoProxy    = "Sets no_proxy for network connections"
)

var (
	configfile string
	ProxyHTTP  string
	ProxyHTTPS string
	ProxyNO    string
	ProxySOCKS string
	Debug      bool
	DebugLevel int64
	GenerateBashCompletion bool
)

type CLICommand struct {
	// Name of the command and passed for use
	Name                   string
	// ShortName used for execution but provides a shorter name
	ShortName              string
	// Usage definition of what this command accomplishes
	Usage                  string
	// Variable used to process a file full of configurations see custom/flgtoml.go as an example used with Hidden:true
	Variable               interface{}
	// Value unused
	Value                  interface{}
	// PreAction perform some action prior to the Action defined
	PreAction              interface{}
	// Action main action to perform for this Command
	Action                 interface{}
	// PostAction perform some action after the main Action
	PostAction             interface{}
	// Flags are command flags local to this command
	Flags                  []CLIFlag
	// FS reserved for internal use
	FS                     *flag.FlagSet
	// BashCompletion should be set to mycli.BashCompletionSub for sub command completion
	BashCompletion         interface{}
	generateBashCompletion bool
	// Hidden stops from showing in help
	Hidden                 bool
	help                   bool
	// SubCommands ability to create sub commands of a top command
	SubCommands            Commands
}

type Commands []*CLICommand

// Convert this section to a JSON string
func (c *CLICommand) RetrieveConfigValue(val interface{}, name string) error {
	treeMap := val.(*toml.Tree).ToMap()
	valS := treeMap[c.Name]
	wrapper := make(map[string]interface{}, 1)
	wrapper[c.Name] = valS
	bytes, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		return err
	}
	// set to Variable here so no need to go further as in other types
	err = json.Unmarshal(bytes, c.Variable)
	if err != nil {
		log.Fatalf("issue unmarshalling\n%+v", err)
	}

	return nil
}

// AppInfo supplies all pertinent information for the application
type AppInfo struct {
	// Version typically v0.0.1 format of version
	Version     string
	// BuildDate typically set to a unix timestamp format
	BuildDate   string
	// GitCommit the short git commit hash
	GitCommit   string
	// GoVersion go version application was built upon
	GoVersion	string
	// Title plain text name for the application
	Title       string
	// Description detailed purpose of the application
	Description string
	Usage       string
	Author      string
	// Copyright typically company or developer copyright i.e. [ (c) 4-digit-year company/user ]
	Copyright   string
}

type UsageAdapter interface {
	UsageText(*CLICommand)
}
type UsageDisplay struct {
	adapter UsageAdapter
}

func (u *UsageDisplay) UsageText(cmd *CLICommand) {
	cmd.FS.Usage()
}

type FatalAdapter interface {
	PrintNotice(string)
	PrintNoticeSubCmd(string, string)
}

type Fatal struct {
	adapter FatalAdapter
}

func (f *Fatal) PrintNotice(name string) {
	log.Fatal(ng.Red("required flag '-%s' not set\n", name))
}
func (f *Fatal) PrintNoticeSubCmd(name, cmd string) {
	log.Fatal(ng.Red("required flag '-%s' not set on sub-command: %s\n", name, cmd))
}

// CLI command line struct
type CLI struct {
	*AppInfo
	// Flgs location to set all global flags
	Flgs                   []CLIFlag
	// Cmds global commands your application supports
	Cmds                   []*CLICommand
	// PostGlblAction runs an action after processing Global flags
	PostGlblAction         interface{}
	// MainAction this is a default if no Command is specified when the application is run
	MainAction             interface{}
	cur                    *CLICommand
	// BashCompletion typically set to the built in default of mycli.BashCompletionMain
	BashCompletion         interface{}
	// VersionPrint an overridable function that prints by default the set Version, BuildDate, GitCommit, GoVersion
	VersionPrint           interface{}
	generateBashCompletion bool
	Writer                 io.Writer
	// DisableEnvVars disable all environment variables
	DisableEnvVars		   bool
	// EnvPrefix a prefix you can define to use on Environment Variables for values used in the application default "T"
	EnvPrefix              string
	// TestMode reserved for internal testing
	TestMode               bool
	fatalAdapter           FatalAdapter
	usageAdapter           UsageAdapter
	help, debug, version   bool
}

// NewCli creates an instance of the CLI application
func NewCli(f FatalAdapter, u UsageAdapter) *CLI {
	a := new(AppInfo)
	t := new(CLI)
	if f != nil {
		t.fatalAdapter = f
	} else {
		t.fatalAdapter = new(Fatal)
	}
	if u != nil {
		t.usageAdapter = u
	} else {
		t.usageAdapter = new(UsageDisplay)
	}
	t.AppInfo = a
	t.Writer = os.Stdout
	t.Flgs = make([]CLIFlag, 0)
	t.Cmds = make([]*CLICommand, 0)

	t.VersionPrint = func() {
		fmt.Printf("\nversion=%s build=%s revision=%s goversion=%s\n\n", a.Version, a.BuildDate, a.GitCommit, a.GoVersion)
	}
	t.BashCompletion = BashCompletionMain
	t.DisableEnvVars = true
	t.EnvPrefix = "T"
	return t
}

func (c *CLI) Help() bool {
	return c.help
}

func (c *CLI) addDefaultFlags() {

	dfFlgs := make([]CLIFlag, 0)
	if !c.findFlag("help", c.Flgs) {
		flg := c.setupHelpFlag()
		dfFlgs = append(dfFlgs, flg)
	}
	if !c.findFlag("debug", c.Flgs) {
		flg := c.setupDebugFlag()
		dfFlgs = append(dfFlgs, flg)
	}
	if !c.findFlag("debugLevel", c.Flgs) {
		flg := c.setupDebugLevelFlag()
		dfFlgs = append(dfFlgs, flg)
	}
	if !c.findFlag("version", c.Flgs) {
		flg := c.setupVersionFlag()
		dfFlgs = append(dfFlgs, flg)
	}
	if !c.findFlag("config", c.Flgs) {
		flg := c.setupConfigFlag()
		dfFlgs = append(dfFlgs, flg)
	}
	if !c.findFlag("proxyhttp", c.Flgs) {
		flgs := c.setupProxyFlags()
		for _, f := range flgs {
			dfFlgs = append(dfFlgs, f)
		}
	}
	for _, d := range c.Flgs {
		dfFlgs = append(dfFlgs, d)
	}
	c.Flgs = dfFlgs
}

// SetupEnvVars Loop through all Flags and Command Flags then set EnvVars based on Prefix and NAME or Override
func (c *CLI) SetupEnvVars() {

	tmp := make([]CLIFlag, 0)
	for i, d := range c.Flgs {
		if !d.GEnvVarExclude() {
			envvar := c.buildEnvVar(d)
			d.SetEnvVar(envvar)
			c.Flgs[i] = d
			tmp = append(tmp, c.Flgs[i])
		}
	}

	for _, j := range c.Cmds {
		for i, d := range j.Flags {
			if !d.GEnvVarExclude() {
				envvar := c.buildEnvVar(d)
				d.SetEnvVar(envvar)
				j.Flags[i] = d
			}
		}
	}

	// setup ENV for SubCommands
	for _, j := range c.Cmds {
		for _, q := range j.SubCommands {
			for i, d := range q.Flags {
				if !d.GEnvVarExclude() {
					envvar := c.buildEnvVar(d)
					d.SetEnvVar(envvar)
					q.Flags[i] = d
				}
			}
		}
	}
}
func (c *CLI) buildEnvVar(f CLIFlag) string {

	if len(f.GEnvVar()) == 0 {
		if len(c.EnvPrefix) == 0 {
			return strings.ToUpper(f.GName())
		} else {
			return strings.ToUpper(c.EnvPrefix + "_" + f.GName())
		}
	} else {
		if len(c.EnvPrefix) == 0 {
			return strings.ToUpper(f.GEnvVar())
		} else {
			return strings.ToUpper(c.EnvPrefix + "_" + f.GEnvVar())
		}
	}
}

func (c *CLI) ValidateFlgKind() error {
	for _, d := range c.Flgs {
		err := d.Kind()
		if err != nil {
			return err
		}
	}
	for _, d := range c.Cmds {
		for _, j := range d.Flags {
			err := j.Kind()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *CLI) ValidateValues(commands bool) error {

	for _, f := range c.Flgs {
		if !f.ValidValue() {
			return &InvalidValueError{f.GName(), f.ValueAsString(), f.GOptions()}
		}
	}

	if commands {
		for _, cmd := range c.Cmds {
			for _, f := range cmd.Flags {
				if !f.ValidValue() {
					return &InvalidValueError{f.GName(), f.ValueAsString(), f.GOptions()}
				}
			}
			for _, s := range cmd.SubCommands {
				for _, f := range s.Flags {
					if !f.ValidValue() {
						return &InvalidValueError{f.GName(), f.ValueAsString(), f.GOptions()}
					}
				}
			}
		}
	}
	return nil
}

func FixPath(path string) string {
	if !filepath.IsAbs(path) {
		pth, _ := filepath.Abs(path)
		return pth
	}
	return path
}

func (c *CLI) parseConfigFile() error {
	debug := false
	// no config file passed, return
	if len(configfile) == 0 {
		return nil
	}
	// if doesn't exist return
	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		log.Printf("!!! config file not found %v\n", configfile)
		return nil
	}
	if len(strings.TrimSpace(configfile)) > 0 {
		configfile = FixPath(configfile)
		tree, err := LoadToml(configfile)
		if err != nil {
			log.Printf("!!! issue loading config toml file %v\n", err)
			return err
		}
		// find any missing values and set them from the tree
		for _, f := range c.Flgs {
			key := f.GName()
			if tree.Has(key) {
				err = f.RetrieveConfigValue(tree, key)
				if Err(err) {
					log.Printf("!!! issue retrieving flag value from config toml file %v\n", err)
					return err
				}
				if key == "debug" && f.GVariableToString() == "true"{
					debug = true
				}
				if debug {
					log.Printf("- config file has global flag %v value found of %v", key, f.GVariableToString())
				}
			}
		}

		for _, cmd := range c.Cmds {
			for _, f := range cmd.Flags {
				key := cmd.Name + "." + f.GName()
				if tree.Has(key) {
					err = f.RetrieveConfigValue(tree, key)
					if Err(err) {
						log.Printf("!!! issue retrieving command value from config toml file %v\n", err)
						return err
					}
					if debug {
						log.Printf("- config file has command flag %v value found of %v", key, f.GVariableToString())
					}
				}
			}
			for _, subcmd := range cmd.SubCommands {
				for _, f := range subcmd.Flags {
					key := cmd.Name + "." + subcmd.Name + "." +f.GName()
					//log.Printf("- looking for %v", key)
					if tree.Has(key) {
						err = f.RetrieveConfigValue(tree, key)
						if Err(err) {
							log.Printf("!!! issue retrieving command value from config toml file %v\n", err)
							return err
						}
						if debug {
							log.Printf("- config file has subcommand flag %v value found of %v", key, f.GVariableToString())
						}
					}
				}
			}
			if cmd.Hidden && tree.Has(cmd.Name) {
				key := cmd.Name
				err = cmd.RetrieveConfigValue(tree, key)
				if Err(err) {
					log.Printf("!!! issue retrieving config value from config toml file %v\n", err)
					return err
				}
				if debug {
					log.Printf("- config file has hidden command %v value found of %v", key, cmd.Variable)
				}
			}
		}
	}
	return nil
}
func (c *CLI) Parse() error {

	// add default flags, help, debug, debuglevel, version, config
	c.addDefaultFlags()
	// Disable all Env Variables
	if !c.DisableEnvVars {
		c.SetupEnvVars()
	}
	err := c.ValidateFlgKind()
	if err != nil {
		log.Fatalf("issue validating flag kind\n%+v", err)
	}
	// Pre process Global Flags
	c.buildFlags(flag.CommandLine, c.Flgs, nil)

	//log.Println("passed parameters: ", os.Args[1:])
	for _,v := range os.Args[1:] {
		if v == "--generate-bash-completion" {
			GenerateBashCompletion = true
		}
	}
	flag.Parse()
	err = c.retrieveEnvVal(c.Flgs)
	if Err(err) {
		return err
	}

	// loop Flags and find environment values, if not set on commandline set value to ENV value

	if c.PostGlblAction != nil {
		err = runAction(c.PostGlblAction)
		if err != nil {
			return err
		}
	}
	if c.version && c.VersionPrint != nil {
		err = runAction(c.VersionPrint)
		if err != nil {
			return err
		}
	}
	//Reset and Process Global and Commands
	ResetForTesting(nil)
	c.buildFlags(flag.CommandLine, c.Flgs, nil)
	c.buildCmds()

	//log.Println("passed parameters: ", os.Args[1:])
	flag.Parse()

	//retrieve environment values if set and flag wasn't passed
	err = c.retrieveEnvVal(c.Flgs)
	if Err(err) {
		return err
	}

	for _, d := range c.Cmds {
		err = c.retrieveEnvVal(d.Flags)
		if Err(err) {
			return err
		}
	}

	for _, d := range c.Cmds {
		for _, q := range d.SubCommands {
			err = c.retrieveEnvVal(q.Flags)
			if Err(err) {
				return err
			}
		}
	}

	// anything not set use config file to set it
	err = c.parseConfigFile()
	if Err(err) {
		return err
	}

	c.checkRequired("", c.Flgs) // see if required ones are set
	if c.help && !GenerateBashCompletion {
		c.printUsage()
		if !c.TestMode {
			os.Exit(1)
		}
	}
	err = c.ValidateValues(false)
	if Err(err) {
		return err
	}

	//check for bash completion flag
	//fmt.Println("args: gen bash? ", os.Args, c.generateBashCompletion)
	if c.generateBashCompletion {
		c.BashCompletion.(func(cli *CLI))(c)
		if !c.TestMode {
			os.Exit(1)
		}
	}

	if Debug && !GenerateBashCompletion {
		for _, f := range c.Flgs {
			fmt.Printf("Flag '%s': %v\n", f.GName(), f.GVariableToString())
		}
	}

	// Process input
	var parentCmd string
	var subCmd string
	var activeCmd *CLICommand
	var pos int

	for _, d := range c.Cmds {
		for i, a := range os.Args {
			if len(os.Args) > 1 && (a == strings.ToLower(d.Name) || a == strings.ToLower(d.ShortName)) {
				pos = i + 1
				t := d
				c.cur = t
				t.FS.Usage = c.flagSetUsage
				//fmt.Printf("Args Size: %d and position 1 %v equals lowered name %s\n",len(os.Args), os.Args[1],strings.ToLower(d.Name))
				activeCmd = t
				parentCmd = t.Name
				//log.Println("Active command ", t.Name)
				// find subcommand to set instead of main command
				subCs := make([]string, 0)
				//foundSub := false
				for _, k := range d.SubCommands {
					subCs = append(subCs, k.Name)
					//log.Println("check sub commands for ", d.Name, " is ", a, " = ", k.Name, " or ", k.ShortName)
					for q, b := range os.Args {
						if b == strings.ToLower(k.Name) || b == strings.ToLower(k.ShortName) {
							pos = q + 1
							t := k
							c.cur = t
							t.FS.Usage = c.flagSetUsage
							//fmt.Printf("Args Size: %d and position 1 %v equals lowered name %s\n",len(os.Args), os.Args[1],strings.ToLower(d.Name))
							activeCmd = t
							subCmd = t.Name
							//log.Println("Active sub command ", t.Name)
							//foundSub = true
						}
					}
				}
				//if !foundSub {
				//	t.Usage = strings.Join(subCs, ",")
				//}

				break
			}
		}
	}

	if activeCmd != nil {
		if Debug && !GenerateBashCompletion {
			if c.MainAction != nil {
				fmt.Println("")
				fmt.Println("- Skipping Main Action and running requested Commands. -")
				fmt.Println("")
			}
			if len(subCmd) > 0 {
				fmt.Printf("Active command : %v %v\n", parentCmd, subCmd)
			} else {
				fmt.Printf("Active command : %v\n", activeCmd.Name)
			}
		}
		err := activeCmd.FS.Parse(os.Args[pos:])
		PanicErr(err)
		err = c.ValidateValues(true)
		if Err(err) {
			return err
		}
		if activeCmd.help {
			c.usageAdapter.UsageText(activeCmd)
			if c.TestMode {
				return nil
			}
			os.Exit(1)
		}

		// If generateBashCompletion flag is true and BashCompletion interface is not nil, execute it
		if activeCmd.generateBashCompletion {
			// If nil then we set to a subcommand completion, main is prepopulated
			if activeCmd.BashCompletion == nil {
				activeCmd.BashCompletion = BashCompletionSub
			}
			activeCmd.BashCompletion.(func(cli *CLI, cm *CLICommand))(c, activeCmd)
			if c.TestMode {
				return nil
			}
			os.Exit(1)
		}
		// If we find generate-bash-completion in the command line exit
		if strings.Index(os.Args[len(os.Args)-1], "generate-bash-completion") > -1 {
			//fmt.Println("generate_bash_completion_IS_BASH_AT_END_EXITING!!!")
			os.Exit(1)
		}
		// If in debug mode print out subcommand
		if Debug && !GenerateBashCompletion {
			for _, f := range activeCmd.Flags {
				fmt.Printf("Subcommand '%s' Flag '%s': %v\n", activeCmd.Name, f.GName(), f.GVariableToString())
			}
		}

		c.checkRequired(activeCmd.FS.Name(), activeCmd.Flags)
		//Execute action
		if activeCmd.PreAction != nil {
			err = runAction(activeCmd.PreAction)
			if err != nil {
				return err
			}
		}

		err = runAction(activeCmd.Action)
		if err != nil {
			return err
		}

		if activeCmd.PostAction != nil {
			err = runAction(activeCmd.PostAction)
			if err != nil {
				return err
			}
		}
	} else if c.MainAction != nil {
		err = runAction(c.MainAction)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("!!! no command set to run")
	}
	return nil
}

func runAction(act interface{}) error {
	var err error
	switch act.(type) {
	case func():
		act.(func())()
	case func() error:
		err = act.(func() error)()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown action type")
	}
	return nil
}
func (c *CLI) findFlag(flgName string, flgs []CLIFlag) bool {
	for _, d := range flgs {
		if d.GName() == flgName {
			return true
		}
	}
	return false
}
func (c *CLI) Flag(flgName string, flgs []CLIFlag) CLIFlag {
	for _, d := range flgs {
		if d.GName() == flgName {
			return d
		}
	}
	return nil
}
func (c *CLI) Command(name string) *CLICommand {
	for _, d := range c.Cmds {
		if d.Name == name {
			return d
		}
	}
	return nil
}
func (c *CLI) setupHelpFlag() CLIFlag {
	return &BoolFlg{Variable: &c.help, Name: "help", ShortName: "h", Usage: "print commands", EnvVarExclude: true, Hidden: true}
}
func (c *CLI) setupDebugFlag() CLIFlag {
	return &BoolFlg{Variable: &Debug, Name: "debug", ShortName: "d", Usage: "flag set to debug", EnvVarExclude: true}
}
func (c *CLI) IsDebug() bool {
	return Debug
}
func (c *CLI) setupDebugLevelFlag() CLIFlag {
	return &Int64Flg{Variable: &DebugLevel, Name: "debugLevel", ShortName: "dbglvl", Usage: "set debug level", EnvVarExclude: true, Value: 0}
}
func (c *CLI) DebugLevel() int64 {
	return DebugLevel
}
func (c *CLI) setupVersionFlag() CLIFlag {
	return &BoolFlg{Variable: &c.version, Name: "version", ShortName: "v", Usage: "flag to show version", EnvVarExclude: true, Hidden: true}
}
func (c *CLI) setupConfigFlag() CLIFlag {
	if !c.DisableEnvVars {
		return &StringFlg{Variable: &configfile, Name: "config", ShortName: "c", EnvVar: "config_filepath", Usage: "config file path"}
	} else {
		return &StringFlg{Variable: &configfile, Name: "config", ShortName: "c", Usage: "config file path"}
	}
}
func (c *CLI) setupProxyFlags() []CLIFlag {

	return []CLIFlag{
		&StringFlg{Variable: &ProxyHTTP, Name: "proxyhttp", EnvVar: "HTTP_PROXY", Usage: UseHTTPProxy},
		&StringFlg{Variable: &ProxyHTTPS, Name: "proxyhttps", EnvVar: "HTTPS_PROXY", Usage: UseHTTPSProxy},
		&StringFlg{Variable: &ProxyNO, Name: "noproxy", EnvVar: "NO_PROXY", Usage: UseNoProxy},
	}
}
func (c *CLI) IsProxySet() bool {
	if len(ProxyHTTP) > 0 || len(ProxyHTTPS) > 0 || len(ProxyNO) > 0 {
		return true
	}
	return false
}
func (c *CLI) GetHttpProxy() string {
	return ProxyHTTP
}
func (c *CLI) GetHttpsProxy() string {
	return ProxyHTTPS
}
func (c *CLI) GetNoProxy() string {
	return ProxyNO
}
func (c *CLI) setupBashFlag(cm *CLICommand) CLIFlag {

	//set flag -generate-bash-completion
	var bf *BoolFlg
	if cm != nil {
		bf = &BoolFlg{Variable: &cm.generateBashCompletion, Name: "generate-bash-completion", Usage: "provides bash completion", Hidden: true}
		bf.Action = cm.BashCompletion
	} else {
		bf = &BoolFlg{Variable: &c.generateBashCompletion, Name: "generate-bash-completion", Usage: "provides bash completion", Hidden: true}
		bf.Action = c.BashCompletion
	}
	return bf
}
func (c *CLI) buildFlags(flgSet *flag.FlagSet, flgs []CLIFlag, cm *CLICommand) {

	if (c.BashCompletion != nil || (cm != nil && cm.BashCompletion != nil)) && !c.findFlag("generate-bash-completion", flgs) {
		flg := c.setupBashFlag(cm)
		flgs = append(flgs, flg)
	}

	if cm != nil && !c.findFlag("help", flgs) {
		flgs = append(flgs, &BoolFlg{Variable: &cm.help, Name: "help", ShortName: "h", Usage: "print commands"})
	}

	for _, f := range flgs {
		if cm != nil {
			f.GCommand(cm.Name)
		}
		f.BuildFlag(flgSet)
	}
}

func (c *CLI) buildCmds() {
	for i, d := range c.Cmds {
		doOnError := flag.ExitOnError
		tmpCommand := flag.NewFlagSet(strings.ToLower(d.Name), doOnError)
		c.Cmds[i].FS = tmpCommand
		c.buildFlags(tmpCommand, d.Flags, d)
		for j, k := range d.SubCommands {
			tmpCommand := flag.NewFlagSet(strings.ToLower(k.Name), doOnError)
			d.SubCommands[j].FS = tmpCommand
			c.buildFlags(tmpCommand, k.Flags, k)
		}
	}
}

func (c *CLI) flagSetUsage() {

	var byt bytes.Buffer
	byt.WriteString("Usage of ")
	byt.WriteString(c.cur.Name)
	byt.WriteString(":\t")
	byt.WriteString("("+c.cur.Usage+")\n")
	for i,sc := range c.cur.SubCommands {
		if i > 0 {
			byt.WriteString(",")
		}
		byt.WriteString("  "+sc.Name)
	}

	for _, f := range c.cur.Flags {
		var s string
		if len(strings.TrimSpace(f.GShortName())) > 0 {
			s = fmt.Sprintf("      -%s, -%s", f.GName(), f.GShortName())
		} else {
			s = fmt.Sprintf("      -%s", f.GName())
		}

		name := f.UnquotedUsage()
		if len(name) > 0 {
			s += "  " + name
		}
		if tmp := fmt.Sprintf("%v", f.GOptions()); len(tmp) > 2 {
			s += fmt.Sprintf("\n    \tOptions: %s", tmp)
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if len(s) <= 4 { // space, space, '-', 'x'.
			s += "\t"
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			s += "\n    \t"
		}
		s += strings.Replace(f.GUsage(), "\n", "\n    \t", -1)
		switch v := f.GValue().(type) {
		case string:
			if len(v) > 0 {
				s += fmt.Sprintf(" (default %v)", f.GValue())
			}
		default:
			s += fmt.Sprintf(" (default %v)", f.GValue())
		}
		s += "\n"
		byt.WriteString(s)
	}
	fmt.Println(byt.String())
}
func (c *CLI) printUsage() {
	szMin := 28
	szMax := 29
	flag.Usage = func() {
		var byt bytes.Buffer
		byt.WriteString("NAME:\n")
		name := "myapp"
		// if there is no path to the name use it, i.e. it's installed
		if strings.Index(os.Args[0], string(filepath.Separator)) < 0 {
			name = os.Args[0]
		} else {
			name = filepath.Base(os.Args[0])
		}
		byt.WriteString("  ")
		byt.WriteString(name)
		byt.WriteString("\n\n")

		//USAGE
		byt.WriteString("USAGE:\n")
		byt.WriteString(fmt.Sprintf("  %s [global options] command [command options] [arguments...]\n\n", name))

		if len(c.Flgs) > 0 {
			byt.WriteString("GLOBAL OPTIONS:\n")
			for _, f := range c.Flgs {
				if f.GHidden() {
					continue
				}
				var s string
				if len(f.GShortName()) > 0 {
					s = fmt.Sprintf("  -%s, -%s", f.GName(), f.GShortName())
				} else {
					s = fmt.Sprintf("  -%s", f.GName())
				}

				typeName := f.UnquotedUsage()
				if len(typeName) > 0 {
					s += "  " + typeName
				}
				if f.GRequired() {
					if len(s) < szMin {
						s += "\t\t(REQUIRED_FLAG)"
					} else if len(s) <= szMax {
						s += "\t(REQUIRED_FLAG)"
					} else {
						s += "  (REQUIRED_FLAG)"
					}

				}
				if len(f.GEnvVar()) > 0 {
					s += fmt.Sprintf("\n    %s\t(as environment var)", f.GEnvVar())
				}
				if tmp := fmt.Sprintf("%v", f.GOptions()); len(tmp) > 2 {
					s += fmt.Sprintf("\n    \tOptions: %s", tmp)
				}
				// Boolean flags of one ASCII letter are so common we
				// treat them specially, putting their usage on the same line.
				if len(s) <= 4 { // space, space, '-', 'x'.
					s += "\t"
				} else {
					// Four spaces before the tab triggers good alignment
					// for both 4- and 8-space tab stops.
					s += "\n    \t"
				}
				s += strings.Replace(f.GUsage(), "\n", "\n    \t", -1)
				tmp := fmt.Sprintf("%v", f.GValue())
				if len(tmp) > 0 {
					s += fmt.Sprintf(" (default %s", tmp)
					s += ")"
				}
				s += "\n\n"
				byt.WriteString(s)
			}
			fmt.Println(byt.String())
		}

		if len(c.Cmds) > 0 {
			byt.Reset()
			byt.WriteString("COMMANDS:\n")
			for _, d := range c.Cmds {
				if d.Hidden {
					continue
				}
				byt.WriteString(fmt.Sprintf("  %s", strings.ToLower(d.Name)))
				if len(d.Usage) > 0 {
					byt.WriteString(fmt.Sprintf(":    (%s)\n", strings.ToLower(d.Usage)))
				}
				for _, f := range d.Flags {
					if f.GHidden() {
						continue
					}
					var s string
					if len(strings.TrimSpace(f.GShortName())) > 0 {
						s = fmt.Sprintf("      -%s, -%s", f.GName(), f.GShortName())
					} else {
						s = fmt.Sprintf("      -%s", f.GName())
					}

					typeName := f.UnquotedUsage()
					if len(typeName) > 0 {
						s += "  " + typeName
					}
					if f.GRequired() {
						if len(s) < szMin {
							s += "\t\t(REQUIRED_FLAG)"
						} else if len(s) <= szMax {
							s += "\t(REQUIRED_FLAG)"
						} else {
							s += "  (REQUIRED_FLAG)"
						}
					}
					if len(f.GEnvVar()) > 0 {
						s += fmt.Sprintf("\n        %s\t(as environment var)", f.GEnvVar())
					}
					if tmp := fmt.Sprintf("%v", f.GOptions()); len(tmp) > 2 {
						s += fmt.Sprintf("\n        \tOptions: %s", tmp)
					}
					// Boolean flags of one ASCII letter are so common we
					// treat them specially, putting their usage on the same line.
					if len(s) <= 4 { // space, space, '-', 'x'.
						s += "\t"
					} else {
						// Four spaces before the tab triggers good alignment
						// for both 4- and 8-space tab stops.
						s += "\n    \t  "
					}
					s += strings.Replace(f.GUsage(), "\n", "\n    \t", -1)
					tmp := fmt.Sprintf("%v", f.GValue())
					if len(tmp) > 0 {
						s += fmt.Sprintf(" (default %s", tmp)
						s += ")"
					}
					s += "\n\n"
					byt.WriteString(s)
				}

				// show subcommands here
				for i, k := range d.SubCommands {
					if i == 0 {
						byt.WriteString("    \n    Sub Commands:\n")
					}
					if k.Hidden {
						continue
					}
					byt.WriteString(fmt.Sprintf("      %s :\t%s\n", strings.ToLower(k.Name), strings.ToLower(k.Usage)))

					for _, f := range k.Flags {
						if f.GHidden() {
							continue
						}
						var s string
						if len(strings.TrimSpace(f.GShortName())) > 0 {
							s = fmt.Sprintf("        -%s, -%s", f.GName(), f.GShortName())
						} else {
							s = fmt.Sprintf("        -%s", f.GName())
						}

						typeName := f.UnquotedUsage()
						if len(typeName) > 0 {
							s += "  " + typeName
						}
						if f.GRequired() {
							if len(s) < szMin {
								s += "\t\t(REQUIRED_FLAG)"
							} else if len(s) <= szMax {
								s += "\t(REQUIRED_FLAG)"
							} else {
								s += "  (REQUIRED_FLAG)"
							}
						}
						if len(f.GEnvVar()) > 0 {
							s += fmt.Sprintf("\n          %s\t(as environment var)", f.GEnvVar())
						}
						if tmp := fmt.Sprintf("%v", f.GOptions()); len(tmp) > 2 {
							s += fmt.Sprintf("\n    \t    Options: %s", tmp)
						}
						// Boolean flags of one ASCII letter are so common we
						// treat them specially, putting their usage on the same line.
						if len(s) <= 4 { // space, space, '-', 'x'.
							s += "\t"
						} else {
							// Four spaces before the tab triggers good alignment
							// for both 4- and 8-space tab stops.
							s += "\n    \t    "
						}
						s += strings.Replace(f.GUsage(), "\n", "\n    \t", -1)
						tmp := fmt.Sprintf("%v", f.GValue())
						if len(tmp) > 0 {
							s += fmt.Sprintf(" (default %s", tmp)
							s += ")"
						}
						s += "\n\n"
						byt.WriteString(s)
					}
				}
			}
			fmt.Println(byt.String())
		}
	}
	flag.Usage()

}

// ResetForTesting clears all flag state and sets the usage function as directed.
// After calling ResetForTesting, parse errors in flag handling will not
// exit the program.
func ResetForTesting(usage func()) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.Usage = commandLineUsage
	flag.Usage = usage
}
func commandLineUsage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func Err(err error) bool {
	if err != nil {
		return true
	}
	return false
}
func PanicErr(err error) {
	if err != nil {
		log.Fatalf("error(s)\n%+v", err)
	}
}
