package model

type ServerConfig struct {
	Addr string `toml:"addr"`
	Env  string `toml:"env"`
}

type DatabaseConfig struct {
	ShowSQL                                        bool `toml:"show_sql"`
	Driver, DBName, Host, Port, UserName, Password string
	MaxIdleConns                                   int `toml:"max_idle_conns"`
	MaxOpenConns                                   int `toml:"max_open_conns"`
}

type RedisConfig struct {
	URL    string `toml:"url"`
	Passwd string `toml:"passwd"`
	DB     int    `toml:"db"`
}

type LogConfig struct {
	DisableCaller bool    `toml:"disable_caller"`
	Console       console `toml:"console"`
	File          file    `toml:"file"`
}

type console struct {
	Enable bool   `toml:"enable"`
	Level  string `toml:"level"`
}

type file struct {
	Enable  bool   `toml:"enable"`
	Level   string `toml:"level"`
	Dir     string `toml:"dir"`
	MaxDays int    `toml:"max_days"`
	MaxSize int64  `toml:"max_size"`
}

type Machine struct {
	ID uint16 `toml:"id"`
}

type SMS struct {
	SecretID   string `toml:"secret_id"`
	SecretKey  string `toml:"secret_key"`
	SdkAppID   string `toml:"sdk_app_id"`
	SignName   string `toml:"sign_name"`
	TemplateID string `toml:"template_id"`
}

type JWT struct {
	IdentityKey      string `toml:"identity_key"`
	Realm            string `toml:"realm"`
	SigningAlgorithm string `toml:"signing_algorithm"`
	Key              string `toml:"key"`
	Timeout          int    `toml:"timeout"`
	MaxRefresh       int    `toml:"max_refresh"`
}

type Kafka struct {
	Server      string `toml:"server"`
	GroupID     string `toml:"group_id"`
	Topic       string `toml:"topic"`
	WorkerCount int    `toml:"worker_count"`
	PollTimeout int    `toml:"poll_timeout"`
	QueueSize   int    `toml:"queue_size"`
}

type Cache struct {
	MaxBytes int `toml:"max_bytes"`
}
