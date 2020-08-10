package ctx

import (
	"context"

	"github.com/go-redis/redis/v8"
	"xorm.io/xorm"
)

type Context struct {
	Context             context.Context
	Mode                string
	Orm                 *xorm.Engine
	RedisDB             *redis.Client
	DBConfig            *DBConfig
	DBConfigFileName    string
	RedisConfig         *RedisConfig
	RedisConfigFileName string
}

// DBConfig : Database Config
type DBConfig struct {
	DBName       string `yaml:"dbName"`
	DBType       string `yaml:"dbType"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	IdleConns    int    `yaml:"idleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
}

type RedisConfig struct {
	Addr         string `yaml:"addr"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	MinIdleConns int    `yaml:"minIdleConns"`
	PoolSize     int    `yaml:"poolSize"`
}
