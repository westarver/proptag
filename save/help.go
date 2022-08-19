package proptag

import (
	"fmt"
	"io"
)

//───────────────┤ getUsage ├───────────────

func getUsage() string {
	var help = `Usage: proptag [command] [-o <outputs>... | --out <outputs>...] [-s <inputs>... | --source <inputs>...]
	
Commands:
*+[help] <topic>   :show help and exit. 
*[generate | gen]  :generate the functions and methods indicated
*[delete | del]	   :delete the tags in the input source files
*[list]			   :list the lines that are tagged

Flags:
[--out | -o] <outputs>...    :filenames to write the generated functions to
[--source | -s] <inputs>...  :filenames to read tags from
	
Long Description:
--source:   Sources are optional.  Missing flag will enable obtaining source file names from stdin.
			Wildcard, filename(s) or blank are all valid.
			Filenames with spaces require quotes.
			Terminate the space separated list with -- if other options follow.

--out:      Outputs are optional.  Missing flag will write output the terminal.
			Terminate the space separated list with -- if other options follow.
 
More:
	examples: 
	proptag <command> will use stdin and stdout
	proptag <command> -o a b -- -s c d  will use a and b for outputs and c and d for input.  
	proptag <command> -o a b will cause a and b to be interpreted as output file and 
	stdin as input.
	proptag [enter] will display this help text and exit
	
	using tags:		
	proptag will read go source files and generate getters and setters
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
	After generating the methods/functions the tags themselves can be deleted 
	by running again using the del or delete command
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
	return help
}

//───────────────┤ ShowHelp ├───────────────

func ShowHelp(w io.Writer, command ...string) {
	var help = `Usage: proptag [command] [-o <outputs>... | --out <outputs>...] [-s <inputs>... | --source <inputs>...]
	
Commands:
[help] <topic>    :show help and exit. 
[generate | gen]  :generate the functions and methods indicated
[delete | del]	  :delete the tags in the input source files
[list]            :list the lines that are tagged to the terminal

Flags:
-o <outputs>... --out <outputs>...   :-filenames to write the generated functions to
-s <inputs>... --source <inputs>...  :-filenames to read tags from
	
Long Description:
--source:   Sources are optional.  Missing flag will enable obtaining source file names from stdin.
			Wildcard, filename(s) or blank are all valid.
			Filenames with spaces require quotes.
			Terminate the space separated list with -- if other options follow.

--out:      Outputs are optional.  Missing flag will write output the terminal.
			Terminate the space separated list with -- if other options follow.
 
More:
	examples: 
	proptag <command> will use stdin and stdout
	proptag <command> -o a b -- -s c d  will use a and b for outputs and c and d for input.  
	proptag <command> -o a b will cause a and b to be interpreted as output file and 
	stdin as input.
	proptag [enter] will display this help text and exit
	
	using tags:		
	proptag will read go source files and generate getters and setters
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
	After generating the methods/functions the tags themselves can be deleted 
	by running again using the del or delete command
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
	if len(command) != 0 {
		fmt.Fprintf(w, "%s\n", command[0])
	} else {
		fmt.Fprintln(w, help)
	}
}
