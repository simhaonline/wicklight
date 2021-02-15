package server

import (
	"encoding/base64"
	"net"
	"net/http"
	"strings"
	"wicklight/config"
)

func checkHost(hr *http.Request) (host, port string) {
	rawHost := hr.Host
	if rawHost == "" {
		rawHost = hr.URL.Host
	}
	host, port, err := net.SplitHostPort(rawHost)
	if err != nil {
		host = hr.Host
		if hr.URL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	return
}

func checkUser(hr *http.Request) (username string, authenticated bool) {
	if len(config.Conf.Users) == 0 {
		return "anonymous", true
	}
	credentials := hr.Header.Get("Proxy-Authorization")
	reqUsername, reqPassword, ok := parseBasicAuth(credentials)
	if ok {
		for _, user := range config.Conf.Users {
			if reqUsername == user.Username && reqPassword == user.Password {
				return reqUsername, true
			}
		}
		return reqUsername, false
	}
	return "", false
}

func isHostInWhiteList(host string) bool {
	if host == config.Conf.Fallback.Host {
		return true
	}
	for _, h := range config.Conf.Fallback.WhiteList {
		if h == host {
			return true
		}
	}
	return false
}

func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

func checkACL(req request) (passed bool) {
	portAllowed := !config.Conf.ACL.WhiteListMode
	hostAllowed := !config.Conf.ACL.WhiteListMode
	var ip net.IP
	reqAddr, _ := net.ResolveIPAddr("ip", req.host)
	if reqAddr != nil {
		ip = reqAddr.IP
	}

	for _, rulePort := range config.Conf.ACL.PortsList {
		if rulePort == req.port {
			portAllowed = !portAllowed
			break
		}
	}

	for _, ruleHost := range config.Conf.ACL.HostsList {
		if _, ruleCIDR, err := net.ParseCIDR(ruleHost); err == nil {
			if ruleCIDR.Contains(ip) {
				hostAllowed = !hostAllowed
				break
			}
		}
		if ruleIP := net.ParseIP(ruleHost); ruleIP != nil {
			if ruleIP.Equal(ip) {
				hostAllowed = !hostAllowed
				break
			}
		}
		if strings.Contains(req.host, ruleHost) {
			hostAllowed = !hostAllowed
			break
		}
	}

	if hostAllowed && ip != nil && !config.Conf.ACL.AllowLocal {
		IPv4 := ip.To4()
		if IPv4 != nil && IPv4[0] == 192 && IPv4[1] == 168 {
			hostAllowed = false
		} else if IPv4 != nil && (IPv4[0] == 10 || IPv4[0] == 127) {
			hostAllowed = false
		} else if IPv4 != nil && IPv4[0] == 172 && (IPv4[1] >= 16 && IPv4[1] <= 31) {
			hostAllowed = false
		} else if ip[0] == 0xfd || ip[0] == 0xfe {
			hostAllowed = false
		} else if !ip.IsGlobalUnicast() {
			hostAllowed = false
		}
	}

	return portAllowed && hostAllowed
}
