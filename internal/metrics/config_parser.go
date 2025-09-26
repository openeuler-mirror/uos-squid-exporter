// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// SquidConfigData 表示解析后的squid配置数据
type SquidConfigData struct {
	HttpPort        int      `json:"http_port"`
	CacheDir        string   `json:"cache_dir"`
	CoreDumpDir     string   `json:"coredump_dir"`
	LocalNetworks   []string `json:"local_networks"`
	SafePorts       []int    `json:"safe_ports"`
	SSLPorts        []int    `json:"ssl_ports"`
	AccessRules     []string `json:"access_rules"`
	RefreshPatterns []string `json:"refresh_patterns"`
	ACLs            []ACL    `json:"acls"`
}

// ACL 表示访问控制列表项
type ACL struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Value   string `json:"value"`
	Comment string `json:"comment"`
}

// SquidConfigParser squid配置文件解析器
type SquidConfigParser struct {
	filePath string
}

// NewSquidConfigParser 创建新的squid配置解析器
func NewSquidConfigParser(filePath string) *SquidConfigParser {
	return &SquidConfigParser{
		filePath: filePath,
	}
}

// Parse 解析squid配置文件
func (p *SquidConfigParser) Parse() (*SquidConfigData, error) {
	file, err := os.Open(p.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", p.filePath, err)
	}
	defer file.Close()

	config := &SquidConfigData{
		LocalNetworks:   make([]string, 0),
		SafePorts:       make([]int, 0),
		SSLPorts:        make([]int, 0),
		AccessRules:     make([]string, 0),
		RefreshPatterns: make([]string, 0),
		ACLs:            make([]ACL, 0),
	}

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析配置项
		if err := p.parseLine(line, config); err != nil {
			return nil, fmt.Errorf("error parsing line %d: %w", lineNumber, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return config, nil
}

// parseLine 解析单行配置
func (p *SquidConfigParser) parseLine(line string, config *SquidConfigData) error {
	// 解析ACL定义
	if strings.HasPrefix(line, "acl ") {
		return p.parseACL(line, config)
	}

	// 解析http_port
	if strings.HasPrefix(line, "http_port ") {
		return p.parseHttpPort(line, config)
	}

	// 解析cache_dir
	if strings.HasPrefix(line, "cache_dir ") {
		return p.parseCacheDir(line, config)
	}

	// 解析coredump_dir
	if strings.HasPrefix(line, "coredump_dir ") {
		return p.parseCoreDumpDir(line, config)
	}

	// 解析http_access
	if strings.HasPrefix(line, "http_access ") {
		return p.parseHttpAccess(line, config)
	}

	// 解析refresh_pattern
	if strings.HasPrefix(line, "refresh_pattern ") {
		return p.parseRefreshPattern(line, config)
	}

	return nil
}

// parseACL 解析ACL定义
func (p *SquidConfigParser) parseACL(line string, config *SquidConfigData) error {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return fmt.Errorf("invalid ACL format: %s", line)
	}

	acl := ACL{
		Name: parts[1],
		Type: parts[2],
	}

	// 处理值部分
	if len(parts) > 3 {
		valueParts := parts[3:]
		// 移除注释部分
		for i, part := range valueParts {
			if strings.HasPrefix(part, "#") {
				valueParts = valueParts[:i]
				break
			}
		}
		acl.Value = strings.Join(valueParts, " ")
	}

	// 检查是否是本地网络ACL
	if acl.Name == "localnet" && acl.Type == "src" {
		config.LocalNetworks = append(config.LocalNetworks, acl.Value)
	}

	// 检查是否是安全端口ACL
	if acl.Name == "Safe_ports" && acl.Type == "port" {
		if ports, err := p.parsePorts(acl.Value); err == nil {
			config.SafePorts = append(config.SafePorts, ports...)
		}
	}

	// 检查是否是SSL端口ACL
	if acl.Name == "SSL_ports" && acl.Type == "port" {
		if ports, err := p.parsePorts(acl.Value); err == nil {
			config.SSLPorts = append(config.SSLPorts, ports...)
		}
	}

	config.ACLs = append(config.ACLs, acl)
	return nil
}

// parseHttpPort 解析http_port配置
func (p *SquidConfigParser) parseHttpPort(line string, config *SquidConfigData) error {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return fmt.Errorf("invalid http_port format: %s", line)
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid port number: %s", parts[1])
	}

	config.HttpPort = port
	return nil
}
