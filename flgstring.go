package mycli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

type StringFlg struct {
	baseFlag
	Variable      interface{}
	Name          string
	ShortName     string
	Usage         string
	EnvVar        string
	EnvVarExclude bool
	Value         string
	CommaSepVal   bool
	Required      bool
	Action        interface{}
	Options       []string
	Hidden        bool
	debug         bool
	debugLevel    int64
}

func (c *StringFlg) AdjustValue(cmd string, flgValues map[string]interface{}) {
	for k, v := range flgValues {
		//if c.Name == "application" && cmd == "weserve_config" && k == cmd+"_"+c.Name {
		//	fmt.Println("application\n")
		//}
		if k == cmd+"_"+c.Name && c.Name != "config" && c.Name != "proxyhttp" && c.Name != "proxyhttps" && c.Name != "noproxy" {
			fld := c.Variable.(*string)
			*fld = v.(string)
			//fmt.Printf("value set for %v to '%v'\n", c.Name, v.(string))
		}
	}
}

func (c *StringFlg) BuildFlag(flgSet *flag.FlagSet, varMap map[string][]FieldPtr, flgValues map[string]interface{}) {
	// obtain variable field pointer
	fld := c.Variable.(*string)
	// set value to variable pointer using golang std lib with the passed in command line name
	flgSet.StringVar(fld, c.Name, c.Value, c.Usage)
	if len(c.ShortName) > 0 {
		// set value to variable using golang std lib with the passed in command line short name
		flgSet.StringVar(fld, c.ShortName, c.Value, c.Usage)
	}
	// set value to memory pointer of variable
	*fld = c.Value
	flgValues[c.Command+"_"+c.Name] = *fld
	// Map Any Duplicate Pointer issues for Variables and warn user
	if v, ok := varMap[fmt.Sprintf("%p", c.Variable)]; ok {
		// Don't add same thing twice
		if v[0].FieldName != c.Name || v[0].Command != c.Command {
			// found add to array
			v = append(v, FieldPtr{FieldName: c.Name, Command: c.Command, Address: fmt.Sprintf("%p", c.Variable), Value: *fld, ValType: "string"})
			varMap[fmt.Sprintf("%p", c.Variable)] = v
		}
	} else {
		// create array
		t := make([]FieldPtr, 0)
		t = append(t, FieldPtr{FieldName: c.Name, Command: c.Command, Address: fmt.Sprintf("%p", c.Variable), Value: *fld, ValType: "string"})
		varMap[fmt.Sprintf("%p", c.Variable)] = t
	}
}
func (c *StringFlg) GCommand(cmd string) {
	c.Command = cmd
}
func (c *StringFlg) GVariable() interface{} {
	return c.Variable
}
func (c *StringFlg) GVariableToString() string {
	return *c.Variable.(*string)
}
func (c *StringFlg) SetEnvVar(envVar string) {
	c.EnvVar = envVar
}
func (c *StringFlg) GName() string {
	return c.Name
}
func (c *StringFlg) GShortName() string {
	return c.ShortName
}
func (c *StringFlg) GUsage() string {
	return c.Usage
}
func (c *StringFlg) GEnvVar() string {
	return c.EnvVar
}
func (c *StringFlg) GEnvVarExclude() bool {
	return c.EnvVarExclude
}
func (c *StringFlg) GValue() interface{} {
	return c.Value
}
func (c *StringFlg) GRequired() bool {
	return c.Required
}
func (c *StringFlg) GAction() interface{} {
	return c.Action
}
func (c *StringFlg) GOptions() interface{} {
	return c.Options
}
func (c *StringFlg) RetrieveEnvValue() error {
	fld := c.Variable.(*string)
	if *fld == c.Value {
		if envVal, found := os.LookupEnv(c.EnvVar); found {
			if c.debug {
				log.Println("overriding " + c.Name + " with env variable setting '" + envVal + "'")
			}
			*fld = envVal
		}
	}
	return nil
}
func (c *StringFlg) RetrieveConfigValue(val *TomlWrapper, name string) error {
	var curVal string
	//name := c.Command + "." + c.Name
	//if len(c.Command) == 0 {
	//	name = c.Name
	//}
	curVal = val.Get(name).(string)
	fld := c.Variable.(*string)
	if *fld == c.Value {
		if c.debug {
			log.Println("overriding " + c.Name + " with CONFIG variable setting '" + curVal + "'")
		}
		*fld = curVal
	}
	return nil
}
func (c *StringFlg) RequiredAndNotSet() bool {
	fld := c.Variable.(*string)
	// if this is the same it wasn't set
	if c.Required && *fld == c.Value {
		return true
	}
	return false
}
func (c *StringFlg) GCommaSepVal() bool {
	if c.CommaSepVal {
		return true
	}
	return false
}
func (c *StringFlg) ValidValue() bool {
	if len(c.Options) > 0 && len(*c.Variable.(*string)) > 0 && c.Value != *c.Variable.(*string) {
		if c.GCommaSepVal() {
			// split values on comma then compare
			t := *c.Variable.(*string)
			vals := strings.Split(t, ",")
			// count each and see if all match
			count := 0
			for _, v := range vals {
				for _, d := range c.Options {
					if v == d {
						count += 1
					}
				}
			}
			if count != len(vals) {
				return false
			}
			return true
		} else {
			// if passed in and has options then validate value is in options
			for _, d := range c.Options {
				if d == *c.Variable.(*string) {
					return true
				}
			}
		}
		return false
	}
	return true
}
func (c *StringFlg) ValueAsString() string {
	return *c.Variable.(*string)
}

// Kind check if this is NOT of type pointer or Nil and return error
func (c *StringFlg) Kind() error {
	rv := reflect.ValueOf(c)
	if rv.Kind() != reflect.Ptr {
		name := rv.FieldByName("Name").String()
		return &InvalidObjectError{reflect.TypeOf(c), "'" + name + "' flag of type"}
	} else if rv.IsNil() {
		return &InvalidObjectError{reflect.TypeOf(c), ""}
	}
	return nil
}
func (c *StringFlg) GHidden() bool {
	return c.Hidden
}
func (c *StringFlg) SetDebug(dbg bool) {
	c.debug = dbg
}
func (c *StringFlg) SetDebugLevel(lvl int64) {
	c.debugLevel = lvl
}
func (c *StringFlg) UnquotedUsage() string {
	return "string"
}
