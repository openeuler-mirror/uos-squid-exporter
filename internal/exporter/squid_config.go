// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package exporter

import (
	"uos-squid-exporter/internal/metrics"

	"github.com/sirupsen/logrus"
)

// InitSquidCollector 初始化Squid收集器
func InitSquidCollector(squidConfigPath string) {
	logrus.Info("Initializing Squid collector...")

	// 创建基础的Squid配置
	squidConfig := createSquidConfig()

	logrus.Infof("Squid collector initialized with hostname: %s, port: %d",
		squidConfig.Hostname, squidConfig.Port)

	// 注册基础指标收集器
	registerBasicCollectors(squidConfig)

	// 注册配置文件收集器
	registerConfigCollector(squidConfigPath)

	logrus.Info("Squid collector initialization completed")
}

// SquidConfig Squid配置结构
type SquidConfig struct {
	Hostname     string
	Port         int
	Login        string
	Password     string
	ExtractTimes bool
	Headers      []string
}

// createSquidConfig 创建默认的Squid配置
func createSquidConfig() *SquidConfig {
	return &SquidConfig{
		Hostname:     "localhost",
		Port:         3128,
		Login:        "",
		Password:     "",
		ExtractTimes: true, // 默认启用服务时间提取
		Headers:      []string{},
	}
}

// registerBasicCollectors 注册基础指标收集器
func registerBasicCollectors(config *SquidConfig) {
	logrus.Debug("Registering basic collectors...")

	// 注册主要的Squid指标收集器
	mainCollector := metrics.NewSquidCollector(&metrics.SquidConfig{
		Hostname:     config.Hostname,
		Port:         config.Port,
		Login:        config.Login,
		Password:     config.Password,
		Headers:      config.Headers,
		ExtractTimes: config.ExtractTimes,
	})
	Register(mainCollector)

	// 注册Squid计数器指标
	counters := metrics.GetSquidCounters()
	for _, counter := range counters {
		Register(counter)
	}
	logrus.Debugf("Registered %d squid counter collectors", len(counters))

	// 注册Squid信息指标
	infos := metrics.GetSquidInfos()
	for _, info := range infos {
		Register(info)
	}
	logrus.Debugf("Registered %d squid info collectors", len(infos))

	// 如果启用了服务时间提取，注册服务时间指标
	if config.ExtractTimes {
		serviceTimes := metrics.GetSquidServiceTimes()
		for _, serviceTime := range serviceTimes {
			Register(serviceTime)
		}
		logrus.Debugf("Registered %d squid service time collectors", len(serviceTimes))
	}

	logrus.Info("Basic collectors registration completed")
}

// registerConfigCollector 注册配置文件收集器
func registerConfigCollector(configPath string) {
	logrus.Debugf("Registering config collector for path: %s", configPath)

	// 创建配置文件收集器
	configCollector := metrics.NewSquidConfigCollector(configPath)

	// 注册到Prometheus注册表
	Register(configCollector)

	logrus.Infof("Config collector registered successfully for: %s", configPath)
}
