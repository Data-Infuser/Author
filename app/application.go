package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gitlab.com/promptech1/infuser-author/app/ctx"
	"gitlab.com/promptech1/infuser-author/constant"
	"gitlab.com/promptech1/infuser-author/database"
	server "gitlab.com/promptech1/infuser-author/grpc"
	"gitlab.com/promptech1/infuser-author/model"
	"gopkg.in/yaml.v2"
	"xorm.io/xorm"
	"xorm.io/xorm/log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Application define a mode of running app
type Application struct {
	Ctx     *ctx.Context
	Context context.Context
	server  *server.Server
}

// New constructor
func New(context context.Context) (*Application, error) {
	var err error

	a := new(Application)
	a.Ctx = new(ctx.Context)
	a.Context = context

	env := os.Getenv("AUTHOR_ENV")
	if len(env) == 0 || env != constant.ServiceProd {
		a.Ctx.Mode = constant.ServiceDev
	} else {
		a.Ctx.Mode = constant.ServiceProd
	}

	a.Ctx.DBConfigFileName = fmt.Sprintf("config/%s/database.yaml", a.Ctx.Mode)
	a.Ctx.RedisConfigFileName = fmt.Sprintf("config/%s/redis.yaml", a.Ctx.Mode)

	if err = a.initConfig(); err != nil {
		return nil, err
	}

	if err = a.initLogger(); err != nil {
		return nil, err
	}

	if err = a.initDB(); err != nil {
		return nil, err
	}

	a.initRedis(a.Context)

	return a, nil
}

// Run starts application
func (a *Application) Run(network, addr string) {
	a.server = server.New(a.Ctx, a.Context)
	a.server.Run(network, addr)
}

func (a *Application) initConfig() error {
	var file []byte
	var err error

	ballast := make([]byte, 10<<30)
	_ = ballast

	if file, err = ioutil.ReadFile("config/config.yaml"); err != nil {
		return err
	}
	if err = yaml.Unmarshal(file, &a.Ctx.Config); err != nil {
		return err
	}

	// Load DB Config
	if file, err = ioutil.ReadFile(a.Ctx.DBConfigFileName); err != nil {
		return err
	}
	if err = yaml.Unmarshal(file, &a.Ctx.DBConfig); err != nil {
		return err
	}

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
	a.Ctx.Logger.Info("DB Config ==========================================")
	a.Ctx.Logger.Info(dbConfig)
	connectURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	a.Ctx.Logger.Info(connectURL)

	if a.Ctx.Orm, err = xorm.NewEngine(dbConfig.DBType, connectURL); err != nil {
		return err
	}

	if a.Ctx.Mode == constant.ServiceDev {
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
	if err = a.Ctx.Orm.Sync2(new(model.Operation)); err != nil {
		return err
	}
	if err = a.Ctx.Orm.Sync2(new(model.Traffic)); err != nil {
		return err
	}

	return nil
}

func (a *Application) initRedis(context context.Context) {
	redisConfig := a.Ctx.RedisConfig
	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisConfig.Addr, redisConfig.Port),
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		MinIdleConns: redisConfig.MinIdleConns,
		PoolSize:     redisConfig.PoolSize,
	})

	a.Ctx.RedisDB = database.NewRedisDB(context, redisClient)
}

func (a *Application) initLogger() error {
	logger := logrus.New()

	if _, err := os.Stat("log"); os.IsNotExist(err) {
		os.Mkdir("log", 0777)
	}

	dir, _ := os.Getwd()
	fmt.Println("CWD:", dir)

	if a.Ctx.Mode == constant.ServiceDev {
		logger.SetLevel(logrus.DebugLevel)
		// logger.SetFormatter(&logrus.JSONFormatter{})
		logger.Out = os.Stdout
	} else {
		logger.SetLevel(logrus.InfoLevel)
		//log
		file, _ := os.OpenFile(a.Ctx.Config.LoggerConfig.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		logger.Out = file
	}

	a.Ctx.Logger = logger.WithFields(logrus.Fields{
		"tag": a.Ctx.Config.LoggerConfig.Tag,
		"id":  os.Getpid(),
	})

	return nil
}
