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
