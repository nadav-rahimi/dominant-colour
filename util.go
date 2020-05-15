package dominant_colour

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"log"
)

// Error checking function to reduce boilerplate code.
// If boolean, false triggers an error, if error, it not
// being equal to nil triggers an error
func chkFatal (e interface{}) {
	switch v := e.(type) {
	case error:
		if v != nil {
			log.Fatal(v)
		}
	case bool:
		if v == false {
			log.Fatal(v)
		}
	}
}

// Prints out the mat.Matrix object, used for debugging
func matPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n", fa)
}
