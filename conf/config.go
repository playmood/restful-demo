package conf

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

// 全局config实例对象
// 该config对象在配置加载时初始化
// 设置为私有变量，防止恶意修改
var config *Config

// 全局mySQL客户端实例
var (
	db *sql.DB
)

// 全局config对象获取函数
func C() *Config {
	return config
}

func NewDefaultConfig() *Config {
	return &Config{
		App:   NewDefaultApp(),
		Log:   NewDefaultLog(),
		MySQL: NewDefaultMySQl(),
	}
}

// Config 应用配置
// 封装为对象，用于外部对接
type Config struct {
	App   *App   `toml:"app"`
	Log   *Log   `toml:"log"`
	MySQL *MySQL `toml:"mysql"`
}

func NewDefaultApp() *App {
	return &App{
		Name: "demo",
		Host: "127.0.0.1",
		Port: "8050",
	}
}

type App struct {
	Name string `toml:"name" env:"APP_NAME"`
	Host string `toml:"host" env:"APP_HOST"`
	Port string `toml:"port" env:"APP_PORT"`
}

func (a *App) HTTPAddr() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

func (a *App) GrpcAddr() string {
	return fmt.Sprintf("%s:%s", a.Host, fmt.Sprintf("1%s", a.Port))
}

func (a *App) RestfulAddr() string {
	return fmt.Sprintf("%s:%s", a.Host, fmt.Sprintf("2%s", a.Port))
}

func NewDefaultMySQl() *MySQL {
	return &MySQL{
		Host:        "127.0.0.1",
		Port:        "3306",
		UserName:    "demo",
		Password:    "123456",
		Database:    "demo",
		MaxOpenConn: 200,
		MaxIdleConn: 100,
	}
}

// MySQL todo
type MySQL struct {
	Host     string `toml:"host" env:"MYSQL_HOST"`
	Port     string `toml:"port" env:"MYSQL_PORT"`
	UserName string `toml:"username" env:"MYSQL_USERNAME"`
	Password string `toml:"password" env:"MYSQL_PASSWORD"`
	Database string `toml:"database" env:"MYSQL_DATABASE"`
	// 配置MySQL连接池，控制程序MySQL打开连接数
	MaxOpenConn int `toml:"max_open_conn" env:"MYSQL_MAX_OPEN_CONN"`
	// 控制MySQL复用，最大复用数
	MaxIdleConn int `toml:"max_idle_conn" env:"MYSQL_MAX_IDLE_CONN"`
	// 一个连接的生命周期，定时重启，保证可用性
	MaxLifeTime int `toml:"max_life_time" env:"MYSQL_MAX_LIFE_TIME"`
	// Idle连接最多允许存活多久
	MaxIdleTime int `toml:"max_idle_time" env:"MYSQL_MAX_idle_TIME"`
	// mutex防止getDB资源竞争，私有变量
	lock sync.Mutex
}

// 懒汉模式，获取资源时初始化
func (m *MySQL) GetDB() *sql.DB {
	// 直接加锁 锁住临界区
	m.lock.Lock()
	defer m.lock.Unlock()
	if db == nil {
		// 如果实例不存在，就初始化一个新的实例
		conn, err := m.getDBConn()
		if err != nil {
			panic(err)
		}
		db = conn
	}
	// 一定合法
	return db
}

// getDBConn use to get db connection pool
func (m *MySQL) getDBConn() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", m.UserName, m.Password, m.Host, m.Port, m.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	db.SetMaxOpenConns(m.MaxOpenConn)
	db.SetMaxIdleConns(m.MaxIdleConn)
	db.SetConnMaxLifetime(time.Second * time.Duration(m.MaxLifeTime))
	db.SetConnMaxIdleTime(time.Second * time.Duration(m.MaxIdleTime))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}
	return db, nil
}

func NewDefaultLog() *Log {
	return &Log{
		// debug, info, error, warn
		Level:  "info",
		Format: TextFormat,
		To:     ToStdout,
	}
}

// Log todo
type Log struct {
	Level   string    `toml:"level" env:"LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
}
