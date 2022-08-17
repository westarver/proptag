package proptag

import (
	"fmt"
	"os"
)

var help = `Usage: prop-tag [-l | --list] [-f <file> |--file <file> ] [-o <outputs>... | --out <outputs>...] [<inputs>...]
prop-tag [-h | --help]

Options:
-h --help                  	show help and exit. 
-f <file> --file <file>	   	file with list of inputs and outputs
-l --list			list the lines that are tagged
-o <outputs>... --out <outputs>...  filenames to write the generated methods to

inputs are optional.   	empty list will enable input from stdin.
wildcard, filename(s) or blank are all valid

outputs are optional. 	empty list will write output to stdout. use
-- to separate output list from input list.

examples: 
prop-tag [enter] will use stdin and stdout
prop-tag  -o a b -- c d  will use a and b for outputs and c and d for input.  
prop-tag -o a b will cause a and b to be interpreted as save-to file and 
stdin as input.
prop-tag -l will generate a listing of tagged lines to stdout.	

using tags:		
prop-tag will read go source files and generate getters and setters
for variables that are tagged with a tag of the proper format.  
These tags can be set on package level variables or inside structs
on fields.  If placed on fields the functions will be generated as
methods on the containing struct, otherwise as stand-alone exported
functions.  Tags can be, and usually are, placed at the end of the
line to be operated on.  In this case the one-liner format is
required.  This tag consists of the opening bracket '<' then the 
tag string starting with 'p'.  If the property is to be assigned 
and its type cannot or will not be assigned using the assignment
operator, then uppercase S can be used in the tag string to
cause an assignment via a user implemented Set method.
Variables that require special treatment, such as cloning or deep
copy in the getter can be signified using an uppercase 'G' as in 
<pG/>.  In this case the variable will be returned using a Get Method
on the variable or field. The tag must be closed using  />.  This 
will result in a tag like so; <pg|G/>, <ps|S/>, <pg|Gs|S/>.  
Blocks of properties can be designated by placing an opening tag
"<tag>" on a line by itself, a list of variable declarations, 
then a closing tag "</prop> on a line by itself. Variables in a var
block will work as well.  Finally, an asterisk after the g|G or s|S
will cause the getting or setting to be done via pointer.
If both the g and s are used in a single tag, then the g must be first

To state the obvious, these tags must be shielded from the go tooling behind comments. 

In its simplest form:		
  // generate a getter and setter
  var some_var some_type  //<pgs/>
  func Some_var() some_type {
  	return some_var
  }
  func Setsome_var(v some_type) {
  	some_var = v
  }
    
 Intermediate:
  var some_var composite_type  //<pGS/>
  // generated code
  func Some_var() *composite_type {
  	return some_var.Get()
  }
  func Setsome_var(v composite_type)  {
  	some_var.Set(v)
  }
  
  Or the most complicated form:
  type Astruct struct {
  	//<pGS>
  	var1 composite_type1
  	var2 composite_type2
  	//</prop>
  	...
  }
  func(s *Astruct) Setvar1(v composite_type1)  {
  	s.var1.Set(v)
  }
  func(s *Astruct) Var1() *composite_type1  {
 	return s.var1.Get()
  }
  // same for var2
`

//───────────────┤ getUsage ├───────────────

func getUsage() string {
	return help
}

//───────────────┤ ShowHelp ├───────────────

func ShowHelp(err error, help string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", getUsage())
		os.Exit(1)
	} else {
		fmt.Println(getUsage())
	}
}

//───────────────┤ showHelp ├───────────────

func showHelp() {
	fmt.Println(getUsage())
}
