// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// ConfigFileInfo 表示配置文件的详细信息
type ConfigFileInfo struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"mod_time"`
	IsDirectory bool      `json:"is_directory"`
	IsRegular   bool      `json:"is_regular"`
	Permissions string    `json:"permissions"`
	Extension   string    `json:"extension"`
}

// SquidConfigFilesCollector squid配置文件列表收集器
type SquidConfigFilesCollector struct {
	configDir string

	// Prometheus指标
	filesCount      prometheus.Gauge
	totalSize       prometheus.Gauge
	lastScanTime    prometheus.Gauge
	scanSuccess     prometheus.Gauge
	fileInfo        *prometheus.Desc
	fileTypesCount  *prometheus.Desc
	recentlyChanged *prometheus.Desc
}

// NewSquidConfigFilesCollector 创建新的squid配置文件列表收集器
func NewSquidConfigFilesCollector(configDir string) *SquidConfigFilesCollector {
	collector := &SquidConfigFilesCollector{
		configDir: configDir,

		// 初始化Prometheus指标
		filesCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "files_count",
			Help:      "Total number of configuration files in squid config directory",
		}),

		totalSize: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "files_total_size_bytes",
			Help:      "Total size of all configuration files in bytes",
		}),

		lastScanTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "last_scan_timestamp",
			Help:      "Timestamp of the last successful directory scan",
		}),

		scanSuccess: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "squid_config",
			Name:      "scan_success",
			Help:      "Whether the last directory scan was successful (1) or not (0)",
		}),

		fileInfo: prometheus.NewDesc(
			"squid_config_file_info",
			"Information about squid configuration files",
			[]string{"file_name", "file_path", "file_extension", "permissions"},
			nil,
		),

		fileTypesCount: prometheus.NewDesc(
			"squid_config_file_types_count",
			"Count of files by extension type",
			[]string{"extension"},
			nil,
		),

		recentlyChanged: prometheus.NewDesc(
			"squid_config_files_recently_changed",
			"Number of files changed in the last hour",
			nil,
			nil,
		),
	}

	return collector
}
