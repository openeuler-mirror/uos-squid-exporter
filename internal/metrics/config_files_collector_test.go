// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// 测试描述方法
func TestSquidConfigFilesCollector_Describe(t *testing.T) {
	collector := NewSquidConfigFilesCollector("/tmp")

	ch := make(chan *prometheus.Desc, 10)
	collector.Describe(ch)

	// 收集所有描述
	var descs []*prometheus.Desc
	for i := 0; i < 7; i++ {
		select {
		case desc := <-ch:
			descs = append(descs, desc)
		case <-time.After(100 * time.Millisecond):
			break
		}
	}

	assert.GreaterOrEqual(t, len(descs), 7, "应至少描述7个指标")
}

// 测试扫描空目录
func TestSquidConfigFilesCollector_ScanEmptyDirectory(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "squid_config_test_empty")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	collector := NewSquidConfigFilesCollector(tmpDir)

	files, err := collector.scanConfigDirectory()
	require.NoError(t, err)
	assert.Empty(t, files, "空目录应返回空文件列表")
}

// 测试扫描包含文件的目录
func TestSquidConfigFilesCollector_ScanDirectoryWithFiles(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "squid_config_test_files")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 创建测试文件
	testFiles := []struct {
		name      string
		content   string
		extension string
	}{
		{"squid.conf", "http_port 3128\ncache_dir ufs /var/spool/squid 100 16 256", "conf"},
		{"acl.conf", "acl localnet src 192.168.0.0/16", "conf"},
		{"mime.conf", "text/html html htm", "conf"},
		{"readme.txt", "Squid configuration files", "txt"},
		{"backup.conf.bak", "backup file", "bak"},
	}

	for _, tf := range testFiles {
		filePath := filepath.Join(tmpDir, tf.name)
		err := os.WriteFile(filePath, []byte(tf.content), 0644)
		require.NoError(t, err)
	}

	collector := NewSquidConfigFilesCollector(tmpDir)

	files, err := collector.scanConfigDirectory()
	require.NoError(t, err)
	assert.Len(t, files, len(testFiles), "应找到所有测试文件")

	// 验证文件信息
	for _, file := range files {
		assert.NotEmpty(t, file.Name, "文件名不应为空")
		assert.NotEmpty(t, file.Path, "文件路径不应为空")
		assert.False(t, file.IsDirectory, "文件不应是目录")
		assert.True(t, file.IsRegular, "文件应是常规文件")
		assert.NotZero(t, file.Size, "文件大小不应为零")
		assert.NotZero(t, file.ModTime, "修改时间不应为零")
		assert.NotEmpty(t, file.Permissions, "权限字符串不应为空")
	}
}

// 测试扫描包含子目录的目录
func TestSquidConfigFilesCollector_ScanDirectoryWithSubdirectories(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "squid_config_test_subdirs")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// 创建子目录
	subDir := filepath.Join(tmpDir, "conf.d")
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	// 在子目录中创建文件
	filePath := filepath.Join(subDir, "extra.conf")
	err = os.WriteFile(filePath, []byte("acl extra src 10.0.0.0/8"), 0644)
	require.NoError(t, err)

	collector := NewSquidConfigFilesCollector(tmpDir)

	files, err := collector.scanConfigDirectory()
	require.NoError(t, err)
	assert.Len(t, files, 2, "应找到子目录和文件") // 子目录本身和其中的文件

	// 验证子目录信息
	var foundSubDir bool
	var foundFile bool
	for _, file := range files {
		if file.IsDirectory {
			foundSubDir = true
			assert.Equal(t, "conf.d", file.Name, "子目录名应匹配")
		} else {
			foundFile = true
			assert.Equal(t, "extra.conf", file.Name, "文件名应匹配")
		}
	}

	assert.True(t, foundSubDir, "应找到子目录")
	assert.True(t, foundFile, "应找到文件")
}

// 测试扫描不存在的目录
// func TestSquidConfigFilesCollector_ScanNonExistentDirectory(t *testing.T) {
// 	nonExistentDir := "/this/path/does/not/exist/12345"
// 	collector := NewSquidConfigFilesCollector(nonExistentDir)

// 	files, err := collector.scanConfigDirectory()
// 	assert.NoError(t, err, "扫描不存在的目录不应返回错误")
// 	assert.Empty(t, files, "不存在的目录应返回空文件列表")
// }
