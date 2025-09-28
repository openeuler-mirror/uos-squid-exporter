// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package metrics

import (
	"sync"
	"time"
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
