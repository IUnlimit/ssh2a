package ssh

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

func Init(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	defer listener.Close()
	if err != nil {
		log.Fatalf("SSH server exists with error: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Warnf("SSH listening error: %v", err)
			continue
		}
		go func() {
			handleConn(conn)
			defer conn.Close()
		}()
	}
}
