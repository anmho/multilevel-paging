package memory

type SegmentInfo struct {
	SegmentIndex int32
	FrameIndex   int32
	Size         int32 // size in words
}
