// main.go v2
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/egp/goCF/cf"
)

func main() {
	p := flag.Int64("p", 355, "numerator")
	q := flag.Int64("q", 113, "denominator (non-zero)")
	n := flag.Int("n", 12, "number of CF terms to print")
	flag.Parse()

	if *q == 0 {
		fmt.Fprintln(os.Stderr, "error: q must be non-zero")
		os.Exit(1)
	}
	if *n < 0 {
		fmt.Fprintln(os.Stderr, "error: n must be non-negative")
		os.Exit(1)
	}

	r, err := cf.NewRational(*p, *q)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid rational %d/%d: %v\n", *p, *q, err)
		os.Exit(1)
	}

	stream := cf.NewRationalCF(r)

	fmt.Printf("r = %s\n", r.String())
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

// main.go v2
