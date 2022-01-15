package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"os"
)

func main() {
	cryptType := flag.String("crypt", "sha256", "Allowed to be sha256, sha384, or sha512")
	flag.Parse()
	info, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Println(*cryptType)
	fmt.Println(info)
	if *cryptType == "sha256" {
		s256 := sha256.Sum256([]byte(info))
		fmt.Printf("%x\n", s256)
	} else if *cryptType == "sha384" {
		s384 := sha512.Sum384([]byte(info))
		fmt.Printf("%x\n", s384)
	} else {
		s512 := sha512.Sum512([]byte(info))
		fmt.Printf("%x\n", s512)
	}
}
