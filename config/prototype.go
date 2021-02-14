package config

import "wicklight/logger"

// Config is the prototype of config
type Config struct {
	Listen string

	TLS      TLSConfig
	Log      LogConfig
	Fallback FallbackConfig
	ACL      ACLConfig
	Users    []UserConfig
}

// TLSConfig is prototype for tls
type TLSConfig struct {
	Certificate string
	PrivateKey  string
	// NextProtos  []string
	// IssueHostname string
	// IssueStorage  string
}

// UserConfig is prototype for users
type UserConfig struct {
	Username string
	Password string
	// Admin    bool
	// Quota    int64
}

// FallbackConfig is prototype for fallback
type FallbackConfig struct {
	Target string

	Host      string
	WhiteList []string

	PACPath string
	PACFile string
	PACHost string
	PACPort string
}

// ACLConfig is prototype for access control
type ACLConfig struct {
	WhiteListMode bool
	AllowLocal    bool

	PortsList []string
	HostsList []string
}

// LogConfig is prototype for logger
type LogConfig struct {
	Level logger.LogLevel
	File  string
}
