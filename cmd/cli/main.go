package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"AutoCDN/cdn"
	"AutoCDN/config"
	"AutoCDN/task"
	"AutoCDN/utils"
)

// detectIPTypeFromFile 从文件中检测IP类型
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

		// 检查是否包含IPv6特征（冒号）
		if strings.Contains(line, ":") {
			return true, nil // IPv6
		}

		// 检查是否包含IPv4特征（点）
		if strings.Contains(line, ".") {
			return false, nil // IPv4
		}
	}

	// 默认返回IPv4
	return false, nil
}

func main() {
	var configPath string
	var printVersion bool

	// 定义命令行参数
	flag.StringVar(&configPath, "c", "config.yaml", "配置文件路径")
	flag.BoolVar(&printVersion, "v", false, "打印程序版本")

	// 其他参数 (使用零值作为默认值，稍后应用配置)
	flag.IntVar(&task.Routines, "n", 0, "延迟测速线程 (默认使用配置文件)")
	flag.IntVar(&task.PingTimes, "t", 0, "延迟测速次数")
	flag.IntVar(&task.TestCount, "dn", 0, "下载测速数量")
	flag.IntVar(&task.TCPPort, "tp", 0, "指定测速端口")
	flag.StringVar(&task.URL, "url", "", "指定测速地址")
	flag.BoolVar(&task.Httping, "httping", false, "切换测速模式")
	flag.IntVar(&task.HttpingStatusCode, "httping-code", 0, "有效状态代码")
	flag.StringVar(&task.HttpingCFColo, "cfcolo", "", "匹配指定地区")

	var maxDelay, minDelay, downloadTime int
	var maxLossRate float64

	flag.IntVar(&maxDelay, "tl", 0, "平均延迟上限")
	flag.IntVar(&minDelay, "tll", 0, "平均延迟下限")
	flag.IntVar(&downloadTime, "dt", 0, "下载测速时间")
	flag.Float64Var(&maxLossRate, "tlr", 0, "丢包几率上限")
	flag.Float64Var(&task.MinSpeed, "sl", 0, "下载速度下限")

	flag.IntVar(&utils.PrintNum, "p", 0, "显示结果数量")
	flag.StringVar(&task.IPFile, "f", "", "IP段数据文件")
	flag.StringVar(&task.IPText, "ip", "", "指定IP段数据")
	flag.StringVar(&utils.Output, "o", "", "输出结果文件")
	flag.BoolVar(&task.Disable, "dd", false, "禁用下载测速")
	flag.BoolVar(&task.TestAll, "allip", false, "测速全部 IP")

	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		if configPath == "config.yaml" && os.IsNotExist(err) {
			fmt.Println("配置文件不存在，使用默认配置")
			cfg = config.NewDefaultConfig()
		} else {
			log.Fatalf("无法加载配置文件 %s: %v", configPath, err)
		}
	}

	// 更新全局配置单例，确保其他包(如 cdn)调用 config.GetConfig() 时能获取到正确的配置
	config.SetConfig(cfg)

	fmt.Printf("DEBUG: Loaded Config - ZoneID: %s, Domains: %v\n", cfg.Cloudflare.ZoneID, cfg.Cloudflare.Domains)

	// 提前验证 API 连通性，避免跑完测速才发现 API 不通
	// 这里通过尝试获取记录列表来验证，如果是 404/401 等错误直接打印并退出（或者警告）
	if _, err := cdn.GetRecordList(cfg.Cloudflare.ZoneID); err != nil {
		log.Printf("警告: API 连接测试失败 (ZoneID: %s): %v", cfg.Cloudflare.ZoneID, err)
		log.Printf("请检查配置文件中的 ZoneID 和 API Key 是否正确，或是否有权限访问该 Zone。")
	} else {
		fmt.Println("DEBUG: API 连接测试成功 (ZoneID 有效)")
	}

	// 将 CLI 参数覆盖到配置中 (如果 CLI 参数不为零值)
	if task.Routines != 0 {
		cfg.SpeedTest.Routines = task.Routines
	}
	if task.PingTimes != 0 {
		cfg.SpeedTest.PingTimes = task.PingTimes
	}
	if task.TestCount != 0 {
		cfg.SpeedTest.TestCount = task.TestCount
	}
	if task.TCPPort != 0 {
		cfg.SpeedTest.TCPPort = task.TCPPort
	}
	if task.URL != "" {
		cfg.SpeedTest.SpeedTestURL = task.URL
	}
	// Httping bool is tricky to override if default is true and user wants false without explicit flag interaction,
	// but here we assume CLI flag presence implies override if true.
	// For bool flags, standard `flag` pkg sets true if present.
	// If config is true and user executes without flag, it stays true.
	// If config is false and user adds -httping, it becomes true.
	if task.Httping {
		cfg.SpeedTest.Httping = true
	}

	if maxDelay != 0 {
		cfg.SpeedTest.MaxDelay = maxDelay
	}
	if minDelay != 0 {
		cfg.SpeedTest.MinDelay = minDelay
	}
	if downloadTime != 0 {
		cfg.SpeedTest.DownloadTime = downloadTime
	}
	if maxLossRate != 0 {
		cfg.SpeedTest.MaxLossRate = maxLossRate
	}
	if task.MinSpeed != 0 {
		cfg.SpeedTest.MinSpeed = task.MinSpeed
	}
	if utils.PrintNum != 0 {
		cfg.SpeedTest.PrintNum = utils.PrintNum
	}
	if task.IPFile != "" {
		cfg.SpeedTest.IPv4File = task.IPFile
		cfg.SpeedTest.IPv6File = task.IPFile
	} // Simple override
	if utils.Output != "" {
		cfg.SpeedTest.Output = utils.Output
	}

	// 应用最终配置到 Global Vars
	task.Routines = cfg.SpeedTest.Routines
	task.PingTimes = cfg.SpeedTest.PingTimes
	task.TestCount = cfg.SpeedTest.TestCount
	task.Timeout = time.Duration(cfg.SpeedTest.DownloadTime) * time.Second
	task.TCPPort = cfg.SpeedTest.TCPPort
	task.URL = cfg.SpeedTest.SpeedTestURL
	task.Httping = cfg.SpeedTest.Httping
	task.HttpingStatusCode = cfg.SpeedTest.HttpingStatusCode
	task.HttpingCFColo = cfg.SpeedTest.HttpingCFColo
	task.MinSpeed = cfg.SpeedTest.MinSpeed

	utils.InputMaxDelay = time.Duration(cfg.SpeedTest.MaxDelay) * time.Millisecond
	utils.InputMinDelay = time.Duration(cfg.SpeedTest.MinDelay) * time.Millisecond
	utils.InputMaxLossRate = float32(cfg.SpeedTest.MaxLossRate)
	utils.PrintNum = cfg.SpeedTest.PrintNum
	utils.Output = cfg.SpeedTest.Output

	// Determine correct IP File based on TestType (if not overridden by -f)
	if task.IPFile == "" {
		// CLI mode default to what config says, or IPv4 if simplified
		if cfg.SpeedTest.TestType == "IPV6" {
			task.IPFile = cfg.SpeedTest.IPv6File
		} else {
			task.IPFile = cfg.SpeedTest.IPv4File
		}
	}

	// 检查IP文件是否存在
	if task.IPFile == "" {
		log.Fatal("未指定IP文件，请在配置文件中设置 ip_file 或使用 -f 参数")
	}

	// 检测IP类型
	isIPv6, err := detectIPTypeFromFile(task.IPFile)
	if err != nil {
		log.Fatalf("检测IP类型失败: %v", err)
	}

	if isIPv6 {
		// 处理IPv6
		if len(cfg.Cloudflare.DomainIPv6s) == 0 {
			log.Fatal("检测到IPv6文件，但未配置IPv6域名，请在配置文件中设置 domainipv6s")
		}

		fmt.Printf("开始处理IPv6域名 (使用配置: %s)...\n", configPath)
		pingData := task.NewPingWithFile(task.IPFile).Run().FilterDelay().FilterLossRate()
		speedData := task.TestDownloadSpeed(pingData)
		utils.ExportCsvToFile(speedData, cfg.SpeedTest.Output)
		speedData.Print() // 打印结果

		// 提取IPv6 IP列表
		var ips []string
		for _, data := range speedData {
			ips = append(ips, data.PingData.IP.String())
		}

		// 获取与IPv6 domains数组长度相同的ipList
		ipList, err := cdn.GetIPListForIPv6Domains(ips, cfg.Cloudflare.DomainIPv6s)
		if err != nil {
			log.Printf("获取IPv6 IP列表失败: %v", err)
		} else {
			// 处理IPv6 DNS记录
			if err := cdn.HandleDNSRecordsIPv6(ipList); err != nil {
				log.Printf("处理IPv6 DNS记录失败: %v", err)
			} else {
				fmt.Println("IPv6 DNS记录处理完成")
			}
		}
	} else {
		// 处理IPv4
		if len(cfg.Cloudflare.Domains) == 0 {
			log.Fatal("检测到IPv4文件，但未配置IPv4域名，请在配置文件中设置 domains")
		}

		fmt.Printf("开始处理IPv4域名 (使用配置: %s)...\n", configPath)
		pingData := task.NewPingWithFile(task.IPFile).Run().FilterDelay().FilterLossRate()
		speedData := task.TestDownloadSpeed(pingData)
		utils.ExportCsvToFile(speedData, cfg.SpeedTest.Output)
		speedData.Print() // 打印结果

		// 提取IP列表
		var ips []string
		for _, data := range speedData {
			ips = append(ips, data.PingData.IP.String())
		}

		// 获取与domains数组长度相同的ipList
		ipList, err := cdn.GetIPListForDomains(ips, cfg.Cloudflare.Domains)
		if err != nil {
			log.Printf("获取IPv4 IP列表失败: %v", err)
		} else {
			// 处理IPv4 DNS记录
			if err := cdn.HandleDNSRecords(ipList); err != nil {
				log.Printf("处理IPv4 DNS记录失败: %v", err)
			} else {
				fmt.Println("IPv4 DNS记录处理完成")
			}
		}
	}

	endPrint()
}

func endPrint() {
	if utils.NoPrintResult() {
		return
	}
	if runtime.GOOS == "windows" {
		fmt.Printf("按下 回车键 或 Ctrl+C 退出。")
		fmt.Scanln()
	}
}
