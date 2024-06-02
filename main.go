package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/anmho/cs143b/project2/memory"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	var initPath string
	var inputPath string
	var outputPath string

	flag.StringVar(&initPath, "init", "", "initialization file")
	flag.StringVar(&inputPath, "input", "", "input file")
	flag.StringVar(&outputPath, "output", "", "output file")
	flag.Parse()

	initFile, inputFile, outFile, err := getFiles(initPath, inputPath, outputPath)
	if err != nil {
		log.Fatalln(err)
	}

	segments, pages, err := parseInitFile(initFile)
	if err != nil {
		log.Fatalln(err)
	}

	physicalMemory := make([]int32, memory.PhysicalMemorySize)

	initializeSegmentTable(physicalMemory, segments)
	initializePageTables(physicalMemory, pages)

	processInput(physicalMemory, inputFile, outFile)
}

func processInput(physicalMemory []int32, in, out *os.File) {
	scanner := bufio.NewScanner(in)

	disk := memory.NewDisk()

	mm := memory.NewManager(physicalMemory, disk)

	for scanner.Scan() {
		line := scanner.Text()

		args := strings.Split(line, " ")

		fmt.Fprintln(os.Stdout, args)

		if len(args) == 0 {
			panic("expected at least 1 arg")
		}

		cmd := args[0]
		switch cmd {
		case "RP":
			if len(args) != 2 {
				panic("expected 2 args")
			}

			physicalAddr, err := strconv.Atoi(args[1])
			if err != nil || physicalAddr < 0 {
				panic("invalid virtual address")
			}

			data := mm.ReadPhysical(int32(physicalAddr))
			fmt.Fprintf(out, "%d ", data)

		case "TA":
			if len(args) != 2 {
				panic("expected 2 args")
			}
			virtualAddr, err := strconv.Atoi(args[1])
			if err != nil || virtualAddr < 0 {
				panic("invalid virtual address")
			}
			physicalAddr, err := mm.TranslateAddress(uint32(virtualAddr))
			fmt.Fprintf(out, "%d ", physicalAddr)
		case "NL":
			fmt.Fprintf(out, "\n")
		}

	}
}

func initializeSegmentTable(physicalMemory []int32, segments []memory.SegmentInfo) {
	for _, segment := range segments {
		// ex 6 3000 4
		// Segment 6 resides in frame 4, size of segment 6 is 3000
		// set the memory at 2 * segment index to the segment size

		physicalMemory[2*segment.SegmentIndex] = segment.Size
		// set the frame number of segment 6 to be
		physicalMemory[2*segment.SegmentIndex+1] = segment.FrameIndex
	}
}

func initializePageTables(physicalMemory []int32, pages []memory.PageInfo) {
	for _, page := range pages {
		// 6 5 9
		// page 5 of segment 5 resides in frame 9
		pageTableFrameIndex := int32(math.Abs(float64(physicalMemory[2*page.SegmentIndex+1])))
		physicalMemory[pageTableFrameIndex*memory.FrameSize+page.PageIndex] = page.FrameIndex
	}
}

func parseInitFile(initFile *os.File) ([]memory.SegmentInfo, []memory.PageInfo, error) {
	scanner := bufio.NewScanner(initFile)

	var segments []memory.SegmentInfo
	// Initialize Segment Table
	if scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		for i := 2; i < len(parts); i += 3 {
			segmentIndex, err := strconv.Atoi(parts[i-2])
			if err != nil {
				return nil, nil, fmt.Errorf("parsing segment index: %w", err)
			}
			segmentSize, err := strconv.Atoi(parts[i-1])
			if err != nil {
				return nil, nil, fmt.Errorf("parsing segment size: %w", err)
			}
			pageTableFrameIndex, err := strconv.Atoi(parts[i])
			if err != nil {
				return nil, nil, fmt.Errorf("parsing frame index: %w", err)
			}
			if pageTableFrameIndex < 0 {
				log.Println("this page of this segment is not residing in actual physical memory", pageTableFrameIndex)

			}

			segmentInfo := memory.SegmentInfo{
				SegmentIndex: int32(segmentIndex),
				FrameIndex:   int32(pageTableFrameIndex),
				Size:         int32(segmentSize),
			}
			segments = append(segments, segmentInfo)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("scanning segments (first) line of init file: %w", err)
	}

	// Initialize Page Tables
	var pages []memory.PageInfo
	if scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, " ")

		for i := 2; i < len(parts); i += 3 {

			// ex: 6 5 9
			// page 5 of segment 6 resides in frame 9
			segmentIndex, err := strconv.Atoi(parts[i-2])
			if err != nil {
				return nil, nil, fmt.Errorf("parsing segment index: %w", err)
			}
			pageIndex, err := strconv.Atoi(parts[i-1])
			if err != nil {
				return nil, nil, fmt.Errorf("parsing page index: %w", err)
			}
			pageFrameIndex, err := strconv.Atoi(parts[i])
			if err != nil {
				return nil, nil, fmt.Errorf("parsing page frame index: %w", err)
			}

			pageInfo := memory.PageInfo{
				PageIndex:    int32(pageIndex),
				SegmentIndex: int32(segmentIndex),
				FrameIndex:   int32(pageFrameIndex),
			}

			pages = append(pages, pageInfo)
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("scanning pages (second) line of init file: %w", err)
	}
	return segments, pages, nil

}

func getFiles(initPath, inputPath, outputPath string) (*os.File, *os.File, *os.File, error) {
	if initPath == "" || inputPath == "" || outputPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	initFile, err := os.Open(initPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("opening init file: %w", err)
	}

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("opening input file: %w", err)
	}

	outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("opening output file: %w", err)
	}

	return initFile, inputFile, outputFile, nil
}
