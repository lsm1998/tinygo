package mysqlx

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

type Options struct {
	conf Config
}

type Option func(*Options)

// Config mysql的配置字段
type Config struct {
	Write       Configx   `json:"write" yaml:"write"`
	ReadOnly    []Configx `json:"read_only" yaml:"read_only"`
	MaxIdleConn int       `json:"max_idle_conn" yaml:"max_idle_conn"`
	MaxOpenConn int       `json:"max_open_conn" yaml:"max_open_conn"`
}

type Configx struct {
	Addr     string `json:"addr" yaml:"addr"`
	Port     int    `json:"port" yaml:"port"`
	DBName   string `json:"db_name" yaml:"db_name"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func WithConfig(conf Config) Option {
	return func(k *Options) {
		k.conf = conf
	}
}

func Must(opts ...Option) *gorm.DB {
	e := &Options{}
	for _, opt := range opts {
		opt(e)
	}
	db, err := gorm.Open(mysql.Open(dsn(e.conf.Write)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	var replicas = make([]gorm.Dialector, 0, len(e.conf.ReadOnly))
	for _, v := range e.conf.ReadOnly {
		replicas = append(replicas, mysql.Open(dsn(v)))
	}
	err = db.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(dsn(e.conf.Write))},
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(e.conf.MaxIdleConn).
			SetMaxOpenConns(e.conf.MaxOpenConn),
	)
	if err != nil {
		panic(err)
	}
	return db
}

func dsn(c Configx) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=Local",
		c.Username,
		c.Password,
		c.Addr,
		c.Port,
		c.DBName)
}
