package main

import (
	"fmt"
)

func foo() {
	inputs, err := getNixFilesContents()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(inputs)
}

func openaiTest() {
	test := "Insteall neovim with a basic configuration"
	nixFile := runPrompt(test)
	fmt.Println(nixFile)
}

func main() {
	openaiTest()
}
