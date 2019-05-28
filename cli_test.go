package mycli

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Tests struct {
	name  string
	tests []Test
}
type Test struct {
	name  string
	field interface{}
	value interface{}
}

func setupTestCase(t *testing.T) func(t *testing.T) {
	oldArgs := os.Args
	return func(t *testing.T) {

		ResetForTesting(nil)
		defer func() { os.Args = oldArgs }()
		cli = nil
	}
}

func setupSubTest(t *testing.T) func(t *testing.T) {
	//t.Log("setup sub test")
	return func(t *testing.T) {
		//t.Log("teardown sub test")
	}
}

var cli *CLI
var printNotice, printNoticeSubCmd, displayedUsage bool

func TestVersionPrint(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	cli = NewCli(nil, nil)
	if cli == nil {
		t.Fail()
	}
	cli.Version = "1.0"
	cli.BuildDate = "09242018"
	cli.GitCommit = "3456789uijy"
	cli.PostGlblAction = func() { fmt.Println("hello") }
	cases := []struct {
		name  string
		value string
		field interface{}
	}{
		{"version", "1.0", cli.Version},
		{"build", "09242018", cli.BuildDate},
		{"gitcommit", "3456789uijy", cli.GitCommit},
	}

	cli.Flgs = []CLIFlag{}

	cli.Cmds = []*CLICommand{}

	os.Args = []string{"cmd", "-version"}
	cli.Parse()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := setupSubTest(t)
			defer teardownSubTest(t)

			t.Log("find", tc.name, "with value", tc.value)
			assert.EqualValues(t, tc.value, tc.field)
		})
	}

}

func TestHelpPrint(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	var protocol string
	cli = NewCli(nil, nil)
	if cli == nil {
		t.Fail()
	}
	cli.TestMode = true
	cli.Flgs = []CLIFlag{}

	cli.Cmds = []*CLICommand{
		{
			Name:           "server",
			ShortName:      "s",
			Usage:          "use as a server",
			Value:          nil,
			Action:         func() { log.Println("cmd action") },
			PreAction:      func() { log.Println("cmd preaction") },
			PostAction:     func() { log.Println("cmd postaction") },
			BashCompletion: BashCompletionSub,
			Flags: []CLIFlag{
				&StringFlg{Variable: &protocol, Name: "protocol", ShortName: "proto", Usage: "Set Protocol http(s)", Value: "http"},
			},
		},
	}

	os.Args = []string{"cmd", "-h"}
	cli.Parse()

}

func TestInitialization(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	cases := []struct {
		name  string
		found bool
	}{
		{"debug", true},
		{"help", true},
		{"version", true},
	}

	cli = NewCli(nil, nil)
	if cli == nil {
		t.Fail()
	}
	cli.TestMode = true
	cli.BashCompletion = BashCompletionMain
	cli.Flgs = []CLIFlag{}
	cli.Cmds = []*CLICommand{}

	os.Args = []string{"cmd", "-h", "-generate-bash-completion"}
	cli.Parse()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := setupSubTest(t)
			defer teardownSubTest(t)

			t.Log("find", tc.name, "with value", tc.found)
			assert.EqualValues(t, tc.found, cli.findFlag(tc.name, cli.Flgs))
		})
	}

	assert.Equal(t, true, cli.Help())
}

func TestFlgTypes(t *testing.T) {

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	var (
		test  bool
		test2 bool
		test3 string
		test5 StringList
		test6 int64
		test7 float64
		test9 uint64
	)

	cli = NewCli(nil, nil)
	if cli == nil {
		t.Fail()
	}
	sl := &StringList{}
	sl.Set("hello,test")

	cli.Flgs = []CLIFlag{
		&BoolFlg{Variable: &test, Name: "test", ShortName: "t", Usage: "use by passing -test"},
		&BoolFlg{Variable: &test2, Name: "test2", ShortName: "t2", EnvVar: "TESTTWO", Usage: "use by passing -test2", Required: true},
		&StringFlg{Variable: &test3, Name: "test3", ShortName: "t3", Usage: "use by passing -test3", Value: "test val"},
		&Int64Flg{Variable: &test6, Name: "test6", ShortName: "t6", Usage: "use by passing -test4", Value: int64(1)},
		&Float64Flg{Variable: &test7, Name: "test7", ShortName: "t7", Usage: "use by passing -test4", Value: float64(1)},
		&Uint64Flg{Variable: &test9, Name: "test9", ShortName: "t9", Usage: "use by passing -test4", Value: uint64(1)},
		&VarFlg{Variable: &test5, Name: "test5", ShortName: "t5", Usage: "custom passing -test5", Value: *sl},
	}

	os.Args = []string{"cmd", "-d", "-test2"}
	cli.Parse()

	flg := cli.Flag("test", cli.Flgs)
	flg2 := cli.Flag("test2", cli.Flgs)
	flg3 := cli.Flag("test3", cli.Flgs)

	flg5 := cli.Flag("test5", cli.Flgs)
	cases := []Tests{
		{"not required bool",
			[]Test{
				{"name", flg.GName(), "test"},
				{"shortname", flg.GShortName(), "t"},
				{"env_var", flg.GEnvVar(), "TEST"},
				{"value", test, false},
				{"required", flg.GRequired(), false},
			},
		},
		{"required bool",
			[]Test{
				{"name(req)", flg2.GName(), "test2"},
				{"shortname(req)", flg2.GShortName(), "t2"},
				{"env_var(req)", flg2.GEnvVar(), "TESTTWO"},
				{"value(req)", test2, true},
				{"required(req)", flg2.GRequired(), true},
			},
		},
		{"not required string",
			[]Test{
				{"name", flg3.GName(), "test3"},
				{"shortname", flg3.GShortName(), "t3"},
				{"env_var", flg3.GEnvVar(), "TEST3"},
				{"value", test3, "test val"},
				{"required", flg3.GRequired(), false},
			},
		},
		{"not required custom",
			[]Test{
				{"name", flg5.GName(), "test5"},
				{"shortname", flg5.GShortName(), "t5"},
				{"env_var", flg5.GEnvVar(), "TEST5"},
				{"value", test5, StringList{"hello", "test"}},
				{"required", flg5.GRequired(), false},
			},
		},
	}

	for _, tc := range cases {
		for _, test := range tc.tests {
			t.Run(tc.name+" "+test.name, func(t *testing.T) {
				teardownSubTest := setupSubTest(t)
				defer teardownSubTest(t)

				t.Log("find", test.name, "with value", test.value)
				assert.EqualValues(t, test.value, test.field)
			})
		}
	}
}

func TestFlgTypesRequired(t *testing.T) {

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	var (
		test  bool
		test2 string
		test5 int64
		test6 float64
		test7 uint64
	)

	cli = NewCli(nil, nil)
	if cli == nil {
		t.Fail()
	}

	cli.Flgs = []CLIFlag{
		&BoolFlg{Variable: &test, Name: "test", ShortName: "t", Usage: "use by passing -test"},
		&StringFlg{Variable: &test2, Name: "test2", ShortName: "t2", Usage: "use by passing -test2", Value: "test val", Required: true},
		&Int64Flg{Variable: &test5, Name: "test5", ShortName: "t5", Usage: "use by passing -test5", Value: int64(-1), Required: true},
		&Float64Flg{Variable: &test6, Name: "test6", ShortName: "t6", Usage: "use by passing -test6", Value: float64(1), Required: true},
		&Uint64Flg{Variable: &test7, Name: "test7", ShortName: "t7", Usage: "use by passing -test7", Value: uint64(1), Required: true},
	}
	os.Args = []string{"cmd", "-test", "-test2=high test", "-test5=1", "-test6=1.2", "-test7=2"}
	cli.Parse()

	flg := cli.Flag("test", cli.Flgs)
	flg2 := cli.Flag("test2", cli.Flgs)
	flg5 := cli.Flag("test5", cli.Flgs)
	flg6 := cli.Flag("test6", cli.Flgs)
	flg7 := cli.Flag("test7", cli.Flgs)
	cases := []Tests{
		{"test1",
			[]Test{
				{"name", flg.GName(), "test"},
				{"value", test, true},
			},
		},
		{"test2",
			[]Test{
				{"name", flg2.GName(), "test2"},
				{"value", test2, "high test"},
			},
		},
		{"test5",
			[]Test{
				{"name", flg5.GName(), "test5"},
				{"value", test5, 1},
			},
		},
		{"test6",
			[]Test{
				{"name", flg6.GName(), "test6"},
				{"value", test6, 1.2},
			},
		},
		{"test7",
			[]Test{
				{"name", flg7.GName(), "test7"},
				{"value", test7, 2},
			},
		},
	}

	for _, tc := range cases {
		for _, test := range tc.tests {
			t.Run(tc.name+" "+test.name, func(t *testing.T) {
				teardownSubTest := setupSubTest(t)
				defer teardownSubTest(t)

				t.Log("find", test.name, "with value", test.value)
				assert.EqualValues(t, test.value, test.field)
			})
		}
	}

}

func TestFlgTypesRequiredFail(t *testing.T) {

	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	var (
		test bool
	)

	cli = NewCli(new(FatalStub), nil)
	cli.TestMode = true
	cli.Flgs = []CLIFlag{
		&BoolFlg{Variable: &test, Name: "test", ShortName: "t", Usage: "use by passing -test", Required: true},
	}

	os.Args = []string{"cmd", "none"}
	cli.Parse()

	flg := cli.Flag("test", cli.Flgs)
	cases := []Tests{
		{"required bool",
			[]Test{
				{"required(req)", flg.GRequired(), true},
				{"parameter NOT passed", printNotice, true},
			},
		},
	}

	for _, tc := range cases {
		for _, test := range tc.tests {
			t.Run(tc.name+" "+test.name, func(t *testing.T) {
				teardownSubTest := setupSubTest(t)
				defer teardownSubTest(t)

				t.Log("find", test.name, "with value", test.value)
				assert.EqualValues(t, test.value, test.field)
			})
		}
	}

}

func TestFlgPfx(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	var (
		test, test2 bool
		test3       string
		test5       StringList
	)

	cli = NewCli(nil, nil)
	if cli == nil {
		t.Fail()
	}
	cli.EnvPrefix = "TST"
	cli.Flgs = []CLIFlag{
		&BoolFlg{Variable: &test, Name: "test", ShortName: "t", Usage: "use by passing -test"},
		&BoolFlg{Variable: &test2, Name: "test2", ShortName: "t2", EnvVar: "TESTTWO", Usage: "use by passing -test2", Required: true},
		&StringFlg{Variable: &test3, Name: "test3", ShortName: "t3", Usage: "use by passing -test3", Value: "test val"},
		&VarFlg{Variable: &test5, Name: "test5", ShortName: "t5", Usage: "custom passing -test5", Value: StringList{"hello", "test"}},
	}

	os.Args = []string{"cmd", "-test2"}
	cli.Parse()

	flg := cli.Flag("test", cli.Flgs)
	flg2 := cli.Flag("test2", cli.Flgs)
	flg3 := cli.Flag("test3", cli.Flgs)
	flg5 := cli.Flag("test5", cli.Flgs)
	cases := []Tests{
		{"not required bool",
			[]Test{
				{"name", flg.GName(), "test"},
				{"shortname", flg.GShortName(), "t"},
				{"env_var", flg.GEnvVar(), "TST_TEST"},
				{"value", test, false},
				{"required", flg.GRequired(), false},
			},
		},
		{"required bool",
			[]Test{
				{"name(req)", flg2.GName(), "test2"},
				{"shortname(req)", flg2.GShortName(), "t2"},
				{"env_var(req)", flg2.GEnvVar(), "TST_TESTTWO"},
				{"value(req)", test2, true},
				{"required(req)", flg2.GRequired(), true},
			},
		},
		{"not required string",
			[]Test{
				{"name", flg3.GName(), "test3"},
				{"shortname", flg3.GShortName(), "t3"},
				{"env_var", flg3.GEnvVar(), "TST_TEST3"},
				{"value", test3, "test val"},
				{"required", flg3.GRequired(), false},
			},
		},

		{"not required custom",
			[]Test{
				{"name", flg5.GName(), "test5"},
				{"shortname", flg5.GShortName(), "t5"},
				{"env_var", flg5.GEnvVar(), "TST_TEST5"},
				{"value", test5, StringList{"hello", "test"}},
				{"required", flg5.GRequired(), false},
			},
		},
	}

	for _, tc := range cases {
		for _, test := range tc.tests {
			t.Run(tc.name+" "+test.name, func(t *testing.T) {
				teardownSubTest := setupSubTest(t)
				defer teardownSubTest(t)

				t.Log("find", test.name, "with value", test.value)
				assert.EqualValues(t, test.value, test.field)
			})
		}
	}
}

func TestCmd(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	var (
		test     bool
		t2       int64
		protocol string
	)
	cli = NewCli(nil, nil)
	if cli == nil {
		t.Fail()
	}
	cli.TestMode = true
	cli.Flgs = []CLIFlag{
		&BoolFlg{Variable: &test, Name: "test", ShortName: "t", Usage: "use by passing -test"},
	}
	cli.Cmds = []*CLICommand{
		{
			Name:           "server",
			ShortName:      "s",
			Usage:          "use as a server",
			Value:          nil,
			Action:         func() { log.Println("cmd action") },
			PreAction:      func() { log.Println("cmd preaction") },
			PostAction:     func() { log.Println("cmd postaction") },
			BashCompletion: BashCompletionSub,
			Flags: []CLIFlag{
				&Int64Flg{Variable: &t2, Name: "port", ShortName: "p", Usage: "server port", Value: 8080, Required: true},
				&StringFlg{Variable: &protocol, Name: "protocol", ShortName: "proto", Usage: "Set Protocol http(s)", Value: "http"},
			},
		},
	}

	os.Args = []string{"cmd", "-test", "server", "-port", "8090"}
	cli.Parse()

	flg := cli.Flag("test", cli.Flgs)
	cmd := cli.Command("server")
	subcmd := cli.Flag("port", cmd.Flags)
	cases := []Tests{
		{"not required bool",
			[]Test{
				{"name", flg.GName(), "test"},
				{"shortname", flg.GShortName(), "t"},
				{"env_var", flg.GEnvVar(), "TEST"},
				{"value", test, true},
				{"required", flg.GRequired(), false},
			},
		},
		{
			"command 1",
			[]Test{
				{"name", cmd.Name, "server"},
				{"shortname", cmd.ShortName, "s"},
				{"usage", cmd.Usage, "use as a server"},
				{"value", cmd.Value, nil},
			},
		},
		{
			"command flag 1",
			[]Test{
				{"name", subcmd.GName(), "port"},
				{"shortname", subcmd.GShortName(), "p"},
				{"usage", subcmd.GUsage(), "server port"},
				{"value", t2, 8090},
			},
		},
	}

	for _, tc := range cases {
		for _, test := range tc.tests {
			t.Run(tc.name+" "+test.name, func(t *testing.T) {
				teardownSubTest := setupSubTest(t)
				defer teardownSubTest(t)

				t.Log("find", test.name, "with value", test.value)
				assert.EqualValues(t, test.value, test.field)
			})
		}
	}
}

func TestCmdHelp(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	var (
		test     bool
		t2       int64
		protocol string
	)
	cli = NewCli(nil, new(UsageDisplayStub))
	if cli == nil {
		t.Fail()
	}
	cli.TestMode = true
	cli.Flgs = []CLIFlag{
		&BoolFlg{Variable: &test, Name: "test", ShortName: "t", Usage: "use by passing -test"},
	}
	cli.Cmds = []*CLICommand{
		{
			Name:           "server",
			ShortName:      "s",
			Usage:          "use as a server",
			Value:          nil,
			Action:         func() { log.Println("cmd action") },
			PreAction:      func() { log.Println("cmd preaction") },
			PostAction:     func() { log.Println("cmd postaction") },
			BashCompletion: BashCompletionSub,
			Flags: []CLIFlag{
				&Int64Flg{Variable: &t2, Name: "port", ShortName: "p", Usage: "server port", Value: 8080, Required: true},
				&StringFlg{Variable: &protocol, Name: "protocol", ShortName: "proto", Usage: "Set Protocol http(s)", Value: "http"},
			},
		},
	}

	os.Args = []string{"cmd", "-test", "server", "-h"}
	cli.Parse()

	flg := cli.Flag("test", cli.Flgs)
	cmd := cli.Command("server")

	cases := []Tests{
		{"not required bool",
			[]Test{
				{"name", flg.GName(), "test"},
				{"shortname", flg.GShortName(), "t"},
				{"env_var", flg.GEnvVar(), "TEST"},
				{"value", test, true},
				{"required", flg.GRequired(), false},
			},
		},
		{
			"command 1",
			[]Test{
				{"name", cmd.Name, "server"},
				{"shortname", cmd.ShortName, "s"},
				{"usage", cmd.Usage, "use as a server"},
				{"value", cmd.Value, nil},
				{"displayedHelp", displayedUsage, true},
			},
		},
	}

	for _, tc := range cases {
		for _, test := range tc.tests {
			t.Run(tc.name+" "+test.name, func(t *testing.T) {
				teardownSubTest := setupSubTest(t)
				defer teardownSubTest(t)

				t.Log("find", test.name, "with value", test.value)
				assert.EqualValues(t, test.value, test.field)
			})
		}
	}
}

func TestCmdAutoComplete(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)

	var (
		test     bool
		t2       int64
		protocol string
	)
	cli = NewCli(new(FatalStub), nil)
	if cli == nil {
		t.Fail()
	}
	cli.TestMode = true
	cli.Flgs = []CLIFlag{
		&BoolFlg{Variable: &test, Name: "test", ShortName: "t", Usage: "use by passing -test"},
	}
	cli.Cmds = []*CLICommand{
		{
			Name:           "server",
			ShortName:      "s",
			Usage:          "use as a server",
			Value:          nil,
			Action:         func() { log.Println("cmd action") },
			PreAction:      func() { log.Println("cmd preaction") },
			PostAction:     func() { log.Println("cmd postaction") },
			BashCompletion: BashCompletionSub,
			Flags: []CLIFlag{
				&Int64Flg{Variable: &t2, Name: "port", ShortName: "p", Usage: "server port", Value: 8080},
				&StringFlg{Variable: &protocol, Name: "protocol", ShortName: "proto", Usage: "Set Protocol http(s)", Value: "http"},
			},
		},
	}

	os.Args = []string{"cmd", "server", "-generate-bash-completion"}
	cli.Parse()

	cmd := cli.Command("server")
	cases := []Tests{
		{
			"command 1",
			[]Test{
				{"name", cmd.Name, "server"},
				{"shortname", cmd.ShortName, "s"},
				{"usage", cmd.Usage, "use as a server"},
				{"value", cmd.Value, nil},
			},
		},
	}

	for _, tc := range cases {
		for _, test := range tc.tests {
			t.Run(tc.name+" "+test.name, func(t *testing.T) {
				teardownSubTest := setupSubTest(t)
				defer teardownSubTest(t)

				t.Log("find", test.name, "with value", test.value)
				assert.EqualValues(t, test.value, test.field)
			})
		}
	}
}

func TestInvalidObjectError_Error(t *testing.T) {
	ioe := InvalidObjectError{}
	ioe.Name = "invalid obj error"
	ioe.Error()

	assert.Equal(t, "invalid obj error", ioe.Name)

	tmp := new(CLIFlag)
	ioe.Type = reflect.TypeOf(&tmp)

	ioe.Error()
}

func TestInvalidValueError_Error(t *testing.T) {
	ive := InvalidValueError{}
	ive.Value = "someval"
	ive.Field = "field1"
	ive.Error()
	assert.Equal(t, "someval", ive.Value)
	assert.Equal(t, "field1", ive.Field)

	ive.Value = ""
	assert.Equal(t, "", ive.Value)
	ive.Error()
}

type FatalStub struct {
	adapter FatalAdapter
}

func (f *FatalStub) PrintNotice(name string) {
	printNotice = true
}
func (f *FatalStub) PrintNoticeSubCmd(name, cmd string) {
	printNoticeSubCmd = true
}

type UsageDisplayStub struct {
	adapter UsageAdapter
}

func (u *UsageDisplayStub) UsageText(cmd *CLICommand) {
	cmd.FS.Usage()
	displayedUsage = true
}
