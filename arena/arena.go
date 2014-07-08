package arena

type Arena struct {
	buf []byte

	size int
}

func NewArena(size int) *Arena {
	a := new(Arena)

	a.size = size
	a.buf = make([]byte, size, size)

	return a
}

func (a *Arena) Make(size int) []byte {
	if a.size < size {
		return make([]byte, size)
	} else if len(a.buf) < size {
		a.buf = make([]byte, a.size)
	}

	b := a.buf[0:size]
	a.buf = a.buf[size:]
	return b
}
