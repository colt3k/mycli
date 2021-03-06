package mycli

import (
	"flag"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/pelletier/go-toml"
)

type Float64Flg struct {
	baseFlag
	Variable      interface{}
	Name          string
	ShortName     string
	Usage         string
	EnvVar        string
	EnvVarExclude bool
	Value         float64
	CommaSepVal   bool
	Required      bool
	Action        interface{}
	Options       []float64
	Hidden        bool
	debug         bool
}

func (c *Float64Flg) BuildFlag(flgSet *flag.FlagSet) {
	fld := c.Variable.(*float64)
	flgSet.Float64Var(fld, c.Name, c.Value, c.Usage)
	if len(c.ShortName) > 0 {
		flgSet.Float64Var(fld, c.ShortName, c.Value, c.Usage)
	}
	*fld = c.Value
}
func (c *Float64Flg) GCommand(cmd string) {
	c.Command = cmd
}
func (c *Float64Flg) GVariable() interface{} {
	return c.Variable
}
func (c *Float64Flg) GVariableToString() string {
	return strconv.FormatFloat(*c.Variable.(*float64), 'f', -1, 64)
}
func (c *Float64Flg) SetEnvVar(envVar string) {
	c.EnvVar = envVar
}
func (c *Float64Flg) GName() string {
	return c.Name
}
func (c *Float64Flg) GShortName() string {
	return c.ShortName
}
func (c *Float64Flg) GUsage() string {
	return c.Usage
}
func (c *Float64Flg) GEnvVar() string {
	return c.EnvVar
}
func (c *Float64Flg) GEnvVarExclude() bool {
	return c.EnvVarExclude
}
func (c *Float64Flg) GValue() interface{} {
	return c.Value
}
func (c *Float64Flg) GRequired() bool {
	return c.Required
}
func (c *Float64Flg) GAction() interface{} {
	return c.Action
}
func (c *Float64Flg) GOptions() interface{} {
	return c.Options
}
func (c *Float64Flg) RetrieveEnvValue() error {
	fld := c.Variable.(*float64)
	if *fld == c.Value {
		if envVal, found := os.LookupEnv(c.EnvVar); found {
			if c.debug {
				log.Println("overriding " + c.Name + " with env variable setting '" + envVal + "'")
			}
			var err error
			*fld, err = strconv.ParseFloat(envVal, 64)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (c *Float64Flg) RetrieveConfigValue(val interface{}) error {
	var curVal float64
	name := c.Command + "." + c.Name
	if len(c.Command) == 0 {
		name = c.Name
	}
	switch val.(type) {
	case *toml.Tree:
		curVal = val.(*toml.Tree).Get(name).(float64)
	}
	fld := c.Variable.(*float64)
	if *fld == c.Value {
		if c.debug {
			log.Println("overriding " + c.Name + " with CONFIG variable setting'" + strconv.FormatFloat(curVal, 'E', -1, 64) + "'")
		}
		*fld = curVal
	}
	return nil
}
func (c *Float64Flg) RequiredAndNotSet() bool {
	fld := c.Variable.(*float64)
	if c.Required && *fld == c.Value {
		return true
	}
	return false
}
func (c *Float64Flg) GCommaSepVal() bool {
	if c.CommaSepVal {
		return true
	}
	return false
}
func (c *Float64Flg) ValidValue() bool {

	if len(c.Options) > 0 && c.Value != *c.Variable.(*float64) {
		for _, d := range c.Options {
			if d == *c.Variable.(*float64) {
				return true
			}
		}
		return false
	}
	return true
}
func (c *Float64Flg) ValueAsString() string {
	return strconv.FormatFloat(*c.Variable.(*float64), 'E', -1, 64)
}

// Kind check if this is NOT of type pointer or Nil and return error
func (c *Float64Flg) Kind() error {
	rv := reflect.ValueOf(c)
	if rv.Kind() != reflect.Ptr {
		name := rv.FieldByName("Name").String()
		return &InvalidObjectError{reflect.TypeOf(c), "'" + name + "' flag of type"}
	} else if rv.IsNil() {
		return &InvalidObjectError{reflect.TypeOf(c), ""}
	}
	return nil
}
func (c *Float64Flg) GHidden() bool {
	return c.Hidden
}
func (c *Float64Flg) SetDebug(dbg bool) {
	c.debug = dbg
}
func (c *Float64Flg) UnquotedUsage() string {
	return "float"
}
