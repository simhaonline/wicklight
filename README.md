# Wicklight

Wicklight is a http(s),http2 server that is written in Golang and is under `MIT License`.

> **Notice: Wicklight is developed only for author's study and testing purposes**

## Features
* A Proxy for HTTP, HTTPS. Unlike other proxies, Wickproxy is an HTTP proxy and no client is needed.
* Highly concealed and Probe-Resistance.
    * Rewrite any illegal requests to a backend server.
    * Wickproxy can work as a frontend server application in the front of  `Caddy` or `Nginx`.
* **HTTP2 Fallback**.
* with or without transport layer security
* One script to install on linux amd64 machiens (see `Install`)
* rule check based on IP, CIDR, domain, port
* Access control list. Allow or deny by IP, ports, domain name, or CIDR. 
* Build for almost all platforms. Wickproxy is compiled for Windows, OS X, Linux, and Freebsd.

### Probe-Resistance
For many HTTP proxy clients, it is common that no authentication information is sent in the first packet, then the server should return a `407 Proxy-Authenticate` to indicate the authentication information should be sent. However, this behavior exposes the fact of wickproxy is a proxy server.

In order to resist probe requests, only requests to `fallback.host`(such as `pr.wickproxy.org` will trick a `407 Proxy-Authenticate` response and other requests will be fallbacked to backend servers. 

However, it is nay not be compatible with some software such as `git` command. A `fallback.whitelist` is introduced to solve this problem. Hosts in `fallback.whitelist` will also trick a `407 Proxy-Authenticate`. `fallback.whitelist` should be used as a workaround and it increases the risk of be detected.

### Fallback Model
It is easy to use Wickproxy as the frontend server listening on port 443 or 80. Any invalid requests will be sent to `fallback.target`. Then, an Nginx or caddy server listen on fallback.


## Install

via curl
```
sudo bash -c "$(curl -fsSL https://github.com/wickproxy/wicklight/blob/main/example/quickstart.sh)"
```
via wget
```
sudo bash -c "$(wget -O- https://github.com/wickproxy/wicklight/blob/main/example/quickstart.sh)"
```
<!--
https://raw.githubusercontent.com/wickproxy/wicktroja/main/example/install.sh
-->

Or download binary manually from: [Release Page](https://github.com/wickproxy/wicklight/releases)

## Usage

Command line usage:
```
wicklight <config.toml> 

```

Please refer to [`example/config.toml`](https://github.com/wickproxy/wicklight/blob/main/example/config.toml) to see how to configure.

## Build
Prerequisites:
* `Golang` 1.12 or above
* `git` to clone this repository

It is easy to build Wicklight using `go` command:
```
git clone https://github.com/wickproxy/wicklight
go build -o build/wicklilght .
```

Another way to compile Wicklight is to use `Make` command:
```
make <platform>       # to build for special platform. Including: linux-amd64, linux-arm64 , darwin-amd64, windows-x64, windows-x86 and freebsd-amd64
make all        # to build for all three platforms
```