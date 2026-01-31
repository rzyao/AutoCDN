# Windows 64位
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o AutoCDN_windows_amd64.exe

go build -o ddns.exe -ldflags="-H=windowsgui" d:\project\ddns\main.go

1. 使用 -ldflags 标志：可以设置 -H=windowsgui，这样在 Windows 上运行时不会显示终端窗口。

# Linux 64位
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o AutoCDN_linux_amd64

# MacOS 64位
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o AutoCDN_darwin_amd64