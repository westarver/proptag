// MIT License

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//--------------------
//Author: James Tarver
//using tags:
//prop-tag will read go source files and generate getters and setters
//for variables that are tagged with a tag of the proper format.
//These tags can be set on package level variables or inside structs
//on fields.  If placed on fields the functions will be generated as
//methods on the containing struct, otherwise as stand-alone exported
//functions.  Tags can be, and usually are, placed at the end of the
//line to be operated on.  In this case the one-liner format is
//required.  This tag consists of the opening bracket '<' then the
//tag string starting with 'p'.  If the property is to be assigned
//and its type cannot or will not be assigned using the assignment
//operator, then uppercase S can be used in the tag string to
//cause an assignment via a user implemented Set function or method.
//Variables that require special treatment, such as cloning or deep
//copy in the getter can be signified using an uppercase 'G' as in
//<pG/>.  In this case a pointer to the new or copied variable will
//be returned by default.The tag must be closed using  />.  This
//will result in a tag like so; <pg|G/>, <ps|S/>, <pg|Gs|S/>.
//Blocks of properties can be designated by placing an opening tag
//"<tag>" on a line by itself, a list of variable declarations,
//then a closing tag "</prop> on a line by itself. Variables in a var
//block will work as well.  Finally, an asterisk after the g|G or s|S
//will cause the getting or setting to be done via pointer.
//If both the g and s are used in a single tag, then the g must be first
//
//To state the obvious, these tags must be shielded from the go tooling behind comments.
package main

import (
	"os"

	msg "github.com/westarver/messenger"
	"github.com/westarver/proptag"
)

func main() {
	writer := msg.New()
	writer.SetLogoutStr()
	exitCode := proptag.Run(writer)
	os.Exit(exitCode)
}
