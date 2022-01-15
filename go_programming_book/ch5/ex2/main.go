package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/html"
)

func main() {
	n, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	count := make(map[string]int)
	count, _ = findCount(count, n)
	fmt.Printf("%+v\n", count)
}

func findCount(count map[string]int, n *html.Node) (map[string]int, error) {
	if n.Type == html.ElementNode {
		count[n.Data]++
		// switch n.Data {
		// case "div":
		// 	count["div"]++
		// case "li":
		// 	count["li"]++
		// case "span":
		// 	count["span"]++
		// case "p":
		// 	count["p"]++
		// case "a":
		// 	count["a"]++
		// default:
		// 	count["other"]++
		// }
	}
	if n.FirstChild != nil {
		count, _ = findCount(count, n.FirstChild)
	}

	if n.NextSibling != nil {
		count, _ = findCount(count, n.NextSibling)
	}
	return count, nil
}
