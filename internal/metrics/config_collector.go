// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
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
