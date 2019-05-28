package custom

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/pelletier/go-toml"

	"github.com/colt3k/mycli"
)

// Clients type of which to store client connections
type Clients struct {
	Hosts []Host `json:"clients"`
}

// Host connection information for a host
type Host struct {
	Name       string     `json:"name"`
	Connection Connection `json:"connection"`
	Cert       Cert       `json:"cert"`
}

// Connection data used by a Host
type Connection struct {
	User     string `json:"user"`
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

// Cert ppk to be used for auth
type Cert struct {
	CertPath string `json:"certpath"`
}

func (h *Host) String() string {
	return fmt.Sprintf("%s:%d", h.Connection.Host, h.Connection.Port)
}

func (c *Clients) String() string {
	var byt bytes.Buffer
	for _, d := range c.Hosts {
		byt.WriteString(fmt.Sprintf("%s:%d", d.Connection.Host, d.Connection.Port))
	}
	return byt.String()
}

// Set used to set a value as Client
func (c *Clients) Set(value string) error {
	// this would be set as json, marshal into Client object
	err := json.Unmarshal([]byte(value), c)
	if err != nil {
		log.Fatalf("error unmarshalling\n%+v", err)
	}

	return nil
}

// TomlFlg flag type to store TOML config data
type TomlFlg struct {
	Variable      interface{}
	Name          string
	ShortName     string
	Usage         string
	EnvVar        string
	EnvVarExclude bool
	Value         Clients
	CommaSepVal   bool
	Required      bool
	Action        interface{}
	Options       []Clients
	Hidden        bool
	debug         bool
	Command       string
}

// BuildFlag build flag for a flagset
func (c *TomlFlg) BuildFlag(flgSet *flag.FlagSet) {
	fld := c.Variable.(*Clients)
	flgSet.Var(fld, c.Name, c.Usage)
	if len(c.ShortName) > 0 {
		flgSet.Var(fld, c.ShortName, c.Usage)
	}
	*fld = c.Value
}

// GCommand set command for flag
func (c *TomlFlg) GCommand(cmd string) {
	c.Command = cmd
}

// GVariable set variable for flag
func (c *TomlFlg) GVariable() interface{} {
	return c.Variable
}

// GVariableToString convert value of variable to string
func (c *TomlFlg) GVariableToString() string {
	return (*c.Variable.(*Clients)).String()
}

// SetEnvVar set environment variable for flag
func (c *TomlFlg) SetEnvVar(envVar string) {
	c.EnvVar = envVar
}

// GName set name for flag
func (c *TomlFlg) GName() string {
	return c.Name
}

// GShortName set short name for flag
func (c *TomlFlg) GShortName() string {
	return c.ShortName
}

// GUsage set usage for flag
func (c *TomlFlg) GUsage() string {
	return c.Usage
}

// GEnvVar get environment variable for flag
func (c *TomlFlg) GEnvVar() string {
	return c.EnvVar
}

// GEnvVarExclude get environment variable exclude for flag
func (c *TomlFlg) GEnvVarExclude() bool {
	return c.EnvVarExclude
}

// GValue get value for flag
func (c *TomlFlg) GValue() interface{} {
	if c.Value.Hosts == nil {
		return ""
	}
	return c.Value
}

// UnquotedUsage get unquoted usage value for flag
func (c *TomlFlg) UnquotedUsage() string {
	return "custom.Clients"
}

// GRequired get required value for flag
func (c *TomlFlg) GRequired() bool {
	return c.Required
}

// GAction get action for flag
func (c *TomlFlg) GAction() interface{} {
	return c.Action
}

// GOptions get options for flag
func (c *TomlFlg) GOptions() interface{} {
	return c.Options
}

// RetrieveEnvValue get environment value for flag
func (c *TomlFlg) RetrieveEnvValue() error {
	// pull from environment variable and convert to store as client object
	fld := c.Variable.(*Clients)
	if reflect.DeepEqual(*fld, c.Value) {
		if envVal, found := os.LookupEnv(c.EnvVar); found {
			if c.debug {
				log.Println("overriding " + c.Name + " with env variable setting '" + envVal + "'")
			}
			s := new(Clients)
			// convert value to
			s.Set(envVal)
			*fld = *s
		}
	}
	return nil
}

// RetrieveConfigValue get config value for flag
func (c *TomlFlg) RetrieveConfigValue(val interface{}) error {
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
		log.Fatalf("error unmarshalling\n%+v", err)
	}

	return nil
}

// RequiredAndNotSet determine if this is required and not set for flag
func (c *TomlFlg) RequiredAndNotSet() bool {
	fld := c.Variable.(*Clients)
	if c.Required && reflect.DeepEqual(*fld, c.GValue().(Clients)) {
		return true
	}
	return false
}
func (c *TomlFlg) GCommaSepVal() bool {
	if c.CommaSepVal {
		return true
	}
	return false
}

// ValidValue determine if this is a valid value for flag
func (c *TomlFlg) ValidValue() bool {
	if len(c.Options) > 0 && len(c.Variable.(*Clients).String()) > 0 && c.Value.String() != c.Variable.(*Clients).String() {
		for _, d := range c.Options {
			if d.String() == c.Value.String() {
				return true
			}
		}
		return false
	}
	return true
}

// ValueAsString get value in string format for flag
func (c *TomlFlg) ValueAsString() string {
	return c.Variable.(*Clients).String()
}

// Kind check if this is NOT of type pointer or Nil and return error
func (c *TomlFlg) Kind() error {
	rv := reflect.ValueOf(c)
	if rv.Kind() != reflect.Ptr {
		name := rv.FieldByName("Name").String()
		return &mycli.InvalidObjectError{Type: reflect.TypeOf(c), Name: "'" + name + "' flag of type"}
	} else if rv.IsNil() {
		return &mycli.InvalidObjectError{Type: reflect.TypeOf(c), Name: ""}
	}
	return nil
}

// GHidden get hidden property for flag
func (c *TomlFlg) GHidden() bool {
	return c.Hidden
}

// SetDebug set debug property for flag
func (c *TomlFlg) SetDebug(dbg bool) {
	c.debug = dbg
}
