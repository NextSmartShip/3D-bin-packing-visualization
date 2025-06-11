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
			Length: 400,
			Width:  210,
			Height: 80,
			Qty:    2,
		},
	}
	result := CanPack(container, items)
	fmt.Println(result)
}
