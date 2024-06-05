package memory

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/zyedidia/generic/list"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	PhysicalMemorySize int32 = 524_288
	FrameSize          int32 = 512
)

type Manager struct {
	physicalMemory []int32
	isFree         []bool
	freeList       *list.List[FrameInfo]
	disk           *Disk
}

func NewManager() *Manager {

	isFree := make([]bool, PhysicalMemorySize/FrameSize)
	for i := 2; i < len(isFree); i++ {
		isFree[i] = true
	}
	m := Manager{
		physicalMemory: make([]int32, PhysicalMemorySize),
		isFree:         isFree,
		freeList:       list.New[FrameInfo](),
		disk:           NewDisk(),
	}

	return &m
}

func (m *Manager) InitializeFromFile(file *os.File) error {
	segments, pageTableEntries, err := parseInitFile(file)
	if err != nil {
		return fmt.Errorf("parsing init file: %w", err)
	}

	for _, segment := range segments {
		// Set the segment info
		m.physicalMemory[2*segment.Index] = segment.Size                  // segment size
		m.physicalMemory[2*segment.Index+1] = segment.pageTableFrameIndex // may or may not be negative
		if segment.PageTableIsResident() {
			m.isFree[segment.pageTableFrameIndex] = false
		}
	}

	for _, entry := range pageTableEntries {
		segmentInfo := m.getSegmentInfo(entry.SegmentIndex)
		//
		if segmentInfo.PageTableIsResident() {
			// the associated page table is in memory so write it to memory

			m.isFree[segmentInfo.PageTableFrameIndexValue()] = false

			m.physicalMemory[segmentInfo.PageTableFrameIndexValue()*FrameSize+entry.PageIndex] = entry.pageFrameIndex
		} else {
			// the associated page table is not in memory so write the page table entry to disk
			// write the page frame index to disk instead

			m.disk.Write(segmentInfo.PageTableFrameIndexValue(), entry.PageIndex, entry.pageFrameIndex)
		}
	}
	for frameIndex := range len(m.isFree) {
		if m.isFree[frameIndex] {
			frame := FrameInfo{addr: int32(frameIndex) * FrameSize}
			m.addFrame(frame)
		} else {
			//log.Println("frame not free", frameIndex)
		}

	}

	return nil
}

// TranslateAddress translates virtual address to physical address
func (m *Manager) TranslateAddress(virtualAddress uint32) (int32, error) {
	va := NewVirtualAddress(virtualAddress)
	log.Printf("virtual addr: %+v\n", va)

	segmentInfo := m.getSegmentInfo(int32(va.SegmentIndex))
	log.Printf("segment: %+v\n", segmentInfo)
	if int64(va.Bound) >= int64(segmentInfo.Size) {
		return -1, errors.New("bound size pw invalid")
	}

	if !segmentInfo.PageTableIsResident() {
		// load the frame/block from disk into the page table frame
		freeFrame := m.getFreeFrame()
		b := segmentInfo.PageTableFrameIndexValue()
		m.disk.ReadBlock(m.physicalMemory, b, freeFrame.Addr())
		m.setSegmentPageTableFrameIndex(segmentInfo.Index, freeFrame.Index())
	}
	segmentInfo = m.getSegmentInfo(int32(va.SegmentIndex))
	//if segmentInfo.PageTableFrameIndexValue() != 3 {
	//	panic("should be 3")
	//}

	log.Printf("segment after potential page fault: %+v\n", segmentInfo)

	pageTableEntry := m.getPageTableEntry(segmentInfo.Index, int32(va.PageIndex))
	//if pageTableEntry.pageFrameIndex != -20 {
	//	panic("should be -20")
	//}
	log.Printf("page table entry: %+v\n", pageTableEntry)
	if !pageTableEntry.PageIsResident() {
		// load from disk
		freeFrame := m.getFreeFrame()
		b := pageTableEntry.PageFrameIndexValue()
		m.disk.ReadBlock(m.physicalMemory, b, freeFrame.Addr())
		//m.setPageTableEntryFrame(segmentInfo.PageTableFrameIndexValue(), pageTableEntry.PageIndex, freeFrame.Index())
		log.Printf("lowest free frame: %d\n", freeFrame.Index())
		m.physicalMemory[segmentInfo.PageTableFrameIndexValue()*FrameSize+pageTableEntry.PageIndex] = freeFrame.Index()
	}
	log.Printf("page table entry after potential fault: %+v\n", pageTableEntry)

	s := segmentInfo.Index
	p := int32(va.PageIndex)
	w := int32(va.Offset)

	// PA = PM[PM[2s + 1]*512 + p]*512 + w
	//physicalAddress := m.physicalMemory[pageTableEntry.PageFrameIndexValue()*FrameSize+pageTableEntry.PageIndex] + int32(va.Offset)
	pageTableFrameIndex := m.physicalMemory[2*s+1]
	log.Println("page table frame index: ", pageTableFrameIndex)
	physicalAddress := m.physicalMemory[pageTableFrameIndex*FrameSize+p]*FrameSize + w

	return physicalAddress, nil
}

// ReadPhysical reads the data at a physical address
func (m *Manager) ReadPhysical(physicalAddress uint32) int32 {
	return m.physicalMemory[physicalAddress]
}

func (m *Manager) addFrame(frame FrameInfo) {
	m.isFree[frame.Index()] = true
	m.freeList.PushBack(frame)
}

func parseInitFile(initFile *os.File) ([]SegmentInfo, []PageTableEntry, error) {
	scanner := bufio.NewScanner(initFile)

	var segments []SegmentInfo
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

			segmentInfo := SegmentInfo{
				Index:               int32(segmentIndex),
				pageTableFrameIndex: int32(pageTableFrameIndex),
				Size:                int32(segmentSize),
			}
			segments = append(segments, segmentInfo)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("scanning segments (first) line of init file: %w", err)
	}

	// Initialize Page Tables
	var pages []PageTableEntry
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

			pageInfo := PageTableEntry{
				PageIndex:      int32(pageIndex),
				SegmentIndex:   int32(segmentIndex),
				pageFrameIndex: int32(pageFrameIndex),
			}

			pages = append(pages, pageInfo)
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("scanning pages (second) line of init file: %w", err)
	}
	return segments, pages, nil

}

// returns the starting address of the free frame
func (m *Manager) getFreeFrame() FrameInfo {
	if m.freeList.Front == nil {
		panic("out of free frames")
	}
	head := m.freeList.Front
	frame := head.Value
	m.freeList.Remove(head)
	m.isFree[frame.Index()] = false

	log.Printf("%+v\n", frame)
	return frame
}

func (m *Manager) getSegmentInfo(segmentIndex int32) SegmentInfo {
	segmentSize := m.physicalMemory[2*segmentIndex]
	pageTableFrameIndex := m.physicalMemory[2*segmentIndex+1]

	s := SegmentInfo{
		Index:               segmentIndex,
		pageTableFrameIndex: pageTableFrameIndex,
		Size:                segmentSize,
	}
	return s
}

func (m *Manager) getPageTableEntry(segmentIndex, pageIndex int32) PageTableEntry {
	segmentInfo := m.getSegmentInfo(segmentIndex)
	pageFrameIndex := m.physicalMemory[segmentInfo.PageTableFrameIndexValue()*FrameSize+pageIndex]

	log.Printf("segment %d info for page %d in frame %d for page table frame %d\n", segmentInfo.Index, pageIndex, pageFrameIndex, segmentInfo.pageTableFrameIndex)

	return PageTableEntry{
		PageIndex:      pageIndex,
		SegmentIndex:   segmentIndex,
		pageFrameIndex: pageFrameIndex,
	}
}

func (m *Manager) setPageTableEntryFrame(pageTableFrameIndex, pageIndex, frameIndex int32) {
	m.physicalMemory[pageTableFrameIndex*FrameSize+pageIndex] = frameIndex
}

func (m *Manager) setSegmentPageTableFrameIndex(segmentIndex, pageTableFrameIndex int32) {
	m.physicalMemory[2*segmentIndex+1] = pageTableFrameIndex
}
