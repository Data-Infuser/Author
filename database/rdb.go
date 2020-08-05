package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/promptech1/infuser-author/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	_ "github.com/go-sql-driver/mysql"
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

func ConnDB() *gorm.DB{
	ballast := make([]byte, 10<<30)
	_ = ballast

	filename, _ := filepath.Abs("config/database.yaml")
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("filename: %v", filename)

	var envConfig DBEnvConfig
	var dbConfig DBConfig

	err = yaml.Unmarshal(file, &envConfig)
	if err != nil {
		panic(err.Error())
	}

	env := os.Getenv("MANAGER_ENV")
	if len(env) == 0 || env != "prod"{
		dbConfig = envConfig.Dev
	} else {
		dbConfig = envConfig.Prod
	}

	connectURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

	db, err := gorm.Open(dbConfig.DBType, connectURL)
	if err != nil {
		panic(err.Error())
	}

	db.DB().SetMaxIdleConns(dbConfig.IdleConns)
	db.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)

	migrate(db)

	return db
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&model.AppToken{})
	db.AutoMigrate(&model.Token{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.App{})
	db.AutoMigrate(&model.AppTokenHistory{})
}