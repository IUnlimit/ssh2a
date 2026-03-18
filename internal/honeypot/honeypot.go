package honeypot

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"

	"github.com/IUnlimit/ssh2a/internal/db"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

const hostKeyPath = "ssh2a_host_key"

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

	hostKey, err := loadOrGenerateHostKey()
	if err != nil {
		log.Errorf("[Honeypot] Failed to load host key: %v", err)
		return
	}
	config.AddHostKey(hostKey)

	// 进行 SSH 握手，这会触发 PasswordCallback
	_, _, _, err = ssh.NewServerConn(conn, config)
	if err != nil {
		// 预期行为：客户端认证失败后断开
		log.Debugf("[Honeypot] Connection from %s ended: %v", ip, err)
	}
}

// loadOrGenerateHostKey 加载或生成 SSH host key
func loadOrGenerateHostKey() (ssh.Signer, error) {
	// 尝试从文件加载
	if data, err := os.ReadFile(hostKeyPath); err == nil {
		return ssh.ParsePrivateKey(data)
	}

	// 生成新的 RSA key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// 保存到文件
	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	if err := os.WriteFile(hostKeyPath, pemData, 0600); err != nil {
		log.Warnf("[Honeypot] Failed to save host key: %v", err)
	}

	return ssh.NewSignerFromKey(key)
}
