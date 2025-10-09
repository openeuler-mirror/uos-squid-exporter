// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
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

// Describe 实现prometheus.Collector接口
func (c *SquidConfigFilesCollector) Describe(ch chan<- *prometheus.Desc) {
	c.filesCount.Describe(ch)
	c.totalSize.Describe(ch)
	c.lastScanTime.Describe(ch)
	c.scanSuccess.Describe(ch)
	ch <- c.fileInfo
	ch <- c.fileTypesCount
	ch <- c.recentlyChanged
}

// Collect 实现prometheus.Collector接口
func (c *SquidConfigFilesCollector) Collect(ch chan<- prometheus.Metric) {
	// 扫描配置目录
	files, err := c.scanConfigDirectory()
	if err != nil {
		logrus.Errorf("Failed to scan squid config directory %s: %v", c.configDir, err)
		c.scanSuccess.Set(0)
		ch <- c.scanSuccess
		return
	}

	// 设置扫描成功指标
	c.scanSuccess.Set(1)
	c.lastScanTime.Set(float64(time.Now().Unix()))

	// 计算统计信息
	var totalSize int64
	extensionCount := make(map[string]int)
	var recentlyChangedCount int

	// 设置文件信息指标
	for _, file := range files {
		totalSize += file.Size

		// 统计文件类型
		if file.Extension != "" {
			extensionCount[file.Extension]++
		} else {
			extensionCount["none"]++
		}

		// 检查最近修改的文件（1小时内）
		if time.Since(file.ModTime) <= time.Hour {
			recentlyChangedCount++
		}

		// 发送文件信息指标
		ch <- prometheus.MustNewConstMetric(
			c.fileInfo,
			prometheus.GaugeValue,
			float64(file.Size),
			file.Name,
			file.Path,
			file.Extension,
			file.Permissions,
		)
	}

	// 设置基础指标
	c.filesCount.Set(float64(len(files)))
	c.totalSize.Set(float64(totalSize))

	// 发送基础指标
	ch <- c.filesCount
	ch <- c.totalSize
	ch <- c.lastScanTime
	ch <- c.scanSuccess

	// 发送文件类型统计指标
	for extension, count := range extensionCount {
		ch <- prometheus.MustNewConstMetric(
			c.fileTypesCount,
			prometheus.GaugeValue,
			float64(count),
			extension,
		)
	}

	// 发送最近修改文件数量指标
	ch <- prometheus.MustNewConstMetric(
		c.recentlyChanged,
		prometheus.GaugeValue,
		float64(recentlyChangedCount),
	)
}

// scanConfigDirectory 扫描配置目录并返回文件信息
func (c *SquidConfigFilesCollector) scanConfigDirectory() ([]ConfigFileInfo, error) {
	var files []ConfigFileInfo

	err := filepath.Walk(c.configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 如果无法访问某个文件或目录，记录警告但继续处理其他文件
			logrus.Warnf("Unable to access %s: %v", path, err)
			return nil
		}

		// 跳过目录本身（只处理目录中的文件）
		if path == c.configDir {
			return nil
		}

		// 构建文件信息
		fileInfo := ConfigFileInfo{
			Name:        info.Name(),
			Path:        path,
			Size:        info.Size(),
			ModTime:     info.ModTime(),
			IsDirectory: info.IsDir(),
			IsRegular:   info.Mode().IsRegular(),
			Permissions: info.Mode().String(),
			Extension:   strings.TrimPrefix(filepath.Ext(info.Name()), "."),
		}

		files = append(files, fileInfo)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// GetConfigDir 获取配置目录路径
func (c *SquidConfigFilesCollector) GetConfigDir() string {
	return c.configDir
}
