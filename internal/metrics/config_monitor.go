// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ConfigFileMonitor 配置文件监控器
type ConfigFileMonitor struct {
	configPath     string
	lastModTime    time.Time
	lastSize       int64
	mu             sync.RWMutex
	refreshChannel chan struct{}
	stopChannel    chan struct{}
	running        bool
}

// NewConfigFileMonitor 创建新的配置文件监控器
func NewConfigFileMonitor(configPath string) *ConfigFileMonitor {
	return &ConfigFileMonitor{
		configPath:     configPath,
		refreshChannel: make(chan struct{}, 1),
		stopChannel:    make(chan struct{}),
		running:        false,
	}
}

// Start 启动配置文件监控
func (m *ConfigFileMonitor) Start(interval time.Duration) {
	if m.running {
		return
	}

	m.running = true
	go m.monitorLoop(interval)
	logrus.Infof("Config file monitor started for: %s", m.configPath)
}

// Stop 停止配置文件监控
func (m *ConfigFileMonitor) Stop() {
	if !m.running {
		return
	}

	close(m.stopChannel)
	m.running = false
	logrus.Info("Config file monitor stopped")
}

// GetRefreshChannel 获取刷新通道
func (m *ConfigFileMonitor) GetRefreshChannel() <-chan struct{} {
	return m.refreshChannel
}

// monitorLoop 监控循环
func (m *ConfigFileMonitor) monitorLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// 初始化文件状态
	m.updateFileState()

	for {
		select {
		case <-ticker.C:
			if m.checkFileChanged() {
				logrus.Infof("Config file changed: %s", m.configPath)
				select {
				case m.refreshChannel <- struct{}{}:
				default:
					// 通道已满，跳过本次通知
				}
			}
		case <-m.stopChannel:
			return
		}
	}
}

// checkFileChanged 检查文件是否发生变化
func (m *ConfigFileMonitor) checkFileChanged() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	info, err := os.Stat(m.configPath)
	if err != nil {
		logrus.Warnf("Failed to stat config file %s: %v", m.configPath, err)
		return false
	}

	// 检查修改时间和文件大小
	if !info.ModTime().Equal(m.lastModTime) || info.Size() != m.lastSize {
		m.lastModTime = info.ModTime()
		m.lastSize = info.Size()
		return true
	}

	return false
}
