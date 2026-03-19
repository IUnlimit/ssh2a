package honeypot

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/IUnlimit/ssh2a/internal/db"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// 常见的 SSH host key 路径，按优先级排列
var sshHostKeyPaths = []string{
	"/etc/ssh/ssh_host_ed25519_key",
	"/etc/ssh/ssh_host_ecdsa_key",
	"/etc/ssh/ssh_host_rsa_key",
}

var cachedHostKey ssh.Signer

// Init 初始化蜜罐，加载 host key
// 优先读取本机 sshd 的 host key，保证蜜罐和转发使用同一指纹
func Init() {
	key, source, err := loadHostKey()
	if err != nil {
		log.Fatalf("[Honeypot] Failed to load any SSH host key: %v", err)
	}
	cachedHostKey = key
	log.Infof("[Honeypot] Using host key from %s", source)
}

// HandleConnection 处理蜜罐 SSH 连接
// 接受连接，记录用户名和密码，然后拒绝认证
func HandleConnection(conn net.Conn, ip string) {
	defer conn.Close()

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			username := c.User()
			password := string(pass)
			log.Warnf("[Honeypot] IP=%s username=%s password=%s", ip, username, password)

			// 记录到数据库
			db.RecordHoneypotCredential(ip, username, password)
			db.RecordSSHAttempt(ip, "honeypot")

			// 始终拒绝
			return nil, fmt.Errorf("password rejected")
		},
		// 允许最多3次尝试以捕获更多凭据
		MaxAuthTries: 3,
	}

	config.AddHostKey(cachedHostKey)

	// 进行 SSH 握手，这会触发 PasswordCallback
	_, _, _, err := ssh.NewServerConn(conn, config)
	if err != nil {
		// 预期行为：客户端认证失败后断开
		log.Debugf("[Honeypot] Connection from %s ended: %v", ip, err)
	}
}

// loadHostKey 按优先级尝试加载本机 sshd 的 host key
func loadHostKey() (ssh.Signer, string, error) {
	var permDenied []string

	tryLoad := func(p string) (ssh.Signer, bool) {
		data, err := os.ReadFile(p)
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				permDenied = append(permDenied, p)
			}
			return nil, false
		}
		signer, err := ssh.ParsePrivateKey(data)
		if err != nil {
			log.Warnf("[Honeypot] Found %s but failed to parse: %v", p, err)
			return nil, false
		}
		return signer, true
	}

	// 1. 优先尝试本机 sshd 的 host key
	for _, p := range sshHostKeyPaths {
		if signer, ok := tryLoad(p); ok {
			return signer, p, nil
		}
	}

	// 2. 尝试 /etc/ssh/ 下任意 host key 文件
	matches, _ := filepath.Glob("/etc/ssh/ssh_host_*_key")
	for _, p := range matches {
		if signer, ok := tryLoad(p); ok {
			return signer, p, nil
		}
	}

	if len(permDenied) > 0 {
		return nil, "", fmt.Errorf(
			"permission denied reading SSH host key(s): %v — run with sudo or grant read access",
			permDenied,
		)
	}
	return nil, "", fmt.Errorf("no SSH host key found in /etc/ssh/")
}
