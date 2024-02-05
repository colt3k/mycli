package mycli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

type BoolFlg struct {
	baseFlag
	Variable      interface{}
	Name          string
	ShortName     string
	Usage         string
	EnvVar        string
	EnvVarExclude bool
	Value         bool
	CommaSepVal   bool
	Required      bool
	Action        interface{}
	Options       []bool
	Hidden        bool
	debug         bool
	debugLevel    int64
}

func (c *BoolFlg) AdjustValue(cmd string, flgValues map[string]interface{}) {
	for k, v := range flgValues {
		if k == cmd+"_"+c.Name && c.Name != "debug" && c.Name != "help" && c.Name != "version" && c.Name != "generate-bash-completion" {
			fld := c.Variable.(*bool)
			*fld = v.(bool)
			//fmt.Printf("value set for %v to '%v'\n", c.Name, v.(bool))
		}
	}
}

func (c *BoolFlg) BuildFlag(flgSet *flag.FlagSet, varMap map[string][]FieldPtr, flgValues map[string]interface{}) {
	// what is this flag below? i.e. command

	fld := c.Variable.(*bool)

	flgSet.BoolVar(fld, c.Name, c.Value, c.Usage)
	if len(c.ShortName) > 0 {
		flgSet.BoolVar(fld, c.ShortName, c.Value, c.Usage)
	}
	*fld = c.Value
	flgValues[c.Command+"_"+c.Name] = *fld
	// Map Any Duplicate Pointer issues for Variables and warn user
	if v, ok := varMap[fmt.Sprintf("%p", c.Variable)]; ok {
		// Don't add same thing twice
		if v[0].FieldName != c.Name || v[0].Command != c.Command {
			// found add to array
			v = append(v, FieldPtr{FieldName: c.Name, Command: c.Command, Address: fmt.Sprintf("%p", c.Variable), Value: *fld, ValType: "bool"})
			varMap[fmt.Sprintf("%p", c.Variable)] = v
		}
	} else {
		// create array
		t := make([]FieldPtr, 0)
		t = append(t, FieldPtr{FieldName: c.Name, Command: c.Command, Address: fmt.Sprintf("%p", c.Variable), Value: *fld, ValType: "bool"})
		varMap[fmt.Sprintf("%p", c.Variable)] = t
	}
}
func (c *BoolFlg) GCommand(cmd string) {
	c.Command = cmd
}
func (c *BoolFlg) GVariable() interface{} {
	return c.Variable
}
func (c *BoolFlg) GVariableToString() string {
	return strconv.FormatBool(*c.Variable.(*bool))
}
func (c *BoolFlg) SetEnvVar(envVar string) {
	c.EnvVar = envVar
}
func (c *BoolFlg) GName() string {
	return c.Name
}
func (c *BoolFlg) GShortName() string {
	return c.ShortName
}
func (c *BoolFlg) GUsage() string {
	return c.Usage
}
func (c *BoolFlg) GEnvVar() string {
	return c.EnvVar
}
func (c *BoolFlg) GEnvVarExclude() bool {
	return c.EnvVarExclude
}
func (c *BoolFlg) GValue() interface{} {
	return c.Value
}
func (c *BoolFlg) GRequired() bool {
	return c.Required
}
func (c *BoolFlg) GAction() interface{} {
	return c.Action
}
func (c *BoolFlg) GOptions() interface{} {
	return c.Options
}
func (c *BoolFlg) RetrieveEnvValue() error {
	fld := c.Variable.(*bool)
	if *fld == c.Value {
		if envVal, found := os.LookupEnv(c.EnvVar); found {
			if c.debug {
				log.Println("overriding " + c.Name + " with env variable setting '" + envVal + "'")
			}
			var err error
			*fld, err = strconv.ParseBool(envVal)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (c *BoolFlg) RetrieveConfigValue(val *TomlWrapper, name string) error {
	var curVal bool
	//name := c.Command + "." + c.Name
	//if len(c.Command) == 0 {
	//	name = c.Name
	//}

	curVal = val.Get(name).(bool)
	fld := c.Variable.(*bool)
	if *fld == c.Value {
		if c.debug {
			log.Println("overriding " + c.Name + " with CONFIG variable setting'" + strconv.FormatBool(curVal) + "'")
		}
		*fld = curVal
	}
	return nil
}
func (c *BoolFlg) RequiredAndNotSet() bool {
	fld := c.Variable.(*bool)
	if c.Required && *fld == c.Value {
		return true
	}
	return false
}
func (c *BoolFlg) GCommaSepVal() bool {
	if c.CommaSepVal {
		return true
	}
	return false
}
func (c *BoolFlg) ValidValue() bool {
	if len(c.Options) > 0 && c.Value != *c.Variable.(*bool) {
		for _, d := range c.Options {
			if d == *c.Variable.(*bool) {
				return true
			}
		}
		return false
	}
	return true
}
func (c *BoolFlg) ValueAsString() string {
	return strconv.FormatBool(*c.Variable.(*bool))
}

// Kind check if this is NOT of type pointer or Nil and return error
func (c *BoolFlg) Kind() error {
	rv := reflect.ValueOf(c)
	if rv.Kind() != reflect.Ptr {
		name := rv.FieldByName("Name").String()
		return &InvalidObjectError{reflect.TypeOf(c), "'" + name + "' flag of type"}
	} else if rv.IsNil() {
		return &InvalidObjectError{reflect.TypeOf(c), ""}
	}
	return nil
}

func (c *BoolFlg) GHidden() bool {
	return c.Hidden
}
func (c *BoolFlg) SetDebug(dbg bool) {
	c.debug = dbg
}
func (c *BoolFlg) SetDebugLevel(lvl int64) {
	c.debugLevel = lvl
}
func (c *BoolFlg) UnquotedUsage() string {
	return ""
}
