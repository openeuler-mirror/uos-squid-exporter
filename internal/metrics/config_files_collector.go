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
