# How to run
1) Install Go 1.22.2
1) Install packages with `go mod download`
2) To build run ```make```
3) To run the program enter command: ```./bin/project2 -i <init-file-path> < <input-file-path> > <output-file-path>```
   - e.g. ```./bin/project2 -i ./rsrc/init-dp.txt < ./rsrc/input-dp.txt > ./rsrc/output-dp.txt```
4) To run the test input w/o demand paging
   - ```./bin/project2 -i ./rsrc/init-no-dp.txt < ./rsrc/input-no-dp.txt > ./rsrc/output-no-dp.txt```
5) To run the test input with demand paging 
   -```./bin/project2 -i ./rsrc/init-dp.txt < ./rsrc/input-dp.txt > ./rsrc/output-dp.txt```

# A list of files in the zip file, with a short description of each one.
- `Makefile`:
- `bin`: Contains the binary executable generated from compilation.
- `go.mod`: Tracks the package dependencies
- `go.sum`: Stores the checksums for the dependencies
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