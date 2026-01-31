package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"AutoCDN/cdn"
	"AutoCDN/config"
	"AutoCDN/task"
	"AutoCDN/utils"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called when the app terminates
func (a *App) shutdown(ctx context.Context) {
	// 强制取消所有正在进行的任务
	utils.SetCancel()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// ConfigListResult return for config list
type ConfigListResult struct {
	Active  string   `json:"active"`
	Configs []string `json:"configs"`
}

// GetConfigList 获取所有配置文件
func (a *App) GetConfigList() ([]string, error) {
	matches, err := filepath.Glob("*.yaml")
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// LoadConfig 加载指定配置
func (a *App) LoadConfig(name string) (*config.Config, error) {
	return config.LoadConfig(name)
}

// SaveConfig 保存配置
func (a *App) SaveConfig(name string, cfg config.Config) error {
	runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Saving config [%s]...", name))
	// Debug print
	fmt.Printf("[DEBUG] SaveConfig received: %+v\n", cfg)
	return cfg.SaveConfig(name)
}

// CreateNewConfig 创建新配置（使用默认值）
func (a *App) CreateNewConfig(name string) error {
	cfg := config.NewDefaultConfig()
	return cfg.SaveConfig(name)
}

// DeleteConfig 删除配置
func (a *App) DeleteConfig(name string) error {
	runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Deleting config [%s]...", name))
	return os.Remove(name)
}

// StartSpeedTest 启动测速 (Sync or Async handled by Wails? Wails calls are async from JS, but blocking in Go)
// We should run this in a goroutine and emit events?
// If we block config.go/task.go logic runs in current goroutine.
// Better to return and run in background, but the user might want "Wait".
// Since we have events, we can block here (Wails runs in a separate goroutine usually for method calls).
func (a *App) StartSpeedTest(configName string, mode string) error {
	// 1. Load Config
	cfg, err := config.LoadConfig(configName)
	if err != nil {
		return fmt.Errorf("load config failed: %w", err)
	}

	// 2. Set Global Vars in config package (Since task package uses globals initialized from config.GetConfig())
	// Wait, config.GetConfig() is a singleton.
	// We might need to manually inject values into task package globals?
	// task.Routines = cfg.SpeedTest.Routines
	// task/tcping.go says `var Routines = defaultRoutines`.
	// main.go did: `flag.IntVar(&task.Routines...`
	// So we must manually set task package globals.

	// Init Cancel State
	utils.ResetCancel()

	task.Routines = cfg.SpeedTest.Routines
	task.PingTimes = cfg.SpeedTest.PingTimes
	task.TestCount = cfg.SpeedTest.TestCount
	task.TCPPort = cfg.SpeedTest.TCPPort
	task.URL = cfg.SpeedTest.SpeedTestURL
	task.Httping = cfg.SpeedTest.Httping
	task.HttpingStatusCode = cfg.SpeedTest.HttpingStatusCode
	task.HttpingCFColo = cfg.SpeedTest.HttpingCFColo
	task.MinSpeed = cfg.SpeedTest.MinSpeed
	task.MinSpeed = cfg.SpeedTest.MinSpeed
	// utils vars
	utils.InputMaxDelay = time.Duration(cfg.SpeedTest.MaxDelay) * time.Millisecond
	utils.InputMinDelay = time.Duration(cfg.SpeedTest.MinDelay) * time.Millisecond
	utils.InputMaxLossRate = float32(cfg.SpeedTest.MaxLossRate)

	runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Delay Filter: %d ~ %d ms, MaxLoss: %.2f", cfg.SpeedTest.MinDelay, cfg.SpeedTest.MaxDelay, cfg.SpeedTest.MaxLossRate))

	// 3. Set Progress Handler
	utils.ProgressHandler = func(current, total int, msg string) {
		runtime.EventsEmit(a.ctx, "progress", map[string]interface{}{
			"current": current,
			"total":   total,
			"msg":     msg,
		})
	}

	if a.ctx == nil {
		return fmt.Errorf("context is nil")
	}

	runtime.EventsEmit(a.ctx, "log", "Start SpeedTest...")

	// 4. Run Ping
	// We need to capture stdout logs too if possible?
	// For now we rely on progress events.

	// Replaced go func with synchronous execution
	func() {
		defer func() {
			if r := recover(); r != nil {
				runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("Panic: %v", r))
			}
		}()

		// Determine Mode and IP File
		testType := strings.ToUpper(cfg.SpeedTest.TestType)
		if testType == "" {
			testType = "IPV4"
		}

		var currentIPFile string
		if testType == "IPV6" {
			currentIPFile = cfg.SpeedTest.IPv6File
			if currentIPFile == "" {
				currentIPFile = "ipv6.txt"
			}
			task.IPFile = currentIPFile
		} else {
			currentIPFile = cfg.SpeedTest.IPv4File
			if currentIPFile == "" {
				currentIPFile = "ip.txt"
			}
			task.IPFile = currentIPFile
		}

		runtime.EventsEmit(a.ctx, "log", fmt.Sprintf("Mode: %s, File: %s", testType, currentIPFile))

		// Check if IP file exists
		if _, err := os.Stat(currentIPFile); os.IsNotExist(err) {
			runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("IP File not found: %s", currentIPFile))
			return
		}

		// Ping
		runtime.EventsEmit(a.ctx, "status", fmt.Sprintf("Starting Ping (%s)...", testType))
		pingData := task.NewPingWithFile(currentIPFile).Run().FilterDelay().FilterLossRate()

		if len(pingData) == 0 {
			runtime.EventsEmit(a.ctx, "status", "Ping found 0 valid IPs. Stopping.")
			runtime.EventsEmit(a.ctx, "error", "延迟测速结果为 0，请检查：1. IP文件内容 2. 网络连接 3. 筛选条件(如最大延迟/丢包率)")
			return
		}

		runtime.EventsEmit(a.ctx, "status", "Starting Download Test...")
		speedData := task.TestDownloadSpeed(pingData)

		runtime.EventsEmit(a.ctx, "status", "Test Finished.")

		// Emit results
		runtime.EventsEmit(a.ctx, "result", speedData)

		// Filter and Update DNS if Auto Mode
		if mode == "auto" {
			runtime.EventsEmit(a.ctx, "status", "Updating DNS...")

			// Extract IPs
			var ips []string
			for _, data := range speedData {
				ips = append(ips, data.PingData.IP.String())
			}

			if testType == "IPV6" {
				// Handle IPv6
				if len(cfg.Cloudflare.DomainIPv6s) == 0 {
					runtime.EventsEmit(a.ctx, "error", "IPv6 Mode but no IPv6 domains configured!")
					return
				}

				ipList, err := cdn.GetIPListForIPv6Domains(ips, cfg.Cloudflare.DomainIPv6s)
				if err != nil {
					runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("Get IPv6 List failed: %v", err))
					return
				}

				if err := cdn.HandleDNSRecordsIPv6(ipList); err != nil {
					runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("Update IPv6 DNS failed: %v", err))
					return
				}
				runtime.EventsEmit(a.ctx, "status", "IPv6 DNS Updated Successfully!")

			} else {
				// Handle IPv4
				if len(cfg.Cloudflare.Domains) == 0 {
					runtime.EventsEmit(a.ctx, "error", "IPv4 Mode but no IPv4 domains configured!")
					return
				}

				ipList, err := cdn.GetIPListForDomains(ips, cfg.Cloudflare.Domains)
				if err != nil {
					runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("Get IP List failed: %v", err))
					return
				}

				if err := cdn.HandleDNSRecords(ipList); err != nil {
					runtime.EventsEmit(a.ctx, "error", fmt.Sprintf("Update DNS failed: %v", err))
					return
				}
				runtime.EventsEmit(a.ctx, "status", "DNS Updated Successfully!")
			}
		}
	}()

	return nil
}

// detectIPTypeFromFile helper
func detectIPTypeFromFile(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.Contains(line, ":") {
			return true, nil // IPv6
		}
		if strings.Contains(line, ".") {
			return false, nil // IPv4
		}
	}
	return false, nil
}

// StopSpeedTest 停止测速
func (a *App) StopSpeedTest() {
	utils.SetCancel()
	runtime.EventsEmit(a.ctx, "status", "正在停止任务...")
	runtime.EventsEmit(a.ctx, "log", "[CONTROL] 接收到停止指令")
}
