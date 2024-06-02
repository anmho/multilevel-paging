package memory

import (
	"errors"
	"log"
	"math"
)

const (
	PhysicalMemorySize int32 = 524_288
	FrameSize          int32 = 512
)

type Manager struct {
	physicalMemory []int32
	disk           *Disk
}

func NewManager(physicalMemory []int32, disk *Disk) *Manager {
	m := Manager{
		physicalMemory: physicalMemory,
		disk:           disk,
	}
	return &m
}

// TranslateAddress translates virtual address to physical address
func (m *Manager) TranslateAddress(virtualAddress uint32) (uint32, error) {
	va := NewVirtualAddress(virtualAddress)

	if va.Bound >= uint32(math.Abs(float64(m.physicalMemory[2*va.SegmentIndex]))) {
		return 0, errors.New("invalid address")
	}
	//m.physicalMemory[va.SegmentIndex]

	// s: segment index
	// p: page index
	// w: page offset

	// the address can page fault on the segment/pagetable

	// or on the page entry int he page table

	segmentIndex := va.SegmentIndex
	//segmentSize := m.physicalMemory[2*va.SegmentIndex]

	if m.physicalMemory[2*segmentIndex+1] < 0 {
		log.Println("page fault on ", m.physicalMemory[2*segmentIndex+1])
	}
	pageTableFrameIndex := math.Abs(float64(m.physicalMemory[2*segmentIndex+1]))
	m.disk.ReadBlock(m.physicalMemory, pageTableFrameIndex, pageTableFrameIndex)

	// for each step make sure its not negative otherwise you need to load the frame from the disk
	physicalAddress := m.physicalMemory[pageTableFrameIndex*FrameSize+int32(va.PageIndex)]*FrameSize + int32(va.PageOffset)
	if physicalAddress < 0 {
		log.Println("page fault on ", physicalAddress)
	}

	return uint32(physicalAddress), nil
}

// ReadPhysical reads the data at a physical address
func (m *Manager) ReadPhysical(physicalAddress int32) int32 {
	if physicalAddress < 0 {
		log.Println("page fault", physicalAddress)
	}
	return m.physicalMemory[physicalAddress]
}
