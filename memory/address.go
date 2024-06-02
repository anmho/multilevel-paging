package memory

type VirtualAddress struct {
	SegmentIndex uint32 // s
	PageIndex    uint32 // p
	PageOffset   uint32 // w
	Bound        uint32 // Bound
}

func NewVirtualAddress(address uint32) VirtualAddress {
	segmentIndex := address >> 18
	offset := address & 0x1FF
	pageNum := (address >> 9) & 0x1FF
	pw := address & 0x3FFFF

	return VirtualAddress{
		SegmentIndex: segmentIndex,
		PageIndex:    pageNum,
		PageOffset:   offset,
		Bound:        pw,
	}
}
