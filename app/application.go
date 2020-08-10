package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
	"gitlab.com/promptech1/infuser-author/app/ctx"
	"gitlab.com/promptech1/infuser-author/constant"
	server "gitlab.com/promptech1/infuser-author/grpc"
	"gitlab.com/promptech1/infuser-author/model"
	"gopkg.in/yaml.v2"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

// Application define a mode of running app
type Application struct {
	Ctx    *ctx.Context
	server *server.Server
}

// New constructor
func New() (*Application, error) {
	var err error

	a := new(Application)
	a.Ctx = new(ctx.Context)

	a.Ctx.Context = context.Background()

	env := os.Getenv("AUTHOR_ENV")
	if len(env) == 0 || env != constant.SERVICE_PROD {
		a.Ctx.Mode = constant.SERVICE_DEV
	} else {
		a.Ctx.Mode = constant.SERVICE_PROD
	}

	a.Ctx.DBConfigFileName = fmt.Sprintf("config/%s/database.yaml", a.Ctx.Mode)
	a.Ctx.RedisConfigFileName = fmt.Sprintf("config/%s/redis.yaml", a.Ctx.Mode)

	if err = a.initConfig(); err != nil {
		return nil, err
	}

	if err = a.initDB(); err != nil {
		return nil, err
	}

	redisConfig := a.Ctx.RedisConfig
	a.Ctx.RedisDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisConfig.Addr, redisConfig.Port),
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		MinIdleConns: redisConfig.MinIdleConns,
		PoolSize:     redisConfig.PoolSize,
	})

	return a, nil
}

// Run starts application
func (a *Application) Run(network, addr string) {
	a.server = server.New(a.Ctx)
	a.server.Run(network, addr)
}

func (a *Application) initConfig() error {
	var file []byte
	var err error

	ballast := make([]byte, 10<<30)
	_ = ballast

	// Load DB Config
	glog.Info("Load DB Config ==============================")
	glog.Info(a.Ctx.DBConfigFileName)
	if file, err = ioutil.ReadFile(a.Ctx.DBConfigFileName); err != nil {
		return err
	}
	if err = yaml.Unmarshal(file, &a.Ctx.DBConfig); err != nil {
		return err
	}
	glog.Info(a.Ctx.DBConfig)

	// Load Redis Config
	if file, err = ioutil.ReadFile(a.Ctx.RedisConfigFileName); err != nil {
		return err
	}
	if err = yaml.Unmarshal(file, &a.Ctx.RedisConfig); err != nil {
		return err
	}

	return nil
}

func (a *Application) initDB() error {
	var err error

	dbConfig := a.Ctx.DBConfig
	glog.Info("DB Config ==========================================")
	glog.Info(dbConfig)
	connectURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	glog.Info(connectURL)

	if a.Ctx.Orm, err = xorm.NewEngine(dbConfig.DBType, connectURL); err != nil {
		return err
	}

	if a.Ctx.Mode == constant.SERVICE_DEV {
		a.Ctx.Orm.ShowSQL(true)
		a.Ctx.Orm.Logger().SetLevel(log.LOG_DEBUG)
	}
	a.Ctx.Orm.SetMaxIdleConns(dbConfig.IdleConns)
	a.Ctx.Orm.SetMaxOpenConns(dbConfig.MaxOpenConns)

	//migrate
	err = a.migrateDB()

	return err
}

func (a *Application) migrateDB() error {
	var err error
	if err = a.Ctx.Orm.Sync2(new(model.App)); err != nil {
		return err
	}
	if err = a.Ctx.Orm.Sync2(new(model.Token)); err != nil {
		return err
	}
	if err = a.Ctx.Orm.Sync2(new(model.AppToken)); err != nil {
		return err
	}
	if err = a.Ctx.Orm.Sync2(new(model.AppTokenHistory)); err != nil {
		return err
	}

	return nil
}
