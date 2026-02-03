package auth

import (
	"encoding/base64"
	"strings"
)

// Authenticator 处理 Basic 认证逻辑
type Authenticator struct {
	username string
	password string
}

// NewAuthenticator 创建一个新的认证器
func NewAuthenticator(username, password string) *Authenticator {
	return &Authenticator{
		username: username,
		password: password,
	}
}

// Verify 验证 Base64 编码的凭证
// 预期格式为 "username:password" 的 Base64 编码
func (a *Authenticator) Verify(encoded string) bool {
	if encoded == "" {
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return false
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return false
	}

	return parts[0] == a.username && parts[1] == a.password
}

// VerifyHTTPHeader 从 HTTP Authorization 头部验证
func (a *Authenticator) VerifyHTTPHeader(header string) bool {
	if !strings.HasPrefix(header, "Basic ") {
		return false
	}
	return a.Verify(header[6:])
}
