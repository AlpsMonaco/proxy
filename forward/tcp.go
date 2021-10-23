package forward

import "net"

type TCPForward struct {
	// addr to listen
	ListenAddr string
	// target host addr
	HostAddr string

	// Error Handle
	OnError func(error)
}

func NewTCPForward(listenAddr, hostAddr string, onError func(error)) *TCPForward {
	return &TCPForward{
		ListenAddr: listenAddr,
		HostAddr:   hostAddr,
		OnError:    onError,
	}
}

func (tf *TCPForward) Start() error {
	listener, err := net.Listen("tcp", tf.ListenAddr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		dst, err := net.Dial("tcp", tf.HostAddr)
		if err != nil {
			tf.onError(err)
		} else {
			go NewForward(dst, conn, tf.OnError).Start()
		}
	}
}

func (tf *TCPForward) onError(err error) {
	if tf.OnError != nil {
		tf.OnError(err)
	}
}
