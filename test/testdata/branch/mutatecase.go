//go:build test
// +build test

package main

import (
	"fmt" // user comment
)

func main() {
	i := 1

	for i != 4 {
		switch {
		case i == 1:
			fmt.Println(i)
		case i == 2:
			fmt.Println(i * 2)
		default:
			fmt.Println(i * 3)
		}

		i++
	}
}
