package memory

type Disk struct {
	blocks [PhysicalMemorySize / FrameSize][FrameSize]int32
}

func NewDisk() *Disk {
	return &Disk{}
}

func (d *Disk) ReadBlock(physicalMemory []int32, blockIndex, memStart int32) {
	block := d.blocks[blockIndex]
	for i := range FrameSize {
		physicalMemory[memStart+i] = block[i]
	}
}
