package rssh

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/IUnlimit/ssh2a/internal/cache"
	"github.com/IUnlimit/ssh2a/internal/db"
	"github.com/IUnlimit/ssh2a/internal/honeypot"
	"github.com/libp2p/go-reuseport"
	log "github.com/sirupsen/logrus"
)

var Cache *cache.IPCache

func Listen(host string, port int, ipCache *cache.IPCache, wg *sync.WaitGroup) {
	defer wg.Done()
	Cache = ipCache

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Infof("SSH forward service starting on %s", addr)

	listener, err := reuseport.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on SSH port %s: %v", addr, err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Errorf("Failed to accept SSH connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Errorf("Error parsing remote address: %v", err)
		return
	}

	state := Cache.CheckSSH(ip)

	switch state {
	case cache.StateVerified:
		// 已验证，转发到本机 22 端口
		db.RecordSSHAttempt(ip, "forwarded")
		log.Infof("Forwarding SSH connection from verified IP %s", ip)
		forwardToSSH(conn, ip)

	case cache.StateRejected:
		// 首次拒绝，记录并关闭
		db.RecordSSHAttempt(ip, "rejected")
		log.Warnf("SSH connection rejected from IP %s (not verified)", ip)

	case cache.StateHoneypot:
		// 进入蜜罐
		log.Warnf("SSH connection from IP %s entering honeypot", ip)
		honeypot.HandleConnection(conn, ip)

	default:
		db.RecordSSHAttempt(ip, "rejected")
		log.Warnf("SSH connection rejected from IP %s (unknown state)", ip)
	}
}

func forwardToSSH(clientConn net.Conn, ip string) {
	targetConn, err := net.DialTimeout("tcp", "127.0.0.1:22", 3*time.Second)
	if err != nil {
		log.Errorf("Failed to connect to local SSH server: %v", err)
		return
	}
	defer targetConn.Close()

	log.Infof("Forwarding connection from %s to SSH server", ip)

	// 双向转发
	done := make(chan struct{}, 2)
	go func() {
		io.Copy(targetConn, clientConn)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(clientConn, targetConn)
		done <- struct{}{}
	}()
	<-done
}
