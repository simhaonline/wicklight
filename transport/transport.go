package transport

import (
	"io"
	"net"
	"net/http"
	"sync"
)

// BufferPool is for data copy
var BufferPool sync.Pool

// PoolInit to initialize BufferPool
func init() {
	makeBuffer := func() interface{} { return make([]byte, 0, 65536) }
	BufferPool = sync.Pool{New: makeBuffer}
}

type closeWriter interface {
	CloseWrite() error
}

// Relay double side copy
func Relay(target net.Conn, clientReader io.ReadCloser, clientWriter io.Writer) int64 {
	stream := func(w io.Writer, r io.Reader) int64 {
		// copy bytes from r to w
		buf := BufferPool.Get().([]byte)
		buf = buf[0:cap(buf)]
		n, _ := flushingIoCopy(w, r, buf)
		BufferPool.Put(buf)
		if cw, ok := w.(closeWriter); ok {
			cw.CloseWrite()
		}
		return n
	}

	go stream(target, clientReader)
	return stream(clientWriter, target)
}

func flushingIoCopy(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	flusher, hasFlusher := dst.(http.Flusher)
	for {
		var nr int
		var er error
		nr, er = src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if hasFlusher {
				flusher.Flush()
			}
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}

	return
}
