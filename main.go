package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/anmho/cs143b/project2/memory"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var initPath string
	var inputPath string
	var outputPath string

	flag.StringVar(&initPath, "i", "", "initialization file")
	//flag.StringVar(&inputPath, "i", "", "input file")
	//flag.StringVar(&outputPath, "o", "", "output file")
	flag.Parse()

	initFile, err := getFiles(initPath)
	if err != nil {
		log.Fatalln(fmt.Errorf("error getting init file: %w", err))
	}

	mm := memory.NewManager()

	err = mm.InitializeFromFile(initFile)
	if err != nil {
		log.Fatalln(fmt.Errorf("error initializing: %w", err))
	}

	//initializeSegmentTable(mm, segments)
	//initializePageTables(mm, pages)

	processInput(mm, os.Stdin, os.Stdout)
}

func processInput(mm *memory.Manager, in, out *os.File) {
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		line := scanner.Text()

		args := strings.Split(line, " ")

		if len(args) == 0 {
			panic("expected at least 1 arg")
		}

		cmd := args[0]
		switch cmd {
		case "RP":
			if len(args) < 2 {
				log.Fatalln("expected 2 args", args)
			}

			physicalAddr, err := strconv.Atoi(args[1])
			if err != nil || physicalAddr < 0 {
				log.Fatalln("invalid virtual address", args[1])
			}

			data := mm.ReadPhysical(uint32(physicalAddr))
			fmt.Fprintf(out, "%d ", data)

		case "TA":
			if len(args) < 2 {
				log.Fatalln("expected 2 args", args)
			}
			virtualAddr, err := strconv.Atoi(args[1])
			if err != nil || virtualAddr < 0 {
				log.Fatalln("invalid virtual address", args)
			}
			physicalAddr, err := mm.TranslateAddress(uint32(virtualAddr))
			fmt.Fprintf(out, "%d ", physicalAddr)
		case "NL":
			fmt.Fprintf(out, "\n")
		}

	}
}

func getFiles(initPath string) (*os.File, error) {
	if initPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	initFile, err := os.Open(initPath)
	if err != nil {
		return nil, fmt.Errorf("opening init file: %w", err)
	}
	return initFile, nil
}
