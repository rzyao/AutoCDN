package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"AutoCDN/cdn"
	"AutoCDN/config"
	"AutoCDN/task"
	"AutoCDN/utils"
)


func init() {
	// 加载配置
	cfg := config.GetConfig()
	var minDelay, maxDelay, downloadTime int
	var maxLossRate float64
	// 命令行参数设置（用于覆盖配置文件的设置）
	flag.IntVar(&task.Routines, "n", cfg.SpeedTest.Routines, "延迟测速线程")
	flag.IntVar(&task.PingTimes, "t", cfg.SpeedTest.PingTimes, "延迟测速次数")
	flag.IntVar(&task.TestCount, "dn", cfg.SpeedTest.TestCount, "下载测速数量")
	flag.IntVar(&downloadTime, "dt", cfg.SpeedTest.DownloadTime, "下载测速时间")
	flag.IntVar(&task.TCPPort, "tp", cfg.SpeedTest.TCPPort, "指定测速端口")
	flag.StringVar(&task.URL, "url", cfg.SpeedTest.SpeedTestURL, "指定测速地址")

	flag.BoolVar(&task.Httping, "httping", cfg.SpeedTest.Httping, "切换测速模式")
	flag.IntVar(&task.HttpingStatusCode, "httping-code", cfg.SpeedTest.HttpingStatusCode, "有效状态代码")
	flag.StringVar(&task.HttpingCFColo, "cfcolo", cfg.SpeedTest.HttpingCFColo, "匹配指定地区")

	flag.IntVar(&maxDelay, "tl", cfg.SpeedTest.MaxDelay, "平均延迟上限")
	flag.IntVar(&minDelay, "tll", cfg.SpeedTest.MinDelay, "平均延迟下限")
	flag.Float64Var(&maxLossRate, "tlr", cfg.SpeedTest.MaxLossRate, "丢包几率上限")
	flag.Float64Var(&task.MinSpeed, "sl", cfg.SpeedTest.MinSpeed, "下载速度下限")

	flag.IntVar(&utils.PrintNum, "p", cfg.SpeedTest.PrintNum, "显示结果数量")
	flag.StringVar(&task.IPFile, "f", cfg.SpeedTest.IPFile, "IP段数据文件")
	flag.StringVar(&task.IPText, "ip", cfg.SpeedTest.IPText, "指定IP段数据")
	flag.StringVar(&utils.Output, "o", cfg.SpeedTest.Output, "输出结果文件")

	flag.BoolVar(&task.Disable, "dd", cfg.SpeedTest.DisableDownload, "禁用下载测速")
	flag.BoolVar(&task.TestAll, "allip", cfg.SpeedTest.TestAllIP, "测速全部 IP")

	var printVersion bool
	flag.BoolVar(&printVersion, "v", false, "打印程序版本")
	
	flag.Parse()

}

func main() {
	cfg := config.GetConfig()
	

	// 开始延迟测速 + 过滤延迟/丢包
	pingData := task.NewPing().Run().FilterDelay().FilterLossRate()
	
	// 开始下载测速
	speedData := task.TestDownloadSpeed(pingData)
	utils.ExportCsv(speedData) // 输出文件
	speedData.Print()          // 打印结果

	// 提取IP列表
	var ips []string
	for _, data := range speedData {
		ips = append(ips, data.PingData.IP.String())
	}

	// 获取与domains数组长度相同的ipList
	ipList, err := cdn.GetIPListForDomains(ips, cfg.Cloudflare.Domains)
	if err != nil {
		log.Fatalf("获取IP列表失败: %v", err)
	}

	// 处理DNS记录
	if err := cdn.HandleDNSRecords(ipList); err != nil {
		log.Fatalf("处理DNS记录失败: %v", err)
	}
	fmt.Println("DNS记录处理完成")

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
