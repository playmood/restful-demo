package conf

// Config 应用配置
// 封装为对象，用于外部对接
type Config struct {
	App   *app   `toml:"app"`
	Log   *Log   `toml:"log"`
	MySQL *MySQL `toml:"mysql"`
}

type app struct {
	Name      string `toml:"name" env:"APP_NAME"`
	Host      string `toml:"host" env:"APP_HOST"`
	Port      string `toml:"port" env:"APP_PORT"`
	Key       string `toml:"key" env:"APP_KEY"`
	EnableSSL bool   `toml:"enable_ssl" env:"APP_ENABLE_SSL"`
	CertFile  string `toml:"cert_file" env:"APP_CERT_FILE"`
	KeyFile   string `toml:"key_file" env:"APP_KEY_FILE"`
}

// Log todo
type Log struct {
	Level   string    `toml:"level" env:"LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
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
}
