package ctx

import (
	"gitlab.com/promptech1/infuser-author/database"
	"xorm.io/xorm"
)

type Context struct {
	Mode                string
	Orm                 *xorm.Engine
	RedisDB             *database.RedisDB
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
