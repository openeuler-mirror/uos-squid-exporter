// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

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
