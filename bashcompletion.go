package mycli

import (
	"flag"
	"fmt"
	"strings"
)

func BashCompletionMain(c *CLI) {

	if flag.NArg() > 0 {
		return
	}
	// get all flags and print them
	for _, d := range c.Flgs {
		low := strings.ToLower(d.GName())
		if !d.GHidden() {
			fmt.Fprintln(c.Writer, strings.ToLower("-"+low))
		} else if low == "version" {
			fmt.Fprintln(c.Writer, "-v,-version")
		}
	}
	for _, d := range c.Cmds {
		low := strings.ToLower(d.Name)
		if !d.Hidden {
			fmt.Fprintln(c.Writer, strings.ToLower(low))
		} else if low == "version" {
			fmt.Fprintln(c.Writer, "v,version")
		}
	}
}

func BashCompletionSub(c *CLI, cm *CLICommand) {

	if cm.FS.NArg() > 0 {
		return
	}
	for _, d := range cm.SubCommands {
		low := strings.ToLower(d.Name)
		if !d.Hidden {
			fmt.Fprintln(c.Writer, strings.ToLower(low))
		} else if low == "version" {
			fmt.Fprintln(c.Writer, "v,version")
		}
	}

	for _, d := range cm.Flags {
		low := strings.ToLower(d.GName())
		if !d.GHidden() && low != "help" {
			fmt.Fprintln(c.Writer, "-"+low)
		} else if low == "help" {
			fmt.Fprintln(c.Writer, "-h,-help")
		}
	}
}
