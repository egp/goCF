// main.go v1
package main

import (
	"flag"
	"fmt"

	"github.com/egp/goCF/cf"
)

func main() {
	p := flag.Int64("p", 355, "numerator")
	q := flag.Int64("q", 113, "denominator (non-zero)")
	n := flag.Int("n", 12, "number of CF terms to print")
	flag.Parse()

	if *q == 0 {
		panic("q must be non-zero")
	}

	r := cf.Rational{P: *p, Q: *q}
	stream := cf.NewRationalCF(r)

	fmt.Printf("r = %d/%d\n", r.P, r.Q)
	fmt.Printf("CF terms (first %d): [", *n)

	for i := 0; i < *n; i++ {
		a, ok := stream.Next()
		if !ok {
			break
		}
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(a)
	}
	fmt.Println("]")
}

// main.go v1
