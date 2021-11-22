package forward

import (
	"errors"
	"io"
	"time"

	"github.com/AlpsMonaco/proxy/util"
)

var ErrConnClosed = errors.New("remote connection is closed")
var ErrRWSizeDismatch = errors.New("read size not equal to write size")

type Forward struct {
	SrcConn io.ReadWriteCloser
	DstConn io.ReadWriteCloser
	OnError func(error)
}

func NewForward(dst io.ReadWriteCloser, src io.ReadWriteCloser, onError func(error)) *Forward {
	return &Forward{
		SrcConn: src,
		DstConn: dst,
		OnError: onError,
	}
}

func (f *Forward) Start() {
	go func() {
		var err error
		for {
			err = communicate(f.DstConn, f.SrcConn)
			if err != nil {
				f.onError(err)
				break
			}
		}
	}()

	var err error
	for {
		err = communicate(f.SrcConn, f.DstConn)
		if err != nil {
			f.onError(err)
			break
		}
	}
	// f.Stop()
}

func (f *Forward) Stop() {
	closeConn(f.SrcConn)
	closeConn(f.DstConn)
}

func (f *Forward) SetErrHandle(cb func(error)) {
	f.OnError = cb
}

func (f *Forward) onError(err error) {
	if f.OnError != nil {
		f.OnError(err)
	}
}

const defaultNetBufSize = 0xFFFF >> 1

func communicate(src io.ReadWriteCloser, dst io.ReadWriteCloser) error {
	var a *util.Allocator = util.GetAlloctor(defaultNetBufSize)
	defer func() {
		util.FreeAllocator(a)
		closeConn(src)
		closeConn(dst)
	}()

	var nr int
	var er, ew, err error

	for {
		nr, er = src.Read(a.GetBytes())
		if nr > 0 {
			_, ew = dst.Write(a.GetByteSize(nr))
			if ew != nil {
				err = ew
				break
			}
		} else {
			if err != nil {
				err = ErrConnClosed
			}
			break
		}

		if er != nil {
			err = er
			break
		}
	}

	return err
}

func closeConn(conn io.ReadWriteCloser) {
	time.Sleep(5 * time.Second)
	err := conn.Close()
	if err != nil {
		_ = conn.Close()
	}
}
