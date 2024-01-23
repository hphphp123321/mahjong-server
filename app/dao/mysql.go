package dao

import (
	"database/sql"
	"fmt"
	"github.com/hphphp123321/mahjong-server/app/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//// DRIVER 指定驱动
//const DRIVER = "mysql"

var DB *sql.DB
var GormDB *gorm.DB

// InitMySql 初始化MySQL数据库连接
func InitMySql() (err error) {
	// 获取yaml配置参数
	conf := global.C.DBConfig

	// 如果没有提供数据库名称，则使用默认值
	if conf.DBName == "" {
		conf.DBName = "mahjong"
	}

	// 首先连接到MySQL服务器（没有指定数据库）
	baseDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		conf.UserName,
		conf.Password,
		conf.Host,
		conf.Port,
	)
	db, err := gorm.Open(mysql.Open(baseDSN), &gorm.Config{})
	if err != nil {
		return err
	}

	// 检查数据库是否存在
	var count int64
	db.Raw("SELECT COUNT(SCHEMA_NAME) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", conf.DBName).Scan(&count)

	// 如果数据库不存在，则创建它
	if count == 0 {
		createDB := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", conf.DBName)
		db.Exec(createDB)
	}

	// 关闭当前连接
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.Close()

	// 现在连接到新创建或验证的数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.UserName,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DBName,
	)
	GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 验证数据库连接是否成功
	DB, err = GormDB.DB()
	if err != nil {
		return err
	}
	return DB.Ping()
}

// Close 关闭数据库连接
func Close() error {
	return DB.Close()
}
