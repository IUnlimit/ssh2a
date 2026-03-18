package db

import "time"

// SSHAttempt 记录所有 SSH 连接尝试
type SSHAttempt struct {
	ID        uint      `gorm:"primaryKey"`
	IP        string    `gorm:"index;size:45"`
	Action    string    `gorm:"size:20"` // rejected, honeypot, forwarded
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// HoneypotCredential 蜜罐捕获的用户名密码
type HoneypotCredential struct {
	ID        uint      `gorm:"primaryKey"`
	IP        string    `gorm:"size:45"`
	Username  string    `gorm:"size:255"`
	Password  string    `gorm:"size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// AuthRecord 认证记录
type AuthRecord struct {
	ID        uint      `gorm:"primaryKey"`
	IP        string    `gorm:"index;size:45"`
	Method    string    `gorm:"size:20"` // basic, authenticator
	Success   bool
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
