package gorm

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Engine struct {
	client *gorm.DB
}

func NewEngine(dbInfo MysqlConf, plugins ...gorm.Plugin) (*Engine, error) {
	db, err := NewClient(dbInfo, plugins...)
	if err != nil {
		return nil, err
	}
	return &Engine{client: db}, nil
}

func NewClient(dbInfo MysqlConf, plugins ...gorm.Plugin) (*gorm.DB, error) {
	config := new(gorm.Config)
	for i, v := range plugins {
		if logger, ok := v.(*Logger); ok {
			plugins = append(plugins[:i], plugins[i+1:]...)
			config.Logger = logger
			break
		}
	}
	db, err := gorm.Open(mysql.Open(getConnURL(dbInfo)), config)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if dbInfo.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(dbInfo.MaxIdle)
	}
	if dbInfo.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(dbInfo.MaxOpen)
	}
	for _, plugin := range plugins {
		err = db.Use(plugin)
		if err != nil {
			return nil, err
		}
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getConnURL(info MysqlConf) (url string) {
	url = fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		info.User,
		info.Password,
		info.Host,
		info.Port,
		info.Database,
		info.Charset)
	return
}
