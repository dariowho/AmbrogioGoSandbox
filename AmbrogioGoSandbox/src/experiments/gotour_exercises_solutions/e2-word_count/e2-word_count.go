package main

import (
	"strings"
	"code.google.com/p/go-tour/wc"
)

func WordCount(s string) map[string]int {
	r := map[string]int{}
	
	for _,w := range strings.Fields(s) {
		r[w] = r[w] + 1
	}
	
	return r
}

func main() {
	wc.Test(WordCount)
}
