package cdn

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"AutoCDN/config"
)

// GetRecordList 获取指定Zone的DNS记录列表
func GetRecordList(zoneID string) (map[string]string, error) {
	cfg := config.GetConfig()
	log.Printf("DEBUG: GetRecordList called with zoneID: '%s' (Config ZoneID: '%s')", zoneID, cfg.Cloudflare.ZoneID)

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Key", cfg.Cloudflare.APIKey)
	req.Header.Set("X-Auth-Email", cfg.Cloudflare.Email)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取DNS记录失败: %s", resp.Status)
	}

	var result struct {
		Result []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	records := make(map[string]string)
	for _, record := range result.Result {
		// 只处理A记录和AAAA记录
		if record.Type == "A" || record.Type == "AAAA" {
			records[record.Name] = record.ID
		}
	}

	return records, nil
}

// GetRecordListWithType 获取指定Zone的DNS记录列表（包含类型信息）
func GetRecordListWithType(zoneID string) (map[string]map[string]string, error) {
	cfg := config.GetConfig()

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Key", cfg.Cloudflare.APIKey)
	req.Header.Set("X-Auth-Email", cfg.Cloudflare.Email)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取DNS记录失败: %s", resp.Status)
	}

	var result struct {
		Result []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	records := make(map[string]map[string]string)
	for _, record := range result.Result {
		// 只处理A记录和AAAA记录
		if record.Type == "A" || record.Type == "AAAA" {
			if records[record.Name] == nil {
				records[record.Name] = make(map[string]string)
			}
			records[record.Name][record.Type] = record.ID
		}
	}

	return records, nil
}

// UpdateDNSRecords 更新指定域名的DNS解析记录
func UpdateDNSRecords(newIP, domain, zoneID, recordID string) error {
	cfg := config.GetConfig()

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)
	reqBody := fmt.Sprintf(`{
		"type": "A",
		"name": "%s",
		"content": "%s",
		"ttl": 3600,
		"proxied": false
	}`, domain, newIP)

	req, err := http.NewRequest("PATCH", url, strings.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Key", cfg.Cloudflare.APIKey)
	req.Header.Set("X-Auth-Email", cfg.Cloudflare.Email)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("更新DNS记录失败: %s", resp.Status)
	}
	return nil
}

// UpdateDNSRecordsIPv6 更新指定域名的IPv6 DNS解析记录
func UpdateDNSRecordsIPv6(newIP, domain, zoneID, recordID string) error {
	cfg := config.GetConfig()

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)
	reqBody := fmt.Sprintf(`{
		"type": "AAAA",
		"name": "%s",
		"content": "%s",
		"ttl": 3600,
		"proxied": false
	}`, domain, newIP)

	req, err := http.NewRequest("PATCH", url, strings.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Key", cfg.Cloudflare.APIKey)
	req.Header.Set("X-Auth-Email", cfg.Cloudflare.Email)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("更新IPv6 DNS记录失败: %s", resp.Status)
	}
	return nil
}

// CreateDNSRecord 创建新的DNS记录
func CreateDNSRecord(newIP, domain, zoneID string) error {
	cfg := config.GetConfig()

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID)
	reqBody := fmt.Sprintf(`{
		"type": "A",
		"name": "%s",
		"content": "%s",
		"ttl": 3600,
		"proxied": false
	}`, domain, newIP)

	req, err := http.NewRequest("POST", url, strings.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Key", cfg.Cloudflare.APIKey)
	req.Header.Set("X-Auth-Email", cfg.Cloudflare.Email)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("创建DNS记录失败: %s", resp.Status)
	}
	return nil
}

// CreateDNSRecordIPv6 创建新的IPv6 DNS记录
func CreateDNSRecordIPv6(newIP, domain, zoneID string) error {
	cfg := config.GetConfig()

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID)
	reqBody := fmt.Sprintf(`{
		"type": "AAAA",
		"name": "%s",
		"content": "%s",
		"ttl": 3600,
		"proxied": false
	}`, domain, newIP)

	req, err := http.NewRequest("POST", url, strings.NewReader(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Key", cfg.Cloudflare.APIKey)
	req.Header.Set("X-Auth-Email", cfg.Cloudflare.Email)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("创建IPv6 DNS记录失败: %s", resp.Status)
	}
	return nil
}

// HandleDNSRecords 处理DNS记录的更新或创建
func HandleDNSRecords(ipList []string) error {
	cfg := config.GetConfig()

	recordList, err := GetRecordList(cfg.Cloudflare.ZoneID)
	if err != nil {
		return fmt.Errorf("获取记录列表失败: %v", err)
	}

	if len(ipList) == 0 {
		return fmt.Errorf("没有可用的 IPv4 IP")
	}

	// 如果 IP 数量少于域名数量，循环复用 IP
	log.Printf("更新 %d 个域名 (可用 IP: %d)", len(cfg.Cloudflare.Domains), len(ipList))

	for i, domain := range cfg.Cloudflare.Domains {
		// 使用取模运算循环分配 IP
		newIP := ipList[i%len(ipList)]

		if recordID, exists := recordList[domain]; exists {
			// 更新记录
			if err := UpdateDNSRecords(newIP, domain, cfg.Cloudflare.ZoneID, recordID); err != nil {
				log.Printf("更新记录失败 %s: %v", domain, err)
				continue
			}
			log.Printf("成功更新记录: %s -> %s", domain, newIP)
		} else {
			// 创建记录
			if err := CreateDNSRecord(newIP, domain, cfg.Cloudflare.ZoneID); err != nil {
				log.Printf("创建记录失败 %s: %v", domain, err)
				continue
			}
			log.Printf("成功创建记录: %s -> %s", domain, newIP)
		}
	}
	return nil
}

// HandleDNSRecordsIPv6 处理IPv6 DNS记录的更新或创建
func HandleDNSRecordsIPv6(ipList []string) error {
	cfg := config.GetConfig()

	recordList, err := GetRecordList(cfg.Cloudflare.ZoneID)
	if err != nil {
		return fmt.Errorf("获取记录列表失败: %v", err)
	}

	if len(ipList) == 0 {
		return fmt.Errorf("没有可用的 IPv6 IP")
	}

	// 如果 IP 数量少于域名数量，循环复用 IP
	log.Printf("更新 %d 个 IPv6 域名 (可用 IP: %d)", len(cfg.Cloudflare.DomainIPv6s), len(ipList))

	for i, domain := range cfg.Cloudflare.DomainIPv6s {
		// 使用取模运算循环分配 IP
		newIP := ipList[i%len(ipList)]

		if recordID, exists := recordList[domain]; exists {
			// 更新记录
			if err := UpdateDNSRecordsIPv6(newIP, domain, cfg.Cloudflare.ZoneID, recordID); err != nil {
				log.Printf("更新IPv6记录失败 %s: %v", domain, err)
				continue
			}
			log.Printf("成功更新IPv6记录: %s -> %s", domain, newIP)
		} else {
			// 创建记录
			if err := CreateDNSRecordIPv6(newIP, domain, cfg.Cloudflare.ZoneID); err != nil {
				log.Printf("创建IPv6记录失败 %s: %v", domain, err)
				continue
			}
			log.Printf("成功创建IPv6记录: %s -> %s", domain, newIP)
		}
	}
	return nil
}

// HandleAllDNSRecords 同时处理IPv4和IPv6 DNS记录的更新或创建
func HandleAllDNSRecords(ipv4List []string, ipv6List []string) error {
	cfg := config.GetConfig()

	// 处理IPv4域名
	if len(cfg.Cloudflare.Domains) > 0 {
		if err := HandleDNSRecords(ipv4List); err != nil {
			log.Printf("处理IPv4 DNS记录失败: %v", err)
		}
	}

	// 处理IPv6域名
	if len(cfg.Cloudflare.DomainIPv6s) > 0 {
		if err := HandleDNSRecordsIPv6(ipv6List); err != nil {
			log.Printf("处理IPv6 DNS记录失败: %v", err)
		}
	}

	return nil
}

// GetIPListForDomains 返回所有可用的 IP，不再截断或报错
func GetIPListForDomains(speedData []string, domains []string) ([]string, error) {
	if len(speedData) == 0 {
		return nil, fmt.Errorf("没有可用的测速数据")
	}
	// 不再限制长度，返回所有可用IP供上层循环分配
	return speedData, nil
}

// GetIPListForIPv6Domains 返回所有可用的 IPv6 IP，不再截断或报错
func GetIPListForIPv6Domains(speedData []string, domains []string) ([]string, error) {
	if len(speedData) == 0 {
		return nil, fmt.Errorf("没有可用的测速数据")
	}
	// 不再限制长度，返回所有可用IP供上层循环分配
	return speedData, nil
}
