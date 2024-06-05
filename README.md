Precompiled executable platform binaries contained in ./bin/project1 (darwin-arm64 and linux-amd64)

# How to run
1) Prerequisites: Install Go 1.22.2
1) Dependencies: Install packages with `go mod download`
2) Build: 
   3) To build for current platform run ```make```
   4) To build only for linux run ```make linux-amd64```
   5) To build only for MacOS run ```make darwin-arm64```
3) To run the program enter command: ```./bin/project2 -i <init-file-path> < <input-file-path> > <output-file-path>```
- e.g.
  - To run the test input w/o demand paging
      - ```./bin/project2 -i ./rsrc/init-no-dp.txt < ./rsrc/input-no-dp.txt > ./rsrc/output-no-dp.txt 2> log.txt```
  - To run the test input with demand paging
       -```./bin/project2 -i ./rsrc/init-dp.txt < ./rsrc/input-dp.txt > ./rsrc/output-dp.txt 2> log.txt```

# A list of files in the zip file, with a short description of each one.
- `Makefile`:
- `bin`: Contains the binary executable generated from compilation.
- `go.mod`: Tracks the package dependencies
- `go.sum`: Stores the checksums for the dependencies to ensure correctness.
- `main.go`: Main entrypoint. Responsible for parsing inputs, routing commands, and initializing the program.
- `memory`: Package for memory related tasks for the segmentation and paging system.
  - `address.go`: Contains structs and utilities for virtual address translation.
  - `disk.go`: Contains `Disk` struct and read/write utilities.
  - `frame.go`: Contains `Frame` struct and methods.
  - `manager.go`: Contains the main memory manager that allocates frames, handles demand paging, address translation and memory reading.
  - `manager_test.go`: Contains integration tests for the memory manager .
  - `paging.go`: Contains `PageTableEntry` struct and additional paging utilities.
  - `segment.go`: Contains `SegmentInfo` struct and related utilities for demand paging.
- `rsrc`: Includes various sample init, input, output files, and project instructions 