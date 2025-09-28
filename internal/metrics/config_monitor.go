// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
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
