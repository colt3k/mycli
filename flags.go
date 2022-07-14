package mycli

import (
	"flag"
	"fmt"
	"strings"
)

type baseFlag struct {
	Command string
}
type CLIFlag interface {
	GAction() interface{}
	GEnvVar() string
	GEnvVarExclude() bool
	GHidden() bool
	GName() string
	GOptions() interface{}
	GRequired() bool
	GShortName() string
	GUsage() string
	GValue() interface{}
	GVariable() interface{}
	GVariableToString() string
	Kind() error
	RetrieveEnvValue() error
	RetrieveConfigValue(val *TomlWrapper, name string) error
	RequiredAndNotSet() bool
	SetDebug(bool)
	SetDebugLevel(int64)
	SetEnvVar(string)
	ValidValue() bool
	GCommaSepVal() bool
	ValueAsString() string
	GCommand(string)
	BuildFlag(*flag.FlagSet, map[string][]FieldPtr)
	UnquotedUsage() string
}

// checkRequired check if a required flag is set before continuing
func (c *CLI) checkRequired(subCmd string, glblFlgs []CLIFlag) {
	for _, f := range glblFlgs {
		if f.RequiredAndNotSet() {
			c.requiredMessaging(subCmd, f)
		}
	}
}

func (c *CLI) requiredMessaging(subCmd string, f CLIFlag) {

	if len(strings.TrimSpace(subCmd)) == 0 {
		c.fatalAdapter.PrintNotice(f.GName())
	} else {
		c.fatalAdapter.PrintNoticeSubCmd(f.GName(), subCmd)
	}
	fmt.Println()
}

// retrieveEnvVal if not set on commandline pull from environment
func (c *CLI) retrieveEnvVal(glblFlgs []CLIFlag) error {
	for _, f := range glblFlgs {
		err := f.RetrieveEnvValue()
		if err != nil {
			return err
		}
	}
	return nil
}
