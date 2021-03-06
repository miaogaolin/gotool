package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var engine *Engine

type Engine struct {
	*gorm.DB
}

func Connect(conf *Config) (err error) {
	if engine != nil {
		return
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DataBase )
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	engine = &Engine{
		db,
	}
	return
}

func (Engine) DataBase() *Engine {
	return engine
}
