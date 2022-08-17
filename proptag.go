package proptag

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/westarver/helper"
	msg "github.com/westarver/messenger"
)

type ActionType int

const (
	CommentString = "//"
)

//─────────────┤ Run ├─────────────

func Run(m *msg.Messenger) int {
	ins, outs, action, exit := prepArgs()
	var err error
	if !exit {
		err = performAction(action, ins, outs)
	}
	if err != nil {
		m.Catch(msg.LOG, err)
	}
	m.InfoMsg(os.Stderr, msg.MESSAGE, "Exiting with exit code %d", int(action))
	return int(action)
}

//─────────────┤ performAction ├─────────────

func performAction(act ActionType, ins, outs []string) error {
	matched := helper.Matchio(ins, outs, ".prop")
	var err error
	for _, m := range matched {
		err = process(act, m.In, m.Out)
	}
	return err
}

//─────────────┤ process ├─────────────

func process(act ActionType, in string, o string) error {
	var stdin, stdout, piping bool
	var inp []byte
	var err error

	if in == "-" {
		stdin = true
	}
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
			inb := make([]byte, 256)
			n, _ := os.Stdin.Read(inb)
			if n > 1 {
				in = string(inb)
				in = strings.Trim(in, "\000\t \n")
			} else {
				return fmt.Errorf("no filename entered at command line")
			}
		}
	}
	if !piping { // we have a filename
		in, err = helper.ValidatePath(in)
		if err != nil {
			return err
		}
		f, err := os.Open(in)
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
