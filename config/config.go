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
	Cloudflare CloudflareConfig `yaml:"cloudflare"`
	SpeedTest  SpeedTestConfig  `yaml:"speed_test"`
}

// CloudflareConfig Cloudflare相关配置
type CloudflareConfig struct {
	APIKey   string   `yaml:"api_key"`
	ZoneID   string   `yaml:"zone_id"`
	ZoneName string   `yaml:"zone_name"`
	Email    string   `yaml:"email"`
	Domains  []string `yaml:"domains"`
}

// SpeedTestConfig 速度测试相关配置
type SpeedTestConfig struct {
	// 延迟测速配置
	Routines     int    `yaml:"routines"`      // 延迟测速线程数
	PingTimes    int    `yaml:"ping_times"`    // 延迟测速次数
	TestCount    int    `yaml:"test_count"`    // 下载测速数量
	DownloadTime int    `yaml:"download_time"` // 下载测速时间
	TCPPort      int    `yaml:"tcp_port"`      // 测速端口
	SpeedTestURL string `yaml:"speed_test_url"`// 测速地址

	// HTTP测速配置
	Httping           bool   `yaml:"httping"`             // 是否启用HTTP测速
	HttpingStatusCode int    `yaml:"httping_status_code"` // HTTP状态码
	HttpingCFColo     string `yaml:"httping_cf_colo"`     // 指定地区

	// 延迟和速度限制
	MaxDelay     int     `yaml:"max_delay"`      // 平均延迟上限
	MinDelay     int     `yaml:"min_delay"`      // 平均延迟下限
	MaxLossRate  float64 `yaml:"max_loss_rate"`  // 丢包几率上限
	MinSpeed     float64 `yaml:"min_speed"`      // 下载速度下限

	// 输出配置
	PrintNum    int    `yaml:"print_num"`    // 显示结果数量
	IPFile      string `yaml:"ip_file"`      // IP段数据文件
	IPText      string `yaml:"ip_text"`      // 指定IP段数据
	Output      string `yaml:"output"`       // 输出文件

	// 其他配置
	DisableDownload bool `yaml:"disable_download"` // 禁用下载测速
	TestAllIP       bool `yaml:"test_all_ip"`      // 测试所有IP
}

// LoadConfig 加载配置文件
func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// GetConfig 获取配置单例
func GetConfig() *Config {
	once.Do(func() {
		// 默认配置文件路径
		configPath := "config.yaml"
		
		// 尝试加载配置文件
		cfg, err := LoadConfig(configPath)
		if err != nil {
			// 如果加载失败，使用默认配置
			config = &Config{
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
					IPFile:            "ip.txt",
					IPText:            "",
					Output:            "result.csv",
					DisableDownload:   false,
					TestAllIP:         false,
				},
			}
		} else {
			config = cfg
		}
	})
	return config
}

// ValidateConfig 验证配置有效性
func (c *Config) ValidateConfig() error {
	// TODO: 添加配置验证逻辑
	return nil
}

// SaveConfig 保存配置到文件
func (c *Config) SaveConfig(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
} 