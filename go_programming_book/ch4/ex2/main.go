package main

import (
	"fmt"
)

func main() {
	var x []int
	fmt.Printf("%v\n", cap(x))
	x = append(x, 1)
	fmt.Printf("%v\n", cap(x))
	x = append(x, 2)
	fmt.Printf("%v\n", cap(x))
	x = append(x, 3)
	fmt.Printf("%v\n", cap(x))
	x = append(x, 4)
	fmt.Printf("%v\n", cap(x))
	x = append(x, 5)
	fmt.Printf("%v\n", cap(x))
	x = append(x, 6)
	fmt.Printf("%v\n", cap(x))
}
