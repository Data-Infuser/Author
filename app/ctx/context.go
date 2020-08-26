package ctx

import (
	log "github.com/sirupsen/logrus"
	"gitlab.com/promptech1/infuser-author/database"
	"xorm.io/xorm"
)

type Context struct {
	Mode                string
	Logger              *log.Entry
	Orm                 *xorm.Engine
	RedisDB             *database.RedisDB
	Config              *Config
	DBConfig            *DBConfig
	DBConfigFileName    string
	RedisConfig         *RedisConfig
	RedisConfigFileName string
}

type Config struct {
	LoggerConfig LoggerConfig `yaml:"logger"`
}

type LoggerConfig struct {
	Mode     string `yaml:"mode"`
	Tag      string `yaml:"tag"`
	FileName string `yaml:"fileName"`
	Id       string
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
