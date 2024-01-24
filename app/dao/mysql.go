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

	// 创建数据库表
	if err := createTables(GormDB); err != nil {
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

// 创建数据库表
func createTables(db *gorm.DB) (err error) {
	// 创建 users 表
	err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id int(50) UNSIGNED NOT NULL,
            name varchar(255) NOT NULL,
            password varchar(255) NOT NULL,
            created_at datetime NOT NULL,
            updated_at datetime NOT NULL,
            deleted_at datetime DEFAULT NULL,
            PRIMARY KEY (id),
            UNIQUE KEY name (name)
        ) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';
    `).Error
	if err != nil {
		return err
	}

	// 创建 logs 表
	err = db.Exec(`
        CREATE TABLE IF NOT EXISTS logs (
            id varchar(255) NOT NULL,
            content json NOT NULL,
            created_at datetime NOT NULL,
            updated_at datetime NOT NULL,
            deleted_at datetime DEFAULT NULL,
            PRIMARY KEY (id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='日志表';
    `).Error
	if err != nil {
		return err
	}

	// 创建 user_logs 关联表
	err = db.Exec(`
        CREATE TABLE IF NOT EXISTS user_logs (
            user_id int(50) UNSIGNED NOT NULL,
            log_id varchar(255) NOT NULL,
            PRIMARY KEY (user_id, log_id),
            KEY log_id (log_id),
            CONSTRAINT user_logs_ibfk_1 FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE,
            CONSTRAINT user_logs_ibfk_2 FOREIGN KEY (log_id) REFERENCES logs (id) ON DELETE CASCADE ON UPDATE CASCADE
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户日志关联表';
    `).Error
	if err != nil {
		return err
	}

	return nil
}
