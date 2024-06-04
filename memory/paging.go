package memory

type PageTableEntry struct {
	PageIndex      int32 // index of the page in its page table
	SegmentIndex   int32 // index of the segment it belongs to
	pageFrameIndex int32 // index of the associated page
}

func (p PageTableEntry) PageAddr() int32 {
	return p.pageFrameIndex*FrameSize + p.PageIndex

}

func (p PageTableEntry) PageFrameIndexValue() int32 {
	return max(p.pageFrameIndex, -p.pageFrameIndex)
}

func (p PageTableEntry) PageIsResident() bool {
	return p.pageFrameIndex > 0
}
