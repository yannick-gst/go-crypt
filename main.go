package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/akamensky/argparse"
	"golang.org/x/term"
)

func passwordPrompt() (string, error) {
	var str string
	fmt.Print("Enter password:")
	for {
		bytes, err := term.ReadPassword(int(syscall.Stdin))
		str = string(bytes)
		if err != nil {
			return "", err
		}
		if str != "" {
			break
		}
	}
	fmt.Println()
	return str, nil
}

func main() {
	parser := argparse.NewParser(
		"gocrypt",
		"A command line tool to encrypt a file, or decrypt " +
		"a file that was encrypted with this tool.")
	inputFile := parser.String("i", "input", &argparse.Options{Required: true, Help: "Path to input file"})
	outputFile := parser.String("o", "output", &argparse.Options{Required: true, Help: "Path to output file"})
	mode := parser.Selector(
		"m", "mode",
		[]string{"encrypt", "decrypt"},
		&argparse.Options{Required: true, Help: "Specifies whether to encrypt or decrypt"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	pwd, err := passwordPrompt()
	if err != nil {
		fmt.Printf("%q\n", err)
	}

	if *mode == "encrypt" {
		if err := encryptFile(inputFile, outputFile, pwd); err != nil {
			fmt.Println("An error occured while encrypting the input file: ", err)
		} else {
			fmt.Println("Encryption succeeded.")
		}
	} else {
		if err := decryptFile(inputFile, outputFile, pwd); err != nil {
			fmt.Println("An error occured while decrypting the input file: ", err)
		} else {
			fmt.Println("Decryption succeeded.")
		}
	}
}
