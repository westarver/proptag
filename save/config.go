package proptag

import (
	"github.com/westarver/boa"
)

const (
	actionHelp ActionType = iota
	actionList
	actionGenerate
	actionDelete
	actionError
)

//────────────────────┤ prepArgs ├────────────────────

func prepArgs() (*boa.CLI, ActionType) {
	cli := boa.FromHelp(getUsage())
	var action ActionType

	_, d := cli.Items["delete"].(boa.CmdLineItem[bool])
	if d {
		action = actionDelete
	}

	_, l := cli.Items["list"].(boa.CmdLineItem[bool])
	if l {
		action = actionList
	}

	_, g := cli.Items["generate"].(boa.CmdLineItem[bool])
	if g {
		action = actionGenerate
	}

	_, h := cli.Items["help"].(boa.CmdLineItem[string])
	if h { //help is unique in that the command causes an early exit
		return cli, actionHelp
	}
	if !l && !d && !g { // default command is help
		return cli, actionHelp
	}

	return cli, action
}
