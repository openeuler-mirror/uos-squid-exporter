// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// SquidConfigCollector squid配置文件指标收集器
type SquidConfigCollector struct {
	configPath string
	configData *SquidConfigData
	parser     *SquidConfigParser
	monitor    *ConfigFileMonitor

	// Prometheus指标
	configUp          prometheus.Gauge
	httpPort          prometheus.Gauge
	cacheDirExists    prometheus.Gauge
	coredumpDirExists prometheus.Gauge
	localNetworks     prometheus.Gauge
	safePorts         prometheus.Gauge
	sslPorts          prometheus.Gauge
	accessRules       prometheus.Gauge
	refreshPatterns   prometheus.Gauge
	acls              prometheus.Gauge

	// 配置摘要指标
	configSummary *prometheus.Desc
}

// NewSquidConfigCollector 创建新的squid配置指标收集器
func NewSquidConfigCollector(configPath string) *SquidConfigCollector {
	collector := &SquidConfigCollector{
		configPath: configPath,
		parser:     NewSquidConfigParser(configPath),
		monitor:    NewConfigFileMonitor(configPath),

		// 初始化Prometheus指标
		configUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "up",
			Help:      "Whether the squid config file is accessible and parseable",
		}),

		httpPort: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "http_port",
			Help:      "HTTP port configured in squid.conf",
		}),

		cacheDirExists: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "cache_dir_exists",
			Help:      "Whether the cache directory exists (1) or not (0)",
		}),

		coredumpDirExists: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "coredump_dir_exists",
			Help:      "Whether the coredump directory exists (1) or not (0)",
		}),

		localNetworks: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "local_networks_count",
			Help:      "Number of local networks defined in ACL",
		}),

		safePorts: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "safe_ports_count",
			Help:      "Number of safe ports defined in ACL",
		}),

		sslPorts: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "ssl_ports_count",
			Help:      "Number of SSL ports defined in ACL",
		}),

		accessRules: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "access_rules_count",
			Help:      "Number of http_access rules defined",
		}),

		refreshPatterns: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "refresh_patterns_count",
			Help:      "Number of refresh_pattern rules defined",
		}),

		acls: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "acls_count",
			Help:      "Number of ACL definitions",
		}),

		configSummary: prometheus.NewDesc(
			"squid_config_summary",
			"Summary of squid configuration",
			[]string{"config_file", "http_port", "cache_dir", "coredump_dir"},
			nil,
		),
	}

	// 启动配置文件监控
	collector.monitor.Start(30 * time.Second) // 每30秒检查一次

	return collector
}

// Describe 实现prometheus.Collector接口
func (c *SquidConfigCollector) Describe(ch chan<- *prometheus.Desc) {
	c.configUp.Describe(ch)
	c.httpPort.Describe(ch)
	c.cacheDirExists.Describe(ch)
	c.coredumpDirExists.Describe(ch)
	c.localNetworks.Describe(ch)
	c.safePorts.Describe(ch)
	c.sslPorts.Describe(ch)
	c.accessRules.Describe(ch)
	c.refreshPatterns.Describe(ch)
	c.acls.Describe(ch)
	ch <- c.configSummary
}

// Collect 实现prometheus.Collector接口
func (c *SquidConfigCollector) Collect(ch chan<- prometheus.Metric) {
	// 解析配置文件
	configData, err := c.parser.Parse()
	if err != nil {
		logrus.Errorf("Failed to parse squid config file %s: %v", c.configPath, err)
		c.configUp.Set(0)
		ch <- c.configUp
		return
	}

	c.configData = configData

	// 验证配置
	if err := configData.Validate(); err != nil {
		logrus.Errorf("Squid config validation failed: %v", err)
		c.configData = nil // 验证失败时清空配置数据
		c.configUp.Set(0)
		ch <- c.configUp
		return
	}

	// 设置配置可用性指标
	c.configUp.Set(1)

	// 设置HTTP端口指标
	c.httpPort.Set(float64(configData.HttpPort))

	// 检查缓存目录是否存在
	cacheDirExists := 0.0
	if configData.CacheDir != "" {
		// 提取缓存目录路径（去掉ufs等参数）
		cacheDirPath := extractCacheDirPath(configData.CacheDir)
		if cacheDirPath != "" && dirExists(cacheDirPath) {
			cacheDirExists = 1.0
		}
	}
	c.cacheDirExists.Set(cacheDirExists)

	// 检查核心转储目录是否存在
	coredumpDirExists := 0.0
	if configData.CoreDumpDir != "" && dirExists(configData.CoreDumpDir) {
		coredumpDirExists = 1.0
	}
	c.coredumpDirExists.Set(coredumpDirExists)

	// 设置计数指标
	c.localNetworks.Set(float64(len(configData.LocalNetworks)))
	c.safePorts.Set(float64(len(configData.SafePorts)))
	c.sslPorts.Set(float64(len(configData.SSLPorts)))
	c.accessRules.Set(float64(len(configData.AccessRules)))
	c.refreshPatterns.Set(float64(len(configData.RefreshPatterns)))
	c.acls.Set(float64(len(configData.ACLs)))

	// 发送所有指标
	ch <- c.configUp
	ch <- c.httpPort
	ch <- c.cacheDirExists
	ch <- c.coredumpDirExists
	ch <- c.localNetworks
	ch <- c.safePorts
	ch <- c.sslPorts
	ch <- c.accessRules
	ch <- c.refreshPatterns
	ch <- c.acls

	// 发送配置摘要指标
	configFile := filepath.Base(c.configPath)
	cacheDir := configData.CacheDir
	if cacheDir == "" {
		cacheDir = "not_configured"
	}
	coredumpDir := configData.CoreDumpDir
	if coredumpDir == "" {
		coredumpDir = "not_configured"
	}

	ch <- prometheus.MustNewConstMetric(
		c.configSummary,
		prometheus.GaugeValue,
		1.0,
		configFile,
		strconv.Itoa(configData.HttpPort),
		cacheDir,
		coredumpDir,
	)
}

// extractCacheDirPath 从cache_dir配置中提取目录路径
func extractCacheDirPath(cacheDirConfig string) string {
	// cache_dir格式通常是: ufs /path/to/dir size dirs subdirs
	parts := strings.Fields(cacheDirConfig)
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}
