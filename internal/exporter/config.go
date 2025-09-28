// SPDX-FileCopyrightText: 2025 UnionTech Software Technology Co., Ltd.
// SPDX-License-Identifier: MIT
package exporter

import (
	"os"
	"time"
	"uos-squid-exporter/pkg/logger"
	"uos-squid-exporter/pkg/utils"

	"github.com/alecthomas/kingpin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	Configfile      *string
	SquidConfigPath *string
	DefaultConfig   = Config{
		Logging: logger.Config{
			Level:   "debug",
			LogPath: "/var/log/uos-exporter/squid_exporter.log",
			MaxSize: "10MB",
			MaxAge:  time.Hour * 24 * 7},
		Address:         "127.0.0.1",
		Port:            8080,
		MetricsPath:     "/metrics",
		SquidConfigPath: "/etc/squid/squid.conf",
	}
)

func init() {
	Configfile = kingpin.Flag("config", "Configuration file").
		Short('c').
		Default("/etc/uos-exporter/squid-exporter.yaml").
		String()

	SquidConfigPath = kingpin.Flag("squid-config", "Path to squid configuration file").
		Default("/etc/squid/squid.conf").
		String()
}

type Config struct {
	Logging         logger.Config `yaml:"log"`
	Address         string        `yaml:"address"`
	Port            int           `yaml:"port"`
	MetricsPath     string        `yaml:"metricsPath"`
	SquidConfigPath string        `yaml:"squidConfigPath"`
}

func Unpack(config interface{}) error {
	if !utils.FileExists(*Configfile) {
		logrus.Errorf("%s file not found", *Configfile)
	} else {
		file, err := os.Open(*Configfile)
		if err != nil {
			logrus.Error("Failed to open config file: ", err)
			return err
		}
		err = yaml.NewDecoder(file).Decode(config)
		if err != nil {
			return err
		}
	}
	return nil
}
