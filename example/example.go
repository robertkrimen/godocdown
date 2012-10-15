/*
Package example is an example package with documentation 

	// Here is some code
	func example() {
		abc := 1 + 1
	}()

Installation

	# This is how to install it:
	$ curl http://example.com
	$ tar xf example.tar.gz -C .
	$ ./example &

*/
package example

// A constant section
const Other = 3

// A variable section
var (
	This = 1

	this = 0


	// A description of That
	That = 2.1
)

// Another constant section
const (
	Another = 0
	Again = "this"
)

// Example is a function that does nothing
func Example() {
}

// ExampleType is a type of nothing
//		
//		// Here is how to use it:
//		return &ExampleType{
//			First: 1,
//			Second: "second",
//			nil,
//		}
type ExampleType struct {
	First int
	Second string
	Third float64
	Parent *ExampleType

	first int
	hidden string
}

func (ExampleType) Set() bool {
	return false
}

func NewExample() *ExampleType {
	return &ExampleType{}
}
