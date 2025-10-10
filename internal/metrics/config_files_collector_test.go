// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试创建新的配置收集器
func TestNewSquidConfigFilesCollector(t *testing.T) {
	configDir := "/etc/squid"
	collector := NewSquidConfigFilesCollector(configDir)

	assert.NotNil(t, collector, "收集器不应为空")
	assert.Equal(t, configDir, collector.GetConfigDir(), "配置目录应匹配")
	assert.NotNil(t, collector.filesCount, "文件计数指标不应为空")
	assert.NotNil(t, collector.totalSize, "总大小指标不应为空")
	assert.NotNil(t, collector.lastScanTime, "最后扫描时间指标不应为空")
	assert.NotNil(t, collector.scanSuccess, "扫描成功指标不应为空")
	assert.NotNil(t, collector.fileInfo, "文件信息描述不应为空")
	assert.NotNil(t, collector.fileTypesCount, "文件类型计数描述不应为空")
	assert.NotNil(t, collector.recentlyChanged, "最近修改文件描述不应为空")
}

// 测试获取配置目录
func TestSquidConfigFilesCollector_GetConfigDir(t *testing.T) {
	configDir := "/tmp/squid_config"
	collector := NewSquidConfigFilesCollector(configDir)

	assert.Equal(t, configDir, collector.GetConfigDir(), "配置目录应匹配")
}
