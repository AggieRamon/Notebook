package main

import (
	"fmt"
)

/*Modify reverse to reverse the characters of a []byte slice that represents a UTF-8 encoded string, in place. can you do it without allocating new memory*/

func main() {
	s := "Hello World"
	x := reverse([]byte(s))
	fmt.Printf("%s\n", x)
}

func reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

/*Write an in-place function that squashes each run of adjacent Unicode spaces (see unicode.IsSpace) in a UTF-8-encoded []byte slice into a single ASCII space*/

/*func main() {
	myString := "Hello   World   What"
	x := squashSpace([]byte(myString))
	fmt.Printf("%s\n", x)
}

func squashSpace(bytes []byte) []byte {
	var out []byte
	var last rune
	for i := 0; i < len(bytes); {
		r, s := utf8.DecodeRune(bytes[i:])
		if !unicode.IsSpace(r) {
			out = append(out, bytes[i:i+s]...)
		} else if unicode.IsSpace(r) && !unicode.IsSpace(last) {
			out = append(out, ' ')
		}
		last = r
		i += s
	}
	return out
}*/

/*Exercise 4.5: Write an in-place function to elimate adjacent duplicates in []string slice*/

/*func main() {
	arr := []string{"abc", "abc", "bcc", "bcc", "cdc"}
	arr = removeAdjDups(arr)
	fmt.Println(cap(arr), len(arr), arr)
}

func removeAdjDups(arr []string) []string {
	for i := 0; i < len(arr)-1; {
		if arr[i] == arr[i+1] {
			arr = append(arr[:i], arr[i+1:]...)
		} else {
			i++
		}
	}
	return arr
}*/

/* Exercise 4.3: Rewrite reverse to use an array pointer instead of a slice */
/*func main() {
	s := [6]int{1, 2, 3, 4, 5, 6}
	reverse(&s)
	fmt.Printf("%v\n", s)
}

func reverse(s *[6]int) {
	fmt.Printf("%v\n%T\n", s[2], s)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}*/
