package storage

import (
	"chatbot/config"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB       *gorm.DB
	DBConfig *config.StorageConfig
)

const (
	DefaultDriveName       = "sqlite3"
	DefaultConnMaxIdleTime = 15 * time.Minute
	DefaultConnMaxLifetime = 1 * time.Hour
	DefaultMaxIdleConns    = 10
	DefaultMaxOpenConns    = 100
)

// maxRetries 是尝试重连的最大次数
const maxRetries = 5

// retryDelayMin 和 retryDelayMax 分别是重试延迟的最小和最大时间范围
const (
	retryDelayMin = 5 * time.Second
	retryDelayMax = 30 * time.Second
)

func InitWithConfig(cfg config.SqlDBConfig) *gorm.DB {
	var gormDB *gorm.DB
	var err error
	switch DBConfig.Provider {
	case "sqlite":
		gormDB, err = ConnectDB(DefaultDriveName, &cfg)
		if err != nil {
			panic(fmt.Sprintf("connect database failed:%v", err))
		}
		log.Info().Msg("connect database success")
	default:
		panic("not supported database type, please check your configuration")
	}
	return gormDB
}

func InitDB() *gorm.DB {
	if db := DB; db != nil {
		return db
	}
	DBConfig = &config.GlobalConfig.Storage
	var gormDB *gorm.DB
	var err error
	switch DBConfig.Provider {
	case "sqlite":
		gormDB, err = ConnectDB(DefaultDriveName, DBConfig.SqlDB)
		if err != nil {
			panic(fmt.Sprintf("connect database failed:%v", err))
		}
		log.Info().Msg("connect database success")
	default:
		panic("not supported database type, please check your configuration")
	}
	return gormDB
}

func ConnectDB(drive string, config *config.SqlDBConfig) (*gorm.DB, error) {
	// prepare database source
	if config.Path == "" || config.Name == "" {
		log.Error().Msg("not configured database path or name yet")
		return nil, fmt.Errorf("not configured database path or name yet")
	}
	if info, err := os.Stat(config.Path); err != nil || !info.IsDir() {
		err := os.Mkdir(config.Path, 0766)
		if err != nil {
			panic(err)
		}
	}
	source := filepath.Join(config.Path, config.Name)
	err := CreateSqlFile(source)
	if err != nil {
		log.Error().Msg("no database file yet")
		return nil, fmt.Errorf("no database file yet: %s", source)
	}
	// connect database
	db, err := sql.Open(drive, source) // 尝试打开数据库连接
	if err != nil {                    // 检查是否成功打开数据库
		log.Error().Err(err).Msg("open database failed")
		db.Close() // 确保函数退出时关闭数据库连接
		return nil, err
	}
	// Set database connection parameters
	db.SetConnMaxIdleTime(DefaultConnMaxIdleTime)
	db.SetConnMaxLifetime(DefaultConnMaxLifetime)
	db.SetMaxIdleConns(DefaultMaxIdleConns)
	db.SetMaxOpenConns(DefaultMaxOpenConns)

	// Check if the database connection is alive
	if err = db.Ping(); err != nil {
		log.Error().Msg("open database failed, check your database configuration")
		db.Close() // 确保函数退出时关闭数据库连接
		return nil, err
	}

	// Use the existing database connection to initialize GORM
	DB, err = gorm.Open(sqlite.Dialector{Conn: db}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Error().Err(err).Msg("open database failed")
		db.Close() // 确保函数退出时关闭数据库连接
		return nil, err
	}
	// 开启一个线程检测连接是否存活
	go isAlive()
	return DB, nil
}

// fixme 不好去测试
func isAlive() {
	for {
		timer := time.NewTimer(time.Minute * 10)
		defer timer.Stop() // 确保定时器被释放，避免资源泄漏

		select {
		case <-timer.C:
			if DB == nil {
				log.Error().Msg("database connection is not alive")
				return
			}
			// 尝试检查数据库连接
			err := checkDBConnection()
			if err != nil {
				log.Error().Err(err).Msg("database connection error")
				// 尝试重连
				if !reconnect(maxRetries) {
					log.Panic().Msg("failed to reconnect to the database after multiple attempts")
					return
				}
			}
		}
	}
}

// checkDBConnection 封装了检查数据库连接是否存活的逻辑
func checkDBConnection() error {
	// 假设这里的DB.Exec("SELECT 1")是标准的数据库查询方法
	// 在实际应用中，应该使用具体数据库库提供的方法
	if err := DB.Exec("SELECT 1").Error; err != nil {
		return err
	}
	return nil
}

// reconnect 尝试重连数据库，最多尝试maxRetries次
func reconnect(maxRetries int) bool {
	var retries int
	for retries < maxRetries {
		if err := Reconnect(DefaultDriveName, DBConfig.SqlDB); err == nil {
			// 重连成功，跳出循环
			return true
		}

		time.Sleep(randomRetryDelay())
		retries++
	}
	return false
}

// randomRetryDelay 返回一个随机的重试延迟时间，介于retryDelayMin和retryDelayMax之间
func randomRetryDelay() time.Duration {
	delay := time.Duration(rand.Intn(int(retryDelayMax-retryDelayMin))) + retryDelayMin
	return delay
}

func Reconnect(drive string, config *config.SqlDBConfig) error {
	// prepare database source
	if config.Path == "" || config.Name == "" {
		log.Error().Msg("not configured database path or name yet")
		return fmt.Errorf("not configured database path or name yet")
	}
	source := filepath.Join(config.Path, config.Name)
	err := CreateSqlFile(source)
	if err != nil {
		log.Error().Msg("no database file yet")
		return fmt.Errorf("no database file yet: %s", source)
	}
	// connect database
	db, err := sql.Open(drive, source) // 尝试打开数据库连接
	if err != nil {                    // 检查是否成功打开数据库
		log.Error().Msg("open database failed")
		db.Close() // 确保函数退出时关闭数据库连接
		return err
	}
	// Set database connection parameters
	db.SetConnMaxIdleTime(DefaultConnMaxIdleTime)
	db.SetConnMaxLifetime(DefaultConnMaxLifetime)
	db.SetMaxIdleConns(DefaultMaxIdleConns)
	db.SetMaxOpenConns(DefaultMaxOpenConns)

	// Check if the database connection is alive
	if err = db.Ping(); err != nil {
		log.Error().Err(err).Msg("open database failed")
		db.Close() // 确保函数退出时关闭数据库连接
		return err
	}

	// Use the existing database connection to initialize GORM
	DB, err = gorm.Open(sqlite.Dialector{Conn: db}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Error().Err(err).Msg("open database failed")
		db.Close() // 确保函数退出时关闭数据库连接
		return err
	}
	return nil
}

func CreateSqlFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 文件不存在
		// 创建文件
		file, err := os.Create(path)
		if err != nil {
			// 处理创建文件时的错误
			return err
		}
		defer file.Close()
		// 文件创建成功，可以在这里写入文件内容
		return nil
	} else {
		return nil
	}
}
