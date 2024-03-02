package notifier

import (
	"TraderHelperCore/api"
	"fmt"
	"net"
)

type tcpNotifier struct {
	port   int
	client net.Conn
}

func NewTcpNotifier(port int) api.Notifier {
	notifier := &tcpNotifier{
		port:   port,
		client: nil,
	}

	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			notifier.client = conn
			_, _ = conn.Write([]byte("Connected\n"))
		}
	}()

	return notifier
}

func (sn tcpNotifier) Notify(title string, content string) {
	if sn.client == nil {
		return
	}
	_, _ = sn.client.Write([]byte(fmt.Sprintf("%s&&%s\n", title, content)))
}
