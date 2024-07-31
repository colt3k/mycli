package mycli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
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
	debugLevel    int64
}

func (c *VarFlg) AdjustValue(cmd string, flgValues map[string]interface{}) {
	for k, v := range flgValues {
		if k == cmd+"_"+c.Name {
			fld := c.Variable.(*StringList)
			*fld = v.(StringList)
			//fmt.Printf("value set for %v to '%v'\n", c.Name, v.(StringList))
		}
	}
}

func (c *VarFlg) BuildFlag(flgSet *flag.FlagSet, varMap map[string][]FieldPtr, flgValues map[string]interface{}) {
	// obtain variable field pointer
	fld := c.Variable.(*StringList)
	// set value to variable pointer using golang std lib with the passed in command line name
	flgSet.Var(fld, c.Name, c.Usage)
	if len(c.ShortName) > 0 {
		// set value to variable using golang std lib with the passed in command line short name
		flgSet.Var(fld, c.ShortName, c.Usage)
	}
	// set value to memory pointer of variable
	*fld = c.Value
	flgValues[c.Command+"_"+c.Name] = *fld
	// Map Any Duplicate Pointer issues for Variables and warn user
	if v, ok := varMap[fmt.Sprintf("%p", c.Variable)]; ok {
		// Don't add same thing twice
		if v[0].FieldName != c.Name || v[0].Command != c.Command {
			// found add to array
			v = append(v, FieldPtr{FieldName: c.Name, Command: c.Command, Address: fmt.Sprintf("%p", c.Variable), Value: *fld, ValType: "stringlist"})
			varMap[fmt.Sprintf("%p", c.Variable)] = v
		}
	} else {
		// create array
		t := make([]FieldPtr, 0)
		t = append(t, FieldPtr{FieldName: c.Name, Command: c.Command, Address: fmt.Sprintf("%p", c.Variable), Value: *fld, ValType: "stringlist"})
		varMap[fmt.Sprintf("%p", c.Variable)] = t
	}
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

func (c *VarFlg) RetrieveConfigValue(val *TomlWrapper, name string) error {
	var curVal StringList
	//name := c.Command + "." + c.Name
	//if len(c.Command) == 0 {
	//	name = c.Name
	//}
	curVal = val.Get(name).(StringList)
	fld := c.Variable.(*StringList)
	if fld.String() == c.Value.String() {
		if c.debug {
			log.Println("overriding " + c.Name + " with CONFIG variable setting '" + curVal.String() + "'")
		}
		*fld = curVal
	}
	return nil
}

func (c *VarFlg) RetrieveConfigValueOrig(val map[string]interface{}, name string) error {
	var curVal StringList
	//name := c.Command + "." + c.Name
	//if len(c.Command) == 0 {
	//	name = c.Name
	//}
	curVal = val[name].(StringList)
	fld := c.Variable.(*StringList)
	if fld.String() == c.Value.String() {
		if c.debug {
			log.Println("overriding " + c.Name + " with CONFIG variable setting '" + curVal.String() + "'")
		}
		*fld = curVal
	}
	return nil
}
func (c *VarFlg) RequiredAndNotSet() bool {
	fld := c.Variable.(*StringList)
	// if this is the same it wasn't set
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
		// if passed in and has options then validate value is in options
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
func (c *VarFlg) SetDebugLevel(lvl int64) {
	c.debugLevel = lvl
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
