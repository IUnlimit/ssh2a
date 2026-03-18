package db

import (
	log "github.com/sirupsen/logrus"
)

// RecordSSHAttempt 记录 SSH 连接尝试
func RecordSSHAttempt(ip, action string) {
	if err := DB.Create(&SSHAttempt{IP: ip, Action: action}).Error; err != nil {
		log.Errorf("Failed to record SSH attempt: %v", err)
	}
}

// RecordHoneypotCredential 记录蜜罐捕获的凭据
func RecordHoneypotCredential(ip, username, password string) {
	if err := DB.Create(&HoneypotCredential{
		IP: ip, Username: username, Password: password,
	}).Error; err != nil {
		log.Errorf("Failed to record honeypot credential: %v", err)
	}
}

// RecordAuth 记录认证尝试
func RecordAuth(ip, method string, success bool) {
	if err := DB.Create(&AuthRecord{
		IP: ip, Method: method, Success: success,
	}).Error; err != nil {
		log.Errorf("Failed to record auth: %v", err)
	}
}

// CredentialStat 凭据统计
type CredentialStat struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Count    int64  `json:"count"`
}

// GetHoneypotCredentials 获取蜜罐凭据统计
func GetHoneypotCredentials(page, pageSize int) ([]CredentialStat, int64, error) {
	var stats []CredentialStat
	var total int64

	sub := DB.Model(&HoneypotCredential{}).
		Select("username, password, count(*) as count").
		Group("username, password")

	if err := DB.Table("(?) as t", sub).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := sub.Order("count desc").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&stats).Error
	return stats, total, err
}

// IPStat IP 统计
type IPStat struct {
	IP    string `json:"ip"`
	Count int64  `json:"count"`
}

// GetRejectedIPs 获取被拒绝的 IP 列表及次数
func GetRejectedIPs(page, pageSize int) ([]IPStat, int64, error) {
	var stats []IPStat
	var total int64

	sub := DB.Model(&SSHAttempt{}).
		Select("ip, count(*) as count").
		Where("action = ?", "rejected").
		Group("ip")

	if err := DB.Table("(?) as t", sub).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := sub.Order("count desc").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&stats).Error
	return stats, total, err
}

// VerifiedIP 已验证 IP
type VerifiedIP struct {
	IP        string `json:"ip"`
	Method    string `json:"method"`
	CreatedAt string `json:"created_at"`
}

// GetVerifiedIPs 获取已验证的 IP 列表
func GetVerifiedIPs(page, pageSize int) ([]VerifiedIP, int64, error) {
	var records []VerifiedIP
	var total int64

	query := DB.Model(&AuthRecord{}).Where("success = ?", true)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Select("ip, method, created_at").
		Order("created_at desc").
		Offset((page - 1) * pageSize).Limit(pageSize).
		Find(&records).Error
	return records, total, err
}

// Stats 总览统计
type Stats struct {
	TotalSSHAttempts     int64 `json:"total_ssh_attempts"`
	TotalRejected        int64 `json:"total_rejected"`
	TotalHoneypot        int64 `json:"total_honeypot"`
	TotalForwarded       int64 `json:"total_forwarded"`
	TotalAuthAttempts    int64 `json:"total_auth_attempts"`
	TotalAuthSuccess     int64 `json:"total_auth_success"`
	UniqueHoneypotCreds  int64 `json:"unique_honeypot_creds"`
	UniqueRejectedIPs    int64 `json:"unique_rejected_ips"`
}

// GetStats 获取总览统计
func GetStats() (*Stats, error) {
	s := &Stats{}

	DB.Model(&SSHAttempt{}).Count(&s.TotalSSHAttempts)
	DB.Model(&SSHAttempt{}).Where("action = ?", "rejected").Count(&s.TotalRejected)
	DB.Model(&SSHAttempt{}).Where("action = ?", "honeypot").Count(&s.TotalHoneypot)
	DB.Model(&SSHAttempt{}).Where("action = ?", "forwarded").Count(&s.TotalForwarded)
	DB.Model(&AuthRecord{}).Count(&s.TotalAuthAttempts)
	DB.Model(&AuthRecord{}).Where("success = ?", true).Count(&s.TotalAuthSuccess)

	DB.Model(&HoneypotCredential{}).
		Select("count(distinct(username || '::' || password))").
		Scan(&s.UniqueHoneypotCreds)
	DB.Model(&SSHAttempt{}).
		Where("action = ?", "rejected").
		Select("count(distinct ip)").
		Scan(&s.UniqueRejectedIPs)

	return s, nil
}

// HasSSHAttempt 检查 IP 是否有 SSH 访问记录
func HasSSHAttempt(ip string) bool {
	var count int64
	DB.Model(&SSHAttempt{}).Where("ip = ?", ip).Count(&count)
	return count > 0
}
