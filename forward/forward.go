package forward

import (
	"errors"
	"io"
	"net"
)

var ErrNetClosed = errors.New("ErrNetClosed")

type Forward struct {
	SrcConn net.Conn
	DstConn net.Conn
	OnError func(error)
}

func NewForward(dst net.Conn, src net.Conn, onError func(error)) *Forward {
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

func communicate(src net.Conn, dst net.Conn) error {
	defer func() {
		closeConn(src)
		closeConn(dst)
	}()
	var n int64
	var err error

	for {
		n, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
		if n == 0 {
			return ErrNetClosed
		}
	}
}

func closeConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		_ = conn.Close()
	}
}
