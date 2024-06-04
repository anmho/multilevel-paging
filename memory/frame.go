package memory

type FrameInfo struct {
	addr int32
}

func (f FrameInfo) Addr() int32 {
	return f.addr
}

func (f FrameInfo) Index() int32 {
	return f.addr / FrameSize
}
