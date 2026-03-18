package cache

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// IPState IP 的当前状态
type IPState int

const (
	StateUnknown   IPState = iota // 未知/首次
	StateRejected                 // 已被拒绝，等待验证
	StateHoneypot                 // 进入蜜罐模式
	StateVerified                 // 已验证，放行中
)

// ipEntry 单个 IP 的缓存条目
type ipEntry struct {
	State         IPState
	FirstRejectAt time.Time // 首次被拒绝的时间
	RejectCount   int       // 被拒绝次数
	VerifiedAt    time.Time // 验证通过的时间
}

// IPCache IP 状态缓存
type IPCache struct {
	mu              sync.RWMutex
	entries         map[string]*ipEntry
	triggerTimeout  time.Duration // 蜜罐触发超时
	validDuration   time.Duration // 验证有效期
}

// NewIPCache 创建 IP 缓存
func NewIPCache(triggerTimeout, validDuration time.Duration) *IPCache {
	c := &IPCache{
		entries:        make(map[string]*ipEntry),
		triggerTimeout: triggerTimeout,
		validDuration:  validDuration,
	}
	// 定期清理过期条目
	go c.cleanup()
	return c
}

// CheckSSH 检查 SSH 连接的 IP 状态，返回应执行的动作
func (c *IPCache) CheckSSH(ip string) IPState {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.entries[ip]

	// 已验证且未过期 → 放行
	if exists && entry.State == StateVerified {
		if time.Since(entry.VerifiedAt) < c.validDuration {
			return StateVerified
		}
		// 验证已过期，重置
		delete(c.entries, ip)
		log.Infof("IP %s verification expired", ip)
		entry = nil
		exists = false
	}

	// 首次连接 → 拒绝并记录
	if !exists {
		c.entries[ip] = &ipEntry{
			State:         StateRejected,
			FirstRejectAt: time.Now(),
			RejectCount:   1,
		}
		return StateRejected
	}

	// 已在蜜罐模式
	if entry.State == StateHoneypot {
		entry.RejectCount++
		return StateHoneypot
	}

	// 已被拒绝，检查是否应进入蜜罐
	entry.RejectCount++
	if time.Since(entry.FirstRejectAt) < c.triggerTimeout {
		// 在超时窗口内重试 → 进入蜜罐
		entry.State = StateHoneypot
		log.Warnf("IP %s entered honeypot mode after %d retries within %v",
			ip, entry.RejectCount, c.triggerTimeout)
		return StateHoneypot
	}

	// 超时窗口已过但未验证，重置窗口
	entry.FirstRejectAt = time.Now()
	entry.RejectCount = 1
	return StateRejected
}

// SetVerified 标记 IP 为已验证
func (c *IPCache) SetVerified(ip string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[ip] = &ipEntry{
		State:      StateVerified,
		VerifiedAt: time.Now(),
	}
	log.Infof("IP %s verified, SSH access granted for %v", ip, c.validDuration)
}

// SetRejected 标记 IP 认证失败
func (c *IPCache) SetRejected(ip string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.entries[ip]
	if !exists {
		c.entries[ip] = &ipEntry{
			State:         StateRejected,
			FirstRejectAt: time.Now(),
			RejectCount:   1,
		}
		return
	}
	if entry.State != StateVerified {
		entry.RejectCount++
	}
}

// IsVerified 检查 IP 是否已验证且未过期
func (c *IPCache) IsVerified(ip string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[ip]
	if !exists {
		return false
	}
	return entry.State == StateVerified && time.Since(entry.VerifiedAt) < c.validDuration
}

// HasRecord 检查 IP 是否有任何记录
func (c *IPCache) HasRecord(ip string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.entries[ip]
	return exists
}

// cleanup 定期清理过期条目
func (c *IPCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for ip, entry := range c.entries {
			switch entry.State {
			case StateVerified:
				if time.Since(entry.VerifiedAt) > c.validDuration {
					delete(c.entries, ip)
				}
			case StateRejected, StateHoneypot:
				// 24小时无活动则清理
				if time.Since(entry.FirstRejectAt) > 24*time.Hour {
					delete(c.entries, ip)
				}
			}
		}
		c.mu.Unlock()
	}
}
