package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Println(isAnagram("becca", "eacbc"))

	test := "Howareyoutoday"
	test_update := comma(test)
	fmt.Println(test_update)

	values := []float64{1.2, 3.4, 5.6}
	fstring := floatToString(values)
	fmt.Println(fstring)
}

func comma(s string) string {
	var buf bytes.Buffer
	for i, val := range s {
		if i%3 == 0 && i != 0 {
			buf.WriteString(",")
		}
		buf.WriteRune(val)
	}
	return buf.String()
}

func floatToString(values []float64) string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, val := range values {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%f", val)
	}
	buf.WriteByte(']')
	return buf.String()
}

func isAnagram(s1 string, s2 string) bool {
	if len(s1) == len(s2) {
		s1Map := make(map[rune]int)
		for _, val := range s1 {
			_, ok := s1Map[val]
			if ok {
				s1Map[val] = s1Map[val] + 1
			} else {
				s1Map[val] = 1
			}
		}

		s2Map := make(map[rune]int)
		for _, val := range s2 {
			_, ok := s2Map[val]
			if ok {
				s2Map[val] = s2Map[val] + 1
			} else {
				s2Map[val] = 1
			}
		}

		if len(s1Map) == len(s2Map) {
			for key := range s1Map {
				_, ok := s2Map[key]
				if !ok || s1Map[key] != s2Map[key] {
					return false
				}
			}
			return true
		} else {
			return false
		}
	}
	return false
}
