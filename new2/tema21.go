package main

import (
	"strings"

	"github.com/chrislusf/glow/flow"
)

func main() {
    flow.New().TextFile(
        "passwd.txt", 3,
    ).Map(func(line string, ch chan string) {
        for _, token := range strings.Split(line, ":") {
            ch <- token
        }
    }).Map(func(key string) int {
        vowels := 0
        others := 0
        for _, char := range key {
            if char == 'a' || char == 'e' || char == 'i' || char == 'o' || char == 'u' {
                vowels++
            } else {
                others++
            }
        }
        if ((vowels % 2) == 0) && others % 3 == 0 {
            return 1
        }
        return 0
    }).Reduce(func(x int, y int) int {
        return x + y
    }).Map(func(x int) {
        println("count:", float64(float64(x)/float64(3)))
    }).Run()
}

