package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
	"wicklight/config"
	"wicklight/logger"
)

// Run run the server
func Run() {
	var err error
	http.DefaultTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 75 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          200,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	srv := &http.Server{
		Addr:        config.Conf.Listen,
		Handler:     &server{},
		IdleTimeout: 75 * time.Second,
	}

	logger.Info("[server] start to serve at", config.Conf.Listen)
	if config.Conf.TLS.Certificate != "" && config.Conf.TLS.PrivateKey != "" {
		cert, err := tls.LoadX509KeyPair(config.Conf.TLS.Certificate, config.Conf.TLS.PrivateKey)
		if err != nil {
			logger.Fatal("[tls] read certificate error:", err)
		}
		srv.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"h2", "http/1.1"},
			MinVersion:   tls.VersionTLS12,
		}
		err = srv.ListenAndServeTLS("", "")
	} else {
		err = srv.ListenAndServe()
	}
	if err != nil {
		logger.Fatal("[server] server fatal:", err)
	}
}
