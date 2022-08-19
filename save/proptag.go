package proptag

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/westarver/boa"
	"github.com/westarver/helper"
	msg "github.com/westarver/messenger"
	trace "github.com/westarver/tracer"
)

type ActionType int

const (
	CommentString = "//"
)

// regex patterns to locate potential variables to process
const (
	LPar      = `\050`
	RPar      = `\051`
	LBrk      = `\133`
	RBrk      = `\135`
	LBrc      = `\123`
	RBrc      = `\125`
	Aster     = `\052`
	WS0       = `[\t ]*`
	WS1       = `[\t ]+`
	IdentList = IdentPat + WS0 + `(,` + WS1 + IdentPat + `)*`
	TParam    = `(` + LBrk + `(_?[A-Za-z_]+[0-9]*[ ,]*)+` + RBrk + `)?`
	IdentPat  = `_?[A-Za-z_]+[0-9]*`
	TypePat   = Aster + `?` + IdentPat + WS0 + TParam
	LinePat   = `<p((g|G)?` + Aster + `?(s|S)?)?` + Aster + `?/>`
	OpenPat   = `<p((g|G)?` + Aster + `?(s|S)?` + Aster + `?)?>`
	ClosePat  = `</prop>`
	VarPat    = WS0 + `var` + WS1 + IdentPat + WS1 + TypePat
	VarBlkPat = WS0 + `var` + WS0 + LPar
	StructPat = WS0 + `type` + WS1 + IdentPat + WS0 + TParam + WS1 + `struct` + WS0 + `{`
	InBlkPat  = WS0 + `^(var)` + IdentPat + WS1 + TypePat
)

//─────────────┤ Run ├─────────────

func Run(m *msg.Messenger) int {
	var trace = trace.New(os.Stderr)                               //<rmv/>
	trace.Trace("----------------------------entering Run\n")      //<rmv/>
	defer trace.Trace("----------------------------leaving Run\n") //<rmv/>

	cli, action := prepArgs()
	var err error
	var ins, outs []string

	if action == actionHelp { //help is unique in that the command causes an early exit
		hlp, b := cli.Items["help"].(boa.CmdLineItem[string])
		if b {
			topic := hlp.Value()
			item, exist := cli.AllHelp[topic]
			if exist {
				ShowHelp(m, item)
			} else {
				ShowHelp(m)
			}
			return int(action)
		}
		ShowHelp(m)
		return int(actionHelp)
	}

	src, b1 := cli.Items["--source"].(boa.CmdLineItem[[]string])
	if b1 {
		ins = src.Value()
		if len(ins) == 0 {
			ins = append(ins, "-")
		}
	} else {
		ins = append(ins, "-")
	}
	trace.Trace("ins ", ins) //<rmv/>
	o, b2 := cli.Items["--out"].(boa.CmdLineItem[[]string])
	if b2 {
		outs = o.Value()
		if len(outs) == 0 {
			ins = append(outs, "-")
		}
	} else {
		outs = append(outs, "-")
	}
	trace.Trace("outs ", outs) //<rmv/>

	err = performCommand(action, ins, outs)
	trace.Trace("error returned from performCommand ", err) //<rmv/>
	if m.Catch(msg.LOG, err) != nil {
		action = actionError
	}

	m.InfoMsg(m.Logout(), msg.MESSAGE, "Exiting with exit code %d", int(action))
	return int(action)
}

//─────────────┤ performCommand ├─────────────

func performCommand(act ActionType, ins, outs []string) error {
	var trace = trace.New(os.Stderr)                                          //<rmv/>
	trace.Trace("----------------------------entering performCommand\n")      //<rmv/>
	defer trace.Trace("----------------------------leaving performCommand\n") //<rmv/>

	matched := helper.Matchio(ins, outs, ".prop")
	trace.Trace("matched ", matched) //<rmv/>
	var err error
	for _, m := range matched {
		err = process(act, m.In, m.Out)
	}
	return err
}

//─────────────┤ process ├─────────────

func process(act ActionType, in string, o string) error {
	var trace = trace.New(os.Stderr)                                   //<rmv/>
	trace.Trace("----------------------------entering process\n")      //<rmv/>
	defer trace.Trace("----------------------------leaving process\n") //<rmv/>
	var stdin, stdout, piping bool
	var inp []byte
	var err error

	trace.Trace("in ", in) //<rmv/>
	if in == "-" {
		stdin = true
	}
	trace.Trace("out ", o) //<rmv/>
	if o == "-" {
		stdout = true
	}
	if stdin {
		// user enters file name at command line
		// or the file is read through stdin piping
		o, _ := os.Stdin.Stat()
		if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice { //Terminal
			piping = false
		} else {
			piping = true
		}
		trace.Trace("piping ", piping) //<rmv/>
		if piping {
			var data bytes.Buffer
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				fmt.Fprintln(&data, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				err = errors.New("reading standard input pipe: " + err.Error())
				return err
			}
			inp = data.Bytes()
		} else { // get a file name
			fmt.Print("Enter a file path > ")
			inb := make([]byte, 1024)
			n, _ := os.Stdin.Read(inb)
			if n > 1 {
				in = string(inb)
				in = strings.Trim(in, "\000\t \n")
			} else {
				return fmt.Errorf("no filename entered via terminal")
			}
		}
	}
	if !piping { // we have a filename
		in, err = helper.ValidatePath(in)
		if err != nil {
			return err
		}
		f, err := os.Open(in)
		trace.Trace("error ", err) //<rmv/>
		if err != nil {
			return err
		}
		st, err := f.Stat()
		if err != nil {
			return err
		}
		sz := st.Size()
		inp = make([]byte, sz+256)
		_, err = f.Read(inp)
		if err != nil {
			return err
		}
		f.Close()
	}
	if in == "-" {
		in = "dev/stdin"
	}

	if act == actionList {
		return doList(inp, in)
	}

	var prop *os.File
	if !stdout {
		o, err = helper.ValidatePath(o)
		trace.Trace("error ", err) //<rmv/>
		if err != nil {
			return err
		}

		prop, err = os.OpenFile(o, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		defer prop.Close()
	} else {
		prop = os.Stdout
	}

	if act == actionDelete {
		return doDelete(inp, prop)
	}

	errs := doGenerate(inp, in, prop)
	var errstr []string
	if len(errs) != 0 {
		for _, e := range errs {
			errstr = append(errstr, e.Error())
		}
		return errors.New(strings.Join(errstr, "\n"))
	}
	return nil
}
