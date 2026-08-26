package main

import (
	"fmt"
	"os"
	"tetris-optimizer/tetris"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Please insert only one Argument.")
		return
	}

	if err := tetris.Run(os.Args[1]); err != nil {
		fmt.Println(err)
	}
}
