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

type Int64Flg struct {
	baseFlag
	Variable      interface{}
	Name          string
	ShortName     string
	Usage         string
	EnvVar        string
	EnvVarExclude bool
	Value         int64
	CommaSepVal   bool
	Required      bool
	Action        interface{}
	Options       []int64
	Hidden        bool
	debug         bool
	debugLevel    int64
}

func (c *Int64Flg) AdjustValue(cmd string, flgValues map[string]interface{}) {
	for k, v := range flgValues {
		if k == cmd+"_"+c.Name && c.Name != "debugLevel" {
			fld := c.Variable.(*int64)
			*fld = v.(int64)
			//fmt.Printf("value set for %v to '%v'\n", c.Name, v.(int64))
		}
	}
}
func (c *Int64Flg) BuildFlag(flgSet *flag.FlagSet, varMap map[string][]FieldPtr, flgValues map[string]interface{}) {
	fld := c.Variable.(*int64)

	flgSet.Int64Var(fld, c.Name, c.Value, c.Usage)
	if len(c.ShortName) > 0 {
		flgSet.Int64Var(fld, c.ShortName, c.Value, c.Usage)
	}
	*fld = c.Value
	flgValues[c.Command+"_"+c.Name] = *fld
	// Map Any Duplicate Pointer issues for Variables and warn user
	if v, ok := varMap[fmt.Sprintf("%p", c.Variable)]; ok {
		// Don't add same thing twice
		if v[0].FieldName != c.Name || v[0].Command != c.Command {
			// found add to array
			v = append(v, FieldPtr{FieldName: c.Name, Command: c.Command, Address: fmt.Sprintf("%p", c.Variable), Value: *fld, ValType: "int64"})
			varMap[fmt.Sprintf("%p", c.Variable)] = v
		}
	} else {
		// create array
		t := make([]FieldPtr, 0)
		t = append(t, FieldPtr{FieldName: c.Name, Command: c.Command, Address: fmt.Sprintf("%p", c.Variable), Value: *fld, ValType: "int64"})
		varMap[fmt.Sprintf("%p", c.Variable)] = t
	}
}
func (c *Int64Flg) GCommand(cmd string) {
	c.Command = cmd
}
func (c *Int64Flg) GVariable() interface{} {
	return c.Variable
}
func (c *Int64Flg) GVariableToString() string {
	return strconv.FormatInt(*c.Variable.(*int64), 10)
}
func (c *Int64Flg) SetEnvVar(envVar string) {
	c.EnvVar = envVar
}
func (c *Int64Flg) GName() string {
	return c.Name
}
func (c *Int64Flg) GShortName() string {
	return c.ShortName
}
func (c *Int64Flg) GUsage() string {
	return c.Usage
}
func (c *Int64Flg) GEnvVar() string {
	return c.EnvVar
}
func (c *Int64Flg) GEnvVarExclude() bool {
	return c.EnvVarExclude
}
func (c *Int64Flg) GValue() interface{} {
	return c.Value
}
func (c *Int64Flg) GRequired() bool {
	return c.Required
}
func (c *Int64Flg) GAction() interface{} {
	return c.Action
}
func (c *Int64Flg) GOptions() interface{} {
	return c.Options
}
func (c *Int64Flg) RetrieveEnvValue() error {
	fld := c.Variable.(*int64)
	if *fld == c.Value {
		if envVal, found := os.LookupEnv(c.EnvVar); found {
			var err error
			if c.debug {
				//log.Println("overriding "+c.Name+" with env variable setting")
			}
			*fld, err = strconv.ParseInt(envVal, 10, 64)
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
func (c *Int64Flg) RetrieveConfigValue(val *TomlWrapper, name string) error {
	var curVal int64
	//name := c.Command+"."+c.Name
	//if len(c.Command) == 0 {
	//	name = c.Name
	//}
	curVal = val.Get(name).(int64)
	fld := c.Variable.(*int64)
	if *fld == c.Value {
		if c.debug {
			log.Println("overriding " + c.Name + " with CONFIG variable setting '" + strconv.FormatInt(curVal, 10) + "'")
		}
		*fld = curVal
	}
	return nil
}
func (c *Int64Flg) RequiredAndNotSet() bool {
	fld := c.Variable.(*int64)
	if c.Required && *fld == c.Value {
		return true
	}
	return false
}
func (c *Int64Flg) GCommaSepVal() bool {
	if c.CommaSepVal {
		return true
	}
	return false
}
func (c *Int64Flg) ValidValue() bool {
	// if passed in and has options then validate value is in options
	if len(c.Options) > 0 && c.Value != *c.Variable.(*int64) {
		for _, d := range c.Options {
			if d == *c.Variable.(*int64) {
				return true
			}
		}
		return false
	}
	return true
}
func (c *Int64Flg) ValueAsString() string {
	return strconv.FormatInt(*c.Variable.(*int64), 10)
}

// Kind check if this is NOT of type pointer or Nil and return error
func (c *Int64Flg) Kind() error {
	rv := reflect.ValueOf(c)
	if rv.Kind() != reflect.Ptr {
		name := rv.FieldByName("Name").String()
		return &InvalidObjectError{reflect.TypeOf(c), "'" + name + "' flag of type"}
	} else if rv.IsNil() {
		return &InvalidObjectError{reflect.TypeOf(c), ""}
	}
	return nil
}
func (c *Int64Flg) GHidden() bool {
	return c.Hidden
}
func (c *Int64Flg) SetDebug(dbg bool) {
	c.debug = dbg
}
func (c *Int64Flg) SetDebugLevel(lvl int64) {
	c.debugLevel = lvl
}
func (c *Int64Flg) UnquotedUsage() string {
	return "int"
}
