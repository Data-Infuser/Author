package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"gitlab.com/promptech1/infuser-author/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

// DBConfig : Database Config
type DBConfig struct {
	DBName   string `yaml:"dbName"`
	DBType   string `yaml:"dbType"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	IdleConns int `yaml:"idleConns"`
	MaxOpenConns int `yaml:"maxOpenConns"`
}

type DBEnvConfig struct {
	Dev DBConfig `yaml:"dev"`
	Prod DBConfig `yaml:"prod"`
}

func ConnDB() *xorm.Engine{
	ballast := make([]byte, 10<<30)
	_ = ballast

	filename, _ := filepath.Abs("config/database.yaml")
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}
	glog.Infof("filename: %v", filename)

	var envConfig DBEnvConfig
	var dbConfig DBConfig

	err = yaml.Unmarshal(file, &envConfig)
	if err != nil {
		panic(err.Error())
	}

	var isProd = false
	env := os.Getenv("AUTHOR_ENV")
	if len(env) == 0 || env != "prod"{
		dbConfig = envConfig.Dev
	} else {
		dbConfig = envConfig.Prod
		isProd = true
	}

	connectURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

	engine, err := xorm.NewEngine(dbConfig.DBType, connectURL)
	if err != nil {
		panic(err.Error())
	}

	if !isProd {
		engine.ShowSQL(true)
		engine.Logger().SetLevel(log.LOG_DEBUG)
	}

	engine.SetMaxIdleConns(dbConfig.IdleConns)
	engine.SetMaxOpenConns(dbConfig.MaxOpenConns)

	migrate(engine)

	return engine
}

func migrate(engine *xorm.Engine) {
	engine.Sync2(new(model.App))
	engine.Sync2(new(model.Token))
	engine.Sync2(new(model.AppToken))
}