package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var DB *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.dbname"))
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("mysql connect err:", zap.Error(err))
		return
	}
	err = DB.Ping()
	if err != nil {
		zap.L().Error("mysql ping err:", zap.Error(err))
		return
	}
	fmt.Println("PONG! Mysql Connect Success!")
	DB.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	DB.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	return
}
