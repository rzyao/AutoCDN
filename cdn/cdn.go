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
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	records := make(map[string]string)
	for _, record := range result.Result {
		records[record.Name] = record.ID
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

// HandleDNSRecords 处理DNS记录的更新或创建
func HandleDNSRecords(ipList []string) error {
	cfg := config.GetConfig()
	
	recordList, err := GetRecordList(cfg.Cloudflare.ZoneID)
	if err != nil {
		return fmt.Errorf("获取记录列表失败: %v", err)
	}

	if len(ipList) < len(cfg.Cloudflare.Domains) {
		return fmt.Errorf("IP列表长度(%d)小于域名数量(%d)", 
			len(ipList), len(cfg.Cloudflare.Domains))
	}

	for i, domain := range cfg.Cloudflare.Domains {
		newIP := ipList[i]
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

// GetIPListForDomains 根据domains的长度截取speedData
func GetIPListForDomains(speedData []string, domains []string) ([]string, error) {
	if len(speedData) < len(domains) {
		return nil, fmt.Errorf("speedData长度(%d)小于domains长度(%d)", 
			len(speedData), len(domains))
	}

	return speedData[:len(domains)], nil
}