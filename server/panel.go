package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
	"wicklight/config"
	"wicklight/logger"
	"wicklight/version"
)

const responseBody = `<html>
<head><title>%v</title></head>
<body>
<center><h1>Wicklight Panel</h1></center>
<center><p>%v</p></center>
<hr><center>Wicklight/%v</center>
</body>
</html>
`

func handlePanel(w http.ResponseWriter, r *http.Request, req request) {

	if config.Conf.Fallback.PACPath != "" && config.Conf.Fallback.PACPath == r.URL.Path {
		handlePAC(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Server", "Wicklight/"+version.Version)
	w.Header().Set("Date", time.Now().Format(time.RFC1123))

	code := http.StatusOK
	msg := ""
	if req.user == "" && !req.authenticated {
		w.Header().Set("Proxy-Authenticate", "Basic realm=\"Wicklight\"")
		code = http.StatusProxyAuthRequired
		msg = "need to login in"
	} else if !req.authenticated {
		code = http.StatusUnauthorized
		msg = "username or password is not correct"
	} else if !req.allowed {
		code = http.StatusForbidden
		msg = "not allowed to visit"
	} else if req.err != nil && req.err == errGateway {
		code = http.StatusBadGateway
		msg = "can not connect to the next hop"
	} else if req.err != nil {
		logger.Debug("[error]", req.err)
		code = http.StatusInternalServerError
		msg = req.err.Error()
	} else {
		msg = fmt.Sprintf("User [%v], welcome to wicklight", req.user)
	}

	text := http.StatusText(code)
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf(responseBody, "Wicklight Panel-"+text, text+": "+msg, version.Version)))
}

func handleReverseProxy(target string, w http.ResponseWriter, r *http.Request, hide bool) {
	targetURL, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	if !hide {
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.Header.Set("X-Real-IP", r.RemoteAddr)
	}
	r.URL.Host = targetURL.Host
	r.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	proxy.ServeHTTP(w, r)
}

func handlePAC(w http.ResponseWriter, r *http.Request) {
	if exists(config.Conf.Fallback.PACFile) {
		http.ServeFile(w, r, config.Conf.Fallback.PACFile)
		return
	}

	const pacFile = `
function FindProxyForURL(url, host) {
	if (host === "127.0.0.1" || host === "::1" || host === "localhost")
		return "DIRECT";
	return "HTTPS %s:%s";
}
`
	fmt.Fprintf(w, pacFile, config.Conf.Fallback.PACHost, config.Conf.Fallback.PACPort)
	return
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
