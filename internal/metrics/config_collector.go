// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
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
