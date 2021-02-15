ifeq ($(OS),Windows_NT)
	builddate = $(shell echo %date:~0,4%%date:~5,2%%date:~8,2%%time:~0,2%%time:~3,2%%time:~6,2%)
else
	builddate = $(shell date +"%Y-%m-%d %H:%M:%S")
endif

Version = 0.0.5
ldflags = -X 'wicklight/version.Version=$(Version)' -X 'wicklight/version.BuildTime=$(builddate)'

build:
	clean
	go build -o build/wicklight .

clean:
	rm -f build/wicklight*
darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64  go build -ldflags "$(ldflags)" -o build/wicklight-darwin-amd64 .
linux-amd64:
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64  go build -ldflags "$(ldflags)" -o build/wicklight-linux-amd64 .
linux-arm64:
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64  go build -ldflags "$(ldflags)" -o build/wicklight-linux-arm64 .
linux-arm:
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm    go build -ldflags "$(ldflags)" -o build/wicklight-linux-arm .
linux-mips64:
	CGO_ENABLED=0 GOOS=linux   GOARCH=mips64 go build -ldflags "$(ldflags)" -o build/wicklight-linux-mips64 .
windows-x64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64  go build -ldflags "$(ldflags)" -o build/wicklight-windows-x64.exe .
windows-x86:
	CGO_ENABLED=0 GOOS=windows GOARCH=386    go build -ldflags "$(ldflags)" -o build/wicklight-windows-x86.exe .
freebsd-amd64:
	CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64  go build -ldflags "$(ldflags)" -o build/wicklight-freebsd-amd64 .

all: clean darwin-amd64 linux-amd64 linux-arm64 linux-arm linux-mips64 windows-x64 windows-x86 freebsd-amd64
