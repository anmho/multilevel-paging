package memory

type VirtualAddress struct {
	SegmentIndex uint32 // s
	PageIndex    uint32 // p
	Offset       uint32 // w
	Bound        uint32 // pw
}

func NewVirtualAddress(address uint32) VirtualAddress {
	segmentIndex := address >> 18       // 14 bit mask on left side
	offset := address & 0x1FF           // 9 bit mask
	pageIndex := (address >> 9) & 0x1FF // 9 bit mask, 2nd group
	pw := address & 0x3FFFF             // 18 bit mask offset and page index

	return VirtualAddress{
		SegmentIndex: segmentIndex,
		PageIndex:    pageIndex,
		Offset:       offset,
		Bound:        pw,
	}
}
