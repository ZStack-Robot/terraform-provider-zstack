// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/util/httputils"
)

type AuthType string

const (
	AuthTypeAccessKey   AuthType = "AccessKey"
	AuthTypeAccount     AuthType = "Account"
	AuthTypeAccountUser AuthType = "AccountUser"
	AuthTypeSession     AuthType = "Session"
)

const (
	WebZStackPort            = 5000
	defaultZStackPort        = 8080
	defaultZStackContextPath = "zstack"
)

type ZSConfig struct {
	hostname    string
	port        int
	contextPath string

	sessionId string

	accessKeyId     string
	accessKeySecret string

	accountName     string
	accountUserName string
	password        string

	authType AuthType

	retryInterval int // unit - second
	retryTimes    int

	readOnly  bool
	debug     bool
	proxyFunc httputils.TransportProxyFunc
}

func NewZSConfig(hostname string, port int, contextPath string) *ZSConfig {
	return &ZSConfig{
		hostname:      hostname,
		port:          port,
		contextPath:   contextPath,
		retryInterval: 2,
		retryTimes:    1800,
	}
}

func DefaultZSConfig(hostname string) *ZSConfig {
	return NewZSConfig(hostname, defaultZStackPort, defaultZStackContextPath)
}

func (config *ZSConfig) AccessKey(accessKeyId, accessKeySecret string) *ZSConfig {
	config.accessKeyId = accessKeyId
	config.accessKeySecret = accessKeySecret
	config.authType = AuthTypeAccessKey
	return config
}

func (config *ZSConfig) Session(sessionId string) *ZSConfig {
	config.sessionId = sessionId
	config.authType = AuthTypeSession
	return config
}

func (config *ZSConfig) LoginAccount(accountName, password string) *ZSConfig {
	config.accountName = accountName
	config.password = password
	config.authType = AuthTypeAccount
	return config
}
func (config *ZSConfig) LoginAccountUser(accountName, accountUserName, password string) *ZSConfig {
	config.accountName = accountName
	config.accountUserName = accountUserName
	config.password = password
	config.authType = AuthTypeAccountUser
	return config
}

func (config *ZSConfig) RetryInterval(retryInterval int) *ZSConfig {
	config.retryInterval = retryInterval
	return config
}

func (config *ZSConfig) RetryTimes(retryTimes int) *ZSConfig {
	config.retryTimes = retryTimes
	return config
}

func (config *ZSConfig) ReadOnly(readOnly bool) *ZSConfig {
	config.readOnly = readOnly
	return config
}

func (config *ZSConfig) Debug(debug bool) *ZSConfig {
	config.debug = debug
	return config
}

func (config *ZSConfig) ProxyFunc(proxyFunc httputils.TransportProxyFunc) *ZSConfig {
	config.proxyFunc = proxyFunc
	return config
}
