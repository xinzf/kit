package main

import (
	diff "github.com/r3labs/diff/v3"
	"log"
)

type Order struct {
	ID    string `diff:"id"`
	Items []int  `diff:"items"`
}

func main() {
	a := Order{
		ID:    "1234",
		Items: []int{1, 2, 4},
	}

	b := Order{
		ID:    "1234",
		Items: []int{1, 2, 4},
	}

	changelog, err := diff.Diff(a, b, diff.DisableStructValues(), diff.AllowTypeMismatch(true))
	//changelog, err := diff.Diff(a, b)
	if err != nil {
		panic(err)
	}
	log.Println(changelog)
}
