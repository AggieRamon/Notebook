package main

import (
	"fmt"
	"os"
	"time"
	//"strings"
)

func main() {
	start := time.Now()
	fmt.Println(os.Args[0])
	//fmt.Println(strings.Join(os.Args[1:], " "))
	for _, arg := range os.Args[1:] {
		fmt.Println(arg)
	}
	fmt.Printf("%fs elapsed \n", time.Since(start).Seconds())
}
