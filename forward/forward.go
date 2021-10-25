package forward

import (
	"errors"
	"io"
	"net"
)

var ErrNetEOF = errors.New("ErrNetEof")

type Forward struct {
	src       net.Conn
	dst       net.Conn
	errHandle func(error)
}

func NewForward(dst net.Conn, src net.Conn, onError func(error)) *Forward {
	return &Forward{
		src:       src,
		dst:       dst,
		errHandle: onError,
	}
}

func (f *Forward) Start() {
	go func() {
		var err error
		for {
			err = communicate(f.dst, f.src)
			if err != nil {
				f.onError(err)
				break
			}
		}
	}()

	var err error
	for {
		err = communicate(f.src, f.dst)
		if err != nil {
			f.onError(err)
			break
		}
	}
	f.Stop()
}

func (f *Forward) Stop() {
	var err error
	err = f.src.Close()
	if err != nil {
		f.onError(err)
	}

	err = f.dst.Close()
	if err != nil {
		f.onError(err)
	}
}

func (f *Forward) ErrHandle(cb func(error)) {
	f.errHandle = cb
}

func (f *Forward) onError(err error) {
	if f.errHandle != nil {
		f.errHandle(err)
	}
}

func communicate(src net.Conn, dst net.Conn) error {
	var n int64
	var err error

	for {
		n, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
		if n == 0 {
			return ErrNetEOF
		}
	}
}
