package proptag

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	trace "github.com/westarver/tracer"
)

//─────────────┤ doDelete ├─────────────

func doDelete(inp []byte, out io.Writer) error {
	var trace = trace.New(os.Stderr)                                    //<rmv/>
	trace.Trace("----------------------------entering doDelete\n")      //<rmv/>
	defer trace.Trace("----------------------------leaving doDelete\n") //<rmv/>
	trace.Trace("out ", out)                                            //<rmv/>
	linepat, _ := regexp.Compile(LinePat)
	openpat, _ := regexp.Compile(OpenPat)
	closepat, _ := regexp.Compile(ClosePat)
	slice := strings.Split(string(inp), "\n")
	var sb strings.Builder
	for _, line := range slice {
		if openpat.MatchString(line) {
			continue
		}
		if closepat.MatchString(line) {
			continue
		}
		if linepat.MatchString(line) {
			line = linepat.ReplaceAllString(line, "")
			line = strings.TrimSuffix(line, "//")
		}
		sb.WriteString(line + "\n")
	}
	_, err := fmt.Fprintln(out, sb.String())
	return err
}
