package proptag

import (
	"os"

	cfg "github.com/westarver/configure"

	"github.com/docopt/docopt-go"
)

const (
	actionHelp ActionType = iota
	actionList
	actionGenerate
	actionUncomment
	actionRemove
	actionMock
	actionError
)

//────────────────────┤ prepArgs ├────────────────────

func prepArgs() ([]string, []string, ActionType, bool) {

	var ins, outs, fi, fo, ci, co []string

	if len(os.Args) < 2 {
		return ins, outs, actionGenerate, false
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		showHelp()
		return ins, outs, actionHelp, true
	}

	var parser = &docopt.Parser{
		HelpHandler:   ShowHelp,
		OptionsFirst:  false,
		SkipHelpFlags: false,
	}

	opts, _ := parser.ParseArgs(getUsage(), os.Args[1:], "")
	// merge ins and outs from file and command line
	file, _ := opts.String("<file>")
	if len(file) != 0 {
		fi, fo = getIOFromFile(file)
	}
	ci, _ = opts["<inputs>"].([]string)
	co, _ = opts["<outputs>"].([]string)
	//the following is compensation for docopts inability to
	//recognize the -- as a separator
	if len(co) > 0 {
		ci, co = splitOnDashDash(ci, co)
	}

	ins = merge(fi, ci)
	outs = merge(fo, co)

	var action ActionType
	var l, g bool
	l, _ = opts.Bool("-l")
	if !l {
		l, _ = opts.Bool("--list")
	}
	g = !l

	if l {
		action = actionList
	}
	if g {
		action = actionGenerate
	}

	return ins, outs, action, false
}

//─────────────┤ merge ├─────────────

func merge(s1, s2 []string) []string {
	var ret []string
	var found bool

	s1 = append(s1, s2...)
	for _, s := range s1 {
		found = false
		for _, t := range ret {
			found = (s == t)
		}
		if !found {
			ret = append(ret, s)
		}
	}
	return ret
}

//─────────────┤ getIOFromFile ├─────────────

func getIOFromFile(file string) ([]string, []string) {
	var i, o []string
	inif, err := cfg.ReadIni(file)
	if err == nil {
		i = inif.Section("inputs").KeyStrings()
		o = inif.Section("outputs").KeyStrings()
	}
	return i, o
}

//─────────────┤ splitOnDashDash ├─────────────

func splitOnDashDash(inputs, outputs []string) ([]string, []string) {
	for i, s := range outputs {
		if s == "--" {
			inputs = append(inputs, outputs[i+1:]...)
			return inputs, outputs[:i]
		}
	}

	return inputs, outputs
}
