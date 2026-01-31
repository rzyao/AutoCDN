package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// 使用单例模式存储配置
var (
	config *Config
	once   sync.Once
)

// Config 总配置结构
type Config struct {
	Cloudflare CloudflareConfig `yaml:"cloudflare" json:"Cloudflare"`
	SpeedTest  SpeedTestConfig  `yaml:"speed_test" json:"SpeedTest"`
}

// CloudflareConfig Cloudflare相关配置
type CloudflareConfig struct {
	APIKey      string   `yaml:"api_key" json:"APIKey"`
	ZoneID      string   `yaml:"zone_id" json:"ZoneID"`
	ZoneName    string   `yaml:"zone_name" json:"ZoneName"`
	Email       string   `yaml:"email" json:"Email"`
	Domains     []string `yaml:"domains" json:"Domains"`
	DomainIPv6s []string `yaml:"domainipv6s" json:"DomainIPv6s"`
}

// SpeedTestConfig 速度测试相关配置
type SpeedTestConfig struct {
	// 延迟测速配置
	Routines     int    `yaml:"routines" json:"Routines"`           // 延迟测速线程数
	PingTimes    int    `yaml:"ping_times" json:"PingTimes"`        // 延迟测速次数
	TestCount    int    `yaml:"test_count" json:"TestCount"`        // 下载测速数量
	DownloadTime int    `yaml:"download_time" json:"DownloadTime"`  // 下载测速时间
	TCPPort      int    `yaml:"tcp_port" json:"TCPPort"`            // 测速端口
	SpeedTestURL string `yaml:"speed_test_url" json:"SpeedTestURL"` // 测速地址

	// HTTP测速配置
	Httping           bool   `yaml:"httping" json:"Httping"`                       // 是否启用HTTP测速
	HttpingStatusCode int    `yaml:"httping_status_code" json:"HttpingStatusCode"` // HTTP状态码
	HttpingCFColo     string `yaml:"httping_cf_colo" json:"HttpingCFColo"`         // 指定地区

	// 延迟和速度限制
	MaxDelay    int     `yaml:"max_delay" json:"MaxDelay"`        // 平均延迟上限
	MinDelay    int     `yaml:"min_delay" json:"MinDelay"`        // 平均延迟下限
	MaxLossRate float64 `yaml:"max_loss_rate" json:"MaxLossRate"` // 丢包几率上限
	MinSpeed    float64 `yaml:"min_speed" json:"MinSpeed"`        // 下载速度下限

	// 输出配置
	PrintNum int    `yaml:"print_num" json:"PrintNum"` // 显示结果数量
	IPv4File string `yaml:"ipv4_file" json:"IPv4File"` // IPv4数据文件
	IPv6File string `yaml:"ipv6_file" json:"IPv6File"` // IPv6数据文件
	TestType string `yaml:"test_type" json:"TestType"` // 测速类型 IPV4/IPV6
	Output   string `yaml:"output" json:"Output"`      // 输出文件

	// 其他配置
	DisableDownload bool `yaml:"disable_download" json:"DisableDownload"` // 禁用下载测速
	TestAllIP       bool `yaml:"test_all_ip" json:"TestAllIP"`            // 测试所有IP
}

// NewDefaultConfig 返回默认配置
func NewDefaultConfig() *Config {
	return &Config{
		SpeedTest: SpeedTestConfig{
			Routines:          200,
			PingTimes:         4,
			TestCount:         10,
			DownloadTime:      10,
			TCPPort:           443,
			SpeedTestURL:      "https://speedtest.ayaoblog.space/file.mp4",
			Httping:           false,
			HttpingStatusCode: 200,
			HttpingCFColo:     "",
			MaxDelay:          9999,
			MinDelay:          0,
			MaxLossRate:       1,
			MinSpeed:          0,
			PrintNum:          10,
			IPv4File:          "ip.txt",
			IPv6File:          "ipv6.txt",
			TestType:          "IPV4",
			Output:            "result.csv",
			DisableDownload:   false,
			TestAllIP:         false,
		},
	}
}

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := NewDefaultConfig()
	if err := yaml.Unmarshal(file, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetConfig 获取配置单例
func GetConfig() *Config {
	once.Do(func() {
		// 如果已经手动设置了配置（例如通过 SetConfig），则跳过默认加载逻辑
		if config != nil {
			return
		}

		// 默认配置文件路径
		configPath := "config.yaml"

		// 尝试加载配置文件
		cfg, err := LoadConfig(configPath)
		if err != nil {
			// 如果加载失败，使用默认配置
			config = NewDefaultConfig()
		} else {
			config = cfg
			// Set defaults if empty
			if config.SpeedTest.IPv4File == "" {
				config.SpeedTest.IPv4File = "ip.txt"
			}
			if config.SpeedTest.IPv6File == "" {
				config.SpeedTest.IPv6File = "ipv6.txt"
			}
			if config.SpeedTest.TestType == "" {
				config.SpeedTest.TestType = "IPV4"
			}
		}
	})
	return config
}

// ValidateConfig 验证配置有效性
func (c *Config) ValidateConfig() error {
	// TODO: 添加配置验证逻辑
	return nil
}

// SetConfig 更新全局配置单例
func SetConfig(cfg *Config) {
	config = cfg
}

// SaveConfig 保存配置到文件
func (c *Config) SaveConfig(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
