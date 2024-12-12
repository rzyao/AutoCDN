# Windows 64位
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o AutoCDN_windows_amd64.exe

# Linux 64位
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o AutoCDN_linux_amd64

# MacOS 64位
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o AutoCDN_darwin_amd64