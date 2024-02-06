package mycli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Uint64Flg struct {
	baseFlag
	Variable      interface{}
	Name          string
	ShortName     string
	Usage         string
	EnvVar        string
	EnvVarExclude bool
	Value         uint64
	CommaSepVal   bool
	Required      bool
	Action        interface{}
	Options       []uint64
	Hidden        bool
	debug         bool
	debugLevel    int64
}

func (c *Uint64Flg) AdjustValue(cmd string, flgValues map[string]interface{}) {
	for k, v := range flgValues {
		if k == cmd+"_"+c.Name {
			fld := c.Variable.(*uint64)
			*fld = v.(uint64)
			//fmt.Printf("value set for %v to '%v'\n", c.Name, v.(uint64))
		}
	}
}

func (c *Uint64Flg) BuildFlag(flgSet *flag.FlagSet, varMap map[string][]FieldPtr, flgValues map[string]interface{}) {
	// obtain variable field pointer
	fld := c.Variable.(*uint64)
	// set value to variable pointer using golang std lib with the passed in command line name
	flgSet.Uint64Var(fld, c.Name, c.Value, c.Usage)
	if len(c.ShortName) > 0 {
		// set value to variable using golang std lib with the passed in command line short name
		flgSet.Uint64Var(fld, c.ShortName, c.Value, c.Usage)
	}
	// set value to memory pointer of variable
	*fld = c.Value
	flgValues[c.Command+"_"+c.Name] = *fld
	// Map Any Duplicate Pointer issues for Variables and warn user
	if v, ok := varMap[fmt.Sprintf("%p", c.Variable)]; ok {
		// Don't add same thing twice
		if v[0].FieldName != c.Name || v[0].Command != c.Command {
			// found add to array
			v = append(v, FieldPtr{FieldName: c.Name, Command: c.Command, Address: fmt.Sprintf("%p", c.Variable), Value: *fld, ValType: "uint64"})
			varMap[fmt.Sprintf("%p", c.Variable)] = v
		}
	} else {
		// create array
		t := make([]FieldPtr, 0)
		t = append(t, FieldPtr{FieldName: c.Name, Command: c.Command, Address: fmt.Sprintf("%p", c.Variable), Value: *fld, ValType: "uint64"})
		varMap[fmt.Sprintf("%p", c.Variable)] = t
	}
}
func (c *Uint64Flg) GCommand(cmd string) {
	c.Command = cmd
}
func (c *Uint64Flg) GVariable() interface{} {
	return c.Variable
}
func (c *Uint64Flg) GVariableToString() string {
	return strconv.FormatUint(*c.Variable.(*uint64), 10)
}
func (c *Uint64Flg) SetEnvVar(envVar string) {
	c.EnvVar = envVar
}
func (c *Uint64Flg) GName() string {
	return c.Name
}
func (c *Uint64Flg) GShortName() string {
	return c.ShortName
}
func (c *Uint64Flg) GUsage() string {
	return c.Usage
}
func (c *Uint64Flg) GEnvVar() string {
	return c.EnvVar
}
func (c *Uint64Flg) GEnvVarExclude() bool {
	return c.EnvVarExclude
}
func (c *Uint64Flg) GValue() interface{} {
	return c.Value
}
func (c *Uint64Flg) GRequired() bool {
	return c.Required
}
func (c *Uint64Flg) GAction() interface{} {
	return c.Action
}
func (c *Uint64Flg) GOptions() interface{} {
	return c.Options
}
func (c *Uint64Flg) RetrieveEnvValue() error {
	fld := c.Variable.(*uint64)
	if *fld == c.Value {
		if envVal, found := os.LookupEnv(c.EnvVar); found {
			if c.debug {
				log.Println("overriding " + c.Name + " with env variable setting '" + envVal + "'")
			}
			var err error
			*fld, err = strconv.ParseUint(envVal, 10, 64)
			if err != nil {
				if strings.Index(err.Error(), "invalid syntax") > -1 {
					return fmt.Errorf("invalid value for '%s' flag value: '%s'", c.Name, envVal)
				}
				return err
			}
		}
	}
	return nil
}
func (c *Uint64Flg) RetrieveConfigValue(val *TomlWrapper, name string) error {
	var curVal uint64
	//name := c.Command + "." + c.Name
	//if len(c.Command) == 0 {
	//	name = c.Name
	//}
	curVal = val.Get(name).(uint64)
	fld := c.Variable.(*uint64)
	if *fld == c.Value {
		if c.debug {
			log.Println("overriding " + c.Name + " with CONFIG variable setting'" + strconv.FormatUint(curVal, 10) + "'")
		}
		*fld = curVal
	}
	return nil
}
func (c *Uint64Flg) RequiredAndNotSet() bool {
	fld := c.Variable.(*uint64)
	// if this is the same it wasn't set
	if c.Required && *fld == c.Value {
		return true
	}
	return false
}
func (c *Uint64Flg) GCommaSepVal() bool {
	if c.CommaSepVal {
		return true
	}
	return false
}
func (c *Uint64Flg) ValidValue() bool {
	// if passed in and has options then validate value is in options
	if len(c.Options) > 0 && c.Value != *c.Variable.(*uint64) {
		for _, d := range c.Options {
			if d == *c.Variable.(*uint64) {
				return true
			}
		}
		return false
	}
	return true
}
func (c *Uint64Flg) ValueAsString() string {
	return strconv.FormatUint(*c.Variable.(*uint64), 10)
}

// Kind check if this is NOT of type pointer or Nil and return error
func (c *Uint64Flg) Kind() error {
	rv := reflect.ValueOf(c)
	if rv.Kind() != reflect.Ptr {
		name := rv.FieldByName("Name").String()
		return &InvalidObjectError{reflect.TypeOf(c), "'" + name + "' flag of type"}
	} else if rv.IsNil() {
		return &InvalidObjectError{reflect.TypeOf(c), ""}
	}
	return nil
}
func (c *Uint64Flg) GHidden() bool {
	return c.Hidden
}
func (c *Uint64Flg) SetDebug(dbg bool) {
	c.debug = dbg
}
func (c *Uint64Flg) SetDebugLevel(lvl int64) {
	c.debugLevel = lvl
}
func (c *Uint64Flg) UnquotedUsage() string {
	return "uint"
}
