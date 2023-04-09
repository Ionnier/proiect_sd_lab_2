package main

import (
	"strings"

	"github.com/chrislusf/glow/flow"
)

func main() {
    flow.New().TextFile(
        "passwd2.txt", 3,
    ).Map(func(line string, ch chan string) {
        for _, token := range strings.Split(line, ":") {
            ch <- token
        }
    }).Map(func(key string) int {
        if Anagram(key, "facultate") {
			return 1
		}
        return 0
    }).Reduce(func(x int, y int) int {
        return x + y
    }).Map(func(x int) {
        println("count:", float64(float64(x)/float64(3)))
    }).Run()
}


func Anagram(s string, t string) bool {
	string1 := len(s)
	string2 := len(t)
	if string1 != string2 {
	   return false
	}

	anagramMap := make(map[string]int)
	
	for i := 0; i < string1; i++ {
	   anagramMap[string(s[i])]++
	}
	
	for i := 0; i < string2; i++ {
	   anagramMap[string(t[i])]--
	}
	
	for i := 0; i < string1; i++ {
	   if anagramMap[string(s[i])] != 0 {
		  return false
	   }
	}
	
	return true
 }
