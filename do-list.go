package proptag

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//─────────────┤ doList ├─────────────

func doList(in []byte, file ...string) error {
	var ln, blkend, strend int
	slice := strings.Split(string(in), "\n")
	f := "dev/stdout"
	if len(file) > 0 {
		_, f = filepath.Split(file[0])
		fmt.Fprintln(os.Stdout, "listing file", f)
		ln = 1 + len("listing file") + len(f)
	}
	if ln == 0 {
		ln = 24
	}
	fmt.Fprintln(os.Stdout, strings.Repeat("-", ln))
	linepat, _ := regexp.Compile(LinePat)
	openpat, _ := regexp.Compile(OpenPat)
	closepat, _ := regexp.Compile(ClosePat)
	varblkpat, _ := regexp.Compile(VarBlkPat)
	structpat, _ := regexp.Compile(StructPat)

	inblk := false
	count := 0
	for i, line := range slice {
		line = strings.TrimLeft(line, "\t ")
		//	fmt.Println("line strend blkend", line, strend, blkend)
		if len(line) == 0 {
			continue
		}
		if blkend > 0 && i >= blkend {
			blkend = 0
			inblk = false
		}
		if strend > 0 && i >= strend {
			strend = 0
			inblk = false
		}
		if inblk {
			fmt.Fprintln(os.Stdout, "tagged line #", i, ": ", line)
			count++
		}

		if openpat.MatchString(line) {
			inblk = true
			continue
		}
		if closepat.MatchString(line) {
			inblk = false
			continue
		}
		if linepat.MatchString(line) {
			fmt.Fprintln(os.Stdout, "tagged line #", i, ": ", line)
			count++
		}
		if varblkpat.MatchString(line) {
			inblk = true
			prn := 0
			blktext := slice[i:]
			blkend = i
			for _, l := range blktext {
				for _, r := range l {
					if r == '(' {
						prn++
					}
					if r == ')' {
						prn--
					}
				}
				blkend++
				if prn == 0 {
					break
				}
				if prn < 0 {
					return fmt.Errorf("file %s line %d mismatched parentheses", filepath.Base(f), blkend)
				}
			}
		}
		if structpat.MatchString(line) {
			inblk = true
			brc := 0
			blktext := slice[i:]
			strend = i
			for _, l := range blktext {
				for _, r := range l {
					if r == '{' {
						brc++
					}
					if r == '}' {
						brc--
					}
				}
				strend++
				if brc == 0 {
					break
				}
				if brc < 0 {
					return fmt.Errorf("file %s line %d mismatched braces", filepath.Base(f), strend)
				}
			}
		}
	}

	if count == 0 {
		fmt.Fprintln(os.Stdout, "found no tagged lines")
	}
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, strings.Repeat("-", ln))
	fmt.Fprintln(os.Stdout, "")
	return nil
}
