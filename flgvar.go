package mycli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/pelletier/go-toml"
)

type VarFlg struct {
	baseFlag
	Variable      interface{}
	Name          string
	ShortName     string
	Usage         string
	EnvVar        string
	EnvVarExclude bool
	Value         StringList
	CommaSepVal   bool
	Required      bool
	Action        interface{}
	Options       []StringList
	Hidden        bool
	debug         bool
}

func (c *VarFlg) BuildFlag(flgSet *flag.FlagSet) {
	fld := c.Variable.(*StringList)
	flgSet.Var(fld, c.Name, c.Usage)
	if len(c.ShortName) > 0 {
		flgSet.Var(fld, c.ShortName, c.Usage)
	}
	*fld = c.Value
}
func (c *VarFlg) GCommand(cmd string) {
	c.Command = cmd
}
func (c *VarFlg) GVariable() interface{} {
	return c.Variable
}
func (c *VarFlg) GVariableToString() string {
	return (*c.Variable.(*StringList)).String()
}
func (c *VarFlg) SetEnvVar(envVar string) {
	c.EnvVar = envVar
}
func (c *VarFlg) GName() string {
	return c.Name
}
func (c *VarFlg) GShortName() string {
	return c.ShortName
}
func (c *VarFlg) GUsage() string {
	return c.Usage
}
func (c *VarFlg) GEnvVar() string {
	return c.EnvVar
}
func (c *VarFlg) GEnvVarExclude() bool {
	return c.EnvVarExclude
}
func (c *VarFlg) GValue() interface{} {
	return c.Value
}
func (c *VarFlg) GRequired() bool {
	return c.Required
}
func (c *VarFlg) GAction() interface{} {
	return c.Action
}
func (c *VarFlg) GOptions() interface{} {
	return c.Options
}
func (c *VarFlg) RetrieveEnvValue() error {
	fld := c.Variable.(*StringList)
	if reflect.DeepEqual(*fld, c.Value) {
		if envVal, found := os.LookupEnv(c.EnvVar); found {
			if c.debug {
				log.Println("overriding " + c.Name + " with env variable setting '" + envVal + "'")
			}
			s := new(StringList)
			s.Set(envVal)
			*fld = *s
		}
	}
	return nil
}
func (c *VarFlg) RetrieveConfigValue(val interface{}) error {
	var curVal StringList
	name := c.Command + "." + c.Name
	if len(c.Command) == 0 {
		name = c.Name
	}
	switch val.(type) {
	case *toml.Tree:
		curVal = val.(*toml.Tree).Get(name).(StringList)
	}
	fld := c.Variable.(*StringList)
	if fld.String() == c.Value.String() {
		if c.debug {
			log.Println("overriding " + c.Name + " with CONFIG variable setting'" + curVal.String() + "'")
		}
		*fld = curVal
	}
	return nil
}
func (c *VarFlg) RequiredAndNotSet() bool {
	fld := c.Variable.(*StringList)
	if c.Required && reflect.DeepEqual(*fld, c.GValue().(StringList)) {
		return true
	}
	return false
}
func (c *VarFlg) GCommaSepVal() bool {
	if c.CommaSepVal {
		return true
	}
	return false
}
func (c *VarFlg) ValidValue() bool {
	if len(c.Options) > 0 && len(c.Variable.(*StringList).String()) > 0 && c.Value.String() != c.Variable.(*StringList).String() {
		for _, d := range c.Options {
			if d.String() == c.Value.String() {
				return true
			}
		}
		return false
	}
	return true
}
func (c *VarFlg) ValueAsString() string {
	return c.Variable.(*StringList).String()
}

// Kind check if this is NOT of type pointer or Nil and return error
func (c *VarFlg) Kind() error {
	rv := reflect.ValueOf(c)
	if rv.Kind() != reflect.Ptr {
		name := rv.FieldByName("Name").String()
		return &InvalidObjectError{reflect.TypeOf(c), "'" + name + "' flag of type"}
	} else if rv.IsNil() {
		return &InvalidObjectError{reflect.TypeOf(c), ""}
	}
	return nil
}

func (c *VarFlg) GHidden() bool {
	return c.Hidden
}
func (c *VarFlg) SetDebug(dbg bool) {
	c.debug = dbg
}
func (c *VarFlg) UnquotedUsage() string {
	return ""
}

type StringList []string

// Implement the flag.Value interface
func (s *StringList) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *StringList) Set(value string) error {
	*s = strings.Split(value, ",")
	log.Printf("%v", *s)
	return nil
}
func (c *StringList) UnquotedUsage() string {
	return "string"
}
