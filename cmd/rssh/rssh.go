package rssh

import (
	"errors"
	"fmt"
	"github.com/IUnlimit/ssh2a/cache"
	"github.com/bluele/gcache"
	"github.com/libp2p/go-reuseport"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
	"time"
)

func Listen(host string, port int, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Infof("SSH forward service starting on %s:%d", host, port)
	listener := multipleAbleHttpListen(host, port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func multipleAbleHttpListen(host string, port int) net.Listener {
	sshListen, err := reuseport.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("Error occurred when listening resuse-ssh port(%d) on host(%s)", port, host)
	}
	return sshListen
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	ip, _, err := net.SplitHostPort(remoteAddr) // ipv4 / ipv6
	if err != nil {
		log.Errorf("Error parsing remote address, %v", err)
		return
	}

	// 检查是否在白名单中
	checkWL := true
	err = cache.UpdateIPStatus(ip, checkWL)
	if err != nil && !errors.Is(err, gcache.KeyNotFoundError) {
		log.Errorf("Error updating ip cache, %v", err)
		return
	}
	if !checkWL {
		log.Warnf("IP %s is not in the whitelist, closing connection", ip)
		return
	}

	originAddress := fmt.Sprintf("%s:%d", "127.0.0.1", 22)
	targetConn, err := net.DialTimeout("tcp", originAddress, 3*time.Second)
	if err != nil {
		log.Errorf("Failed to connect to target SSH server: %v", err)
		return
	}

	log.Infof("Forwarding connection from %s to SSH server", ip)
	forwardConnection(conn, targetConn)
}

func forwardConnection(clientConn, targetConn net.Conn) {
	defer clientConn.Close()
	defer targetConn.Close()

	// Bidirectional forwarding
	go io.Copy(targetConn, clientConn)
	io.Copy(clientConn, targetConn)
}
