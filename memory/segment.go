package memory

type SegmentInfo struct {
	Index               int32 // index of the segment
	pageTableFrameIndex int32 // frame which the associated page table resides in
	Size                int32 // size in words
}

func (s SegmentInfo) PageTableIsResident() bool {
	return s.pageTableFrameIndex > 0
}

func (s SegmentInfo) PageTableFrameIndexValue() int32 {
	return max(s.pageTableFrameIndex, -s.pageTableFrameIndex)
}
