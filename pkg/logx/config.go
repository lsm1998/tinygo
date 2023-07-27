package logx

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"time"
)

type Config struct {
	Level      string `json:"level" yaml:"level"`
	Formatter  string `json:"formatter" yaml:"formatter"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups"`
	MaxAge     int    `json:"max_age" yaml:"max_age"`
	MaxSize    int    `json:"max_size" yaml:"max_size"`
	Compress   bool   `json:"compress" yaml:"compress"`
	Filename   string `json:"filename" yaml:"filename"`
}

func (c Config) GetLevel() logrus.Level {
	level, err := logrus.ParseLevel(c.Level)
	if err == nil {
		return level
	}
	level = logrus.Level(cast.ToUint32(c.Level))
	if level > logrus.TraceLevel {
		return logrus.ErrorLevel // 默认为ErrorLevel
	}
	return level
}

func (c Config) GetFormatter() logrus.Formatter {
	switch c.Formatter {
	case "text":
		return &logrus.TextFormatter{}
	default:
		return &logrus.JSONFormatter{} // 默认为JSON格式
	}
}

func (c Config) GetMaxBackups() int {
	if c.MaxBackups <= 0 {
		return 3
	}
	return c.MaxBackups
}

func (c Config) GetMaxAge() int {
	if c.MaxAge <= 0 {
		return 30
	}
	return c.MaxAge
}

func (c Config) GetMaxSize() int {
	if c.MaxSize <= 0 {
		return 100
	}
	return c.MaxSize
}

func (c Config) GetCompress() bool {
	return c.Compress
}

func (c Config) GetFilename() string {
	if c.Filename == "" || c.Filename == "default" {
		return time.Now().Format("2006-01-02") + ".log"
	}
	return c.Filename
}
