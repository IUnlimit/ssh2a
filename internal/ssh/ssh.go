package ssh

import (
	"github.com/IUnlimit/ssh2a/internal/cache"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
)

func handleConn(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	bytes := make([]byte, 1024)
	n, err := conn.Read(bytes)
	if err != nil {
		log.Errorf("Error occurred when read data from %s: %v", addr, err)
		return
	}
	request := string(bytes[:n])

	// ssh format: 'SSH-{2.0}-{name}'
	if !strings.HasPrefix(request, "SSH-") {
		return
	}
	log.Debugf("Received: '%q' from %s\n", request, addr)

	// 暂时只支持 ipv4
	addrSplits := strings.Split(addr, ":")
	if len(addrSplits) == 1 {
		log.Warnf("Unsupport connection addr: %s", addr)
		return
	}

	ip := addrSplits[0]
	if !cache.ContainsAccessMap(ip) {
		count := cache.AddTempCache(ip)
		log.Debugf("IP %s tried to establish ssh connection %d times", ip, count)
		return
	}
	forwardConn(request, conn)
}

func forwardConn(request string, clientConn net.Conn) {
	// 连接到目标服务器
	localSSH := "127.0.0.1:22"
	serverConn, err := net.Dial("tcp", localSSH)
	if err != nil {
		log.Errorf("Error connecting to destination: %v", err)
		return
	}
	defer serverConn.Close()

	// 开启goroutine从客户端读取数据并转发到目标服务器
	go func() {
		_, err := serverConn.Write([]byte(request))
		if err != nil {
			log.Errorf("Error occurred when start SSH request to server: %v", err)
			return
		}

		log.Infof("SSH forward start from %s to %s", clientConn.RemoteAddr().String(), localSSH)
		_, err = io.Copy(serverConn, clientConn)
		if err != nil {
			log.Errorf("Error copying from client to server: %v", err)
		}
	}()

	// block
	// 从目标服务器读取数据并转发到客户端
	_, err = io.Copy(clientConn, serverConn)
	if err != nil {
		log.Errorf("Error copying from server to client: %v", err)
	}
}
