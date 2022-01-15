package main

import (
	"fmt"
	"os"
	"strconv"

	"aggieramon.com/conv/lib/conv"
)

func main() {
	for _, arg := range os.Args[1:] {
		t, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		in := conv.Inch(t)
		cm := conv.Centimeter(t)

		fmt.Printf("%s = %s, %s = %s\n", in, conv.InToCm(in), cm, conv.CmToIn(cm))
	}
}
