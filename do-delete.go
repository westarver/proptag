package proptag

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

//─────────────┤ doDelete ├─────────────

func doDelete(inp []byte, out io.Writer) error {
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
