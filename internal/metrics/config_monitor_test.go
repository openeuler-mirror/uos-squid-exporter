// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
	"os"
	"testing"
	"time"
)

func TestConfigFileMonitor_New(t *testing.T) {
	configPath := "/etc/squid/squid.conf"
	monitor := NewConfigFileMonitor(configPath)

	if monitor == nil {
		t.Fatal("Expected monitor to be created, got nil")
	}

	if monitor.configPath != configPath {
		t.Errorf("Expected config path %s, got %s", configPath, monitor.configPath)
	}

	if monitor.running {
		t.Error("Expected monitor to not be running initially")
	}
}

func TestConfigFileMonitor_StartStop(t *testing.T) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "squid_monitor_test.conf")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	monitor := NewConfigFileMonitor(tmpFile.Name())

	// 测试启动
	monitor.Start(100 * time.Millisecond) // 100ms检查间隔
	if !monitor.IsRunning() {
		t.Error("Expected monitor to be running after Start()")
	}

	// 等待一段时间
	time.Sleep(200 * time.Millisecond)

	// 测试停止
	monitor.Stop()
	if monitor.IsRunning() {
		t.Error("Expected monitor to not be running after Stop()")
	}
}

func TestConfigFileMonitor_FileChangeDetection(t *testing.T) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "squid_monitor_change_test.conf")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	monitor := NewConfigFileMonitor(tmpFile.Name())

	// 启动监控
	monitor.Start(50 * time.Millisecond) // 50ms检查间隔
	defer monitor.Stop()

	// 获取刷新通道
	refreshCh := monitor.GetRefreshChannel()

	// 等待初始状态稳定
	time.Sleep(100 * time.Millisecond)

	// 修改文件
	err = os.WriteFile(tmpFile.Name(), []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// 等待检测到变化
	select {
	case <-refreshCh:
		// 成功检测到变化
	case <-time.After(200 * time.Millisecond):
		t.Error("Expected to detect file change within 200ms")
	}
}

func TestConfigFileMonitor_GetFileInfo(t *testing.T) {
	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "squid_monitor_info_test.conf")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入一些内容
	content := "test content for file info"
	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	monitor := NewConfigFileMonitor(tmpFile.Name())

	// 获取文件信息
	modTime, size, err := monitor.GetFileInfo()
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if size != int64(len(content)) {
		t.Errorf("Expected size %d, got %d", len(content), size)
	}

	if modTime.IsZero() {
		t.Error("Expected modification time to not be zero")
	}
}

func TestConfigFileMonitor_NonExistentFile(t *testing.T) {
	monitor := NewConfigFileMonitor("/nonexistent/file.conf")

	// 获取文件信息应该返回错误
	_, _, err := monitor.GetFileInfo()
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// 启动监控应该不会panic
	monitor.Start(100 * time.Millisecond)
	defer monitor.Stop()

	// 等待一段时间确保没有panic
	time.Sleep(200 * time.Millisecond)
}
