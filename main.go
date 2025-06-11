package main

import (
	"fmt"
)

func main() {
	container := Container{
		Length: 600,
		Width:  400,
		Height: 400,
	}
	items := []*Item{
		{
			Length: 380,
			Width:  320,
			Height: 100,
			Qty:    2,
		},
		{
			Length: 380,
			Width:  320,
			Height: 220,
			Qty:    1,
		},
		{
			Length: 380,
			Width:  320,
			Height: 200,
			Qty:    1,
		},
		{
			Length: 40,
			Width:  210,
			Height: 80,
			Qty:    12,
		},
	}
	result := CanPack(container, items)
	if result {
		fmt.Println("The items can fit in the container.")
	} else {
		fmt.Println("The items cannot fit in the container.")
	}
}
