package queues

import "io"

type Message struct {
	Body []byte
	pos  int
}

func (m *Message) Read(p []byte) (n int, err error) {
	if m.pos >= len(m.Body) {
		return 0, io.EOF
	}

	n = copy(p, m.Body[m.pos:])
	m.pos += n
	return
}

// Seek sets the offset for the next Read to offset, which should be relative to whence
func (m *Message) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		m.pos = int(offset)
	case io.SeekCurrent:
		m.pos += int(offset)
	case io.SeekEnd:
		m.pos = len(m.Body) + int(offset)
	}

	return int64(m.pos), nil
}
