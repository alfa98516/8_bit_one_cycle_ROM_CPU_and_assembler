package main

import (
	"fmt"
	"os"

	. "Assembler/backend"
)

func main() {
	argc := len(os.Args)
	if argc < 2 {
		fmt.Printf("Not enough arguments supplied.\n")
		usage()
		os.Exit(1)
	}
	src, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	var ofile *os.File
	if argc == 2 {
		ofile, err = os.Create("a.out")
	} else {
		ofile, err = os.Create(os.Args[3])
	}
	if err != nil {
		fmt.Println("Error creating file:", err)
	}
	defer ofile.Close()
	var assembler Assembler
	assembler.Init(ofile, src)
}

func usage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("./cosmo [INPUT FILE]")
	fmt.Printf("./cosmo [INPUT FILE] -o [OUTPUT FILE]")
}
