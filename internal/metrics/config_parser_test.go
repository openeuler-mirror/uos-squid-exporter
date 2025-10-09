// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
	"os"
	"testing"
)

func TestSquidConfigParser_Parse(t *testing.T) {
	// 创建临时测试配置文件
	testConfig := `#
# Recommended minimum configuration:
#

# Example rule allowing access from your local networks.
acl localnet src 10.0.0.0/8     # RFC1918 possible internal network
acl localnet src 172.16.0.0/12  # RFC1918 possible internal network
acl localnet src 192.168.0.0/16 # RFC1918 possible internal network

acl SSL_ports port 443
acl Safe_ports port 80          # http
acl Safe_ports port 21          # ftp
acl Safe_ports port 443         # https
acl Safe_ports port 1025-65535  # unregistered ports

#
# Recommended minimum Access Permission configuration:
#
http_access deny !Safe_ports
http_access deny CONNECT !SSL_ports
http_access allow localhost manager
http_access deny manager
http_access allow localnet
http_access allow localhost
http_access deny all

# Squid normally listens to port 3128
http_port 3128

# Uncomment and adjust the following to add a disk cache directory.
cache_dir ufs /var/spool/squid 100 16 256

# Leave coredumps in the first cache dir
coredump_dir /var/spool/squid

#
# Add any of your own refresh_pattern entries above these.
#
refresh_pattern ^ftp:           1440    20%     10080
refresh_pattern ^gopher:        1440    0%      1440
refresh_pattern .               0       20%     4320
`

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "squid_test.conf")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入测试配置
	if _, err := tmpFile.WriteString(testConfig); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	tmpFile.Close()

	// 测试解析器
	parser := NewSquidConfigParser(tmpFile.Name())
	config, err := parser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	// 验证解析结果
	if config.HttpPort != 3128 {
		t.Errorf("Expected http_port 3128, got %d", config.HttpPort)
	}

	if config.CacheDir != "ufs /var/spool/squid 100 16 256" {
		t.Errorf("Expected cache_dir 'ufs /var/spool/squid 100 16 256', got '%s'", config.CacheDir)
	}

	if config.CoreDumpDir != "/var/spool/squid" {
		t.Errorf("Expected coredump_dir '/var/spool/squid', got '%s'", config.CoreDumpDir)
	}

	// 验证本地网络
	expectedNetworks := []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}
	if len(config.LocalNetworks) != len(expectedNetworks) {
		t.Errorf("Expected %d local networks, got %d", len(expectedNetworks), len(config.LocalNetworks))
	}

	// 验证安全端口
	if len(config.SafePorts) < 4 {
		t.Errorf("Expected at least 4 safe ports, got %d", len(config.SafePorts))
	}

	// 验证SSL端口
	if len(config.SSLPorts) != 1 || config.SSLPorts[0] != 443 {
		t.Errorf("Expected SSL port 443, got %v", config.SSLPorts)
	}

	// 验证访问规则
	if len(config.AccessRules) < 5 {
		t.Errorf("Expected at least 5 access rules, got %d", len(config.AccessRules))
	}

	// 验证刷新模式
	if len(config.RefreshPatterns) != 3 {
		t.Errorf("Expected 3 refresh patterns, got %d", len(config.RefreshPatterns))
	}

	// 验证ACL
	if len(config.ACLs) < 5 {
		t.Errorf("Expected at least 5 ACLs, got %d", len(config.ACLs))
	}
}

func TestSquidConfigParser_ParsePorts(t *testing.T) {
	parser := &SquidConfigParser{}

	// 测试单个端口
	ports, err := parser.parsePorts("80")
	if err != nil {
		t.Fatalf("Failed to parse single port: %v", err)
	}
	if len(ports) != 1 || ports[0] != 80 {
		t.Errorf("Expected port 80, got %v", ports)
	}

	// 测试端口范围
	ports, err = parser.parsePorts("1025-1027")
	if err != nil {
		t.Fatalf("Failed to parse port range: %v", err)
	}
	expected := []int{1025, 1026, 1027}
	if len(ports) != len(expected) {
		t.Errorf("Expected ports %v, got %v", expected, ports)
	}
}

func TestSquidConfigData_Validate(t *testing.T) {
	config := &SquidConfigData{
		HttpPort:      3128,
		LocalNetworks: []string{"192.168.0.0/16"},
		SafePorts:     []int{80, 443},
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Valid config should not fail validation: %v", err)
	}

	// 测试无效端口
	config.HttpPort = 0
	err = config.Validate()
	if err == nil {
		t.Error("Invalid port should fail validation")
	}

	// 测试无本地网络
	config.HttpPort = 3128
	config.LocalNetworks = []string{}
	err = config.Validate()
	if err == nil {
		t.Error("Empty local networks should fail validation")
	}

	// 测试无安全端口
	config.LocalNetworks = []string{"192.168.0.0/16"}
	config.SafePorts = []int{}
	err = config.Validate()
	if err == nil {
		t.Error("Empty safe ports should fail validation")
	}
}

func TestSquidConfigData_GetConfigSummary(t *testing.T) {
	config := &SquidConfigData{
		HttpPort:        3128,
		CacheDir:        "/var/spool/squid",
		CoreDumpDir:     "/var/spool/squid",
		LocalNetworks:   []string{"192.168.0.0/16"},
		SafePorts:       []int{80, 443},
		SSLPorts:        []int{443},
		AccessRules:     []string{"http_access allow localhost"},
		RefreshPatterns: []string{"refresh_pattern . 0 20% 4320"},
		ACLs:            []ACL{{Name: "localnet", Type: "src", Value: "192.168.0.0/16"}},
	}

	summary := config.GetConfigSummary()

	if summary["http_port"] != 3128 {
		t.Errorf("Expected http_port 3128, got %v", summary["http_port"])
	}

	if summary["local_networks"] != 1 {
		t.Errorf("Expected 1 local network, got %v", summary["local_networks"])
	}

	if summary["safe_ports"] != 2 {
		t.Errorf("Expected 2 safe ports, got %v", summary["safe_ports"])
	}
}
