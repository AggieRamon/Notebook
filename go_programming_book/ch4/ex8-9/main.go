package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ex 4.9: Write a program wordfreq to report the frequency of each word in an input text file. Call input.Split(bufio.ScanWords) before the first call to Scan to break the input into words instead

func main() {
	count := make(map[string]int)

	f, err := os.Open("./test.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	s.Split(bufio.ScanWords)

	for s.Scan() {
		count[strings.ToLower(s.Text())]++
	}

	fmt.Println(count)
}

// func wordFreq()

// Ex 4.8: Modify charcount to count letters, digits, and so on in their unicode categories using functions like unicode.IsLetter

/*func main() {
	myString := "Hello World **&^%$ 1234jkdj5n6j7fjde8e"
	count := charcount(myString)
	fmt.Println(count)
}

func charcount(s string) map[string]int {
	count := make(map[string]int)
	for _, v := range s {
		if unicode.IsLetter(v) {
			count["letter"]++
		} else if unicode.IsDigit(v) {
			count["number"]++
		} else if unicode.IsSpace(v) {
			count["space"]++
		} else if unicode.IsSymbol(v) {
			count["symbol"]++
		}
	}

	return count
}*/
