package conf

type Config struct {
	Server    Server `mapstructure:"server"`
	MachineId int64  `mapstructure:"machine_id"`
	Data      Data   `mapstructure:"data"`
	Log       Log    `mapstructure:"log"`
}

type Server struct {
	Http Http `mapstructure:"http"`
}

type Http struct {
	Addr    string `mapstructure:"addr"`
	Timeout string `mapstructure:"timeout"`
}

type Data struct {
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
}

type Database struct {
	Addr                    string `mapstructure:"addr"`
	User                    string `mapstructure:"user"`
	Password                string `mapstructure:"password"`
	DbName                  string `mapstructure:"dbname"`
	MaxIdleConn             int    `mapstructure:"max_idle_conn"`
	MaxOpenConn             int    `mapstructure:"max_open_conn"`
	MaxIdleTime             int    `mapstructure:"max_idle_time"`
	LowThresholdMillisecond int    `mapstructure:"low_threshold_millisecond"`
}

type Redis struct {
	Addr            string  `mapstructure:"addr"`
	Password        string  `mapstructure:"password"`
	Db              int     `mapstructure:"db"`
	PoolSize        int     `mapstructure:"pool_size"`
	BloomFilterSize uint    `mapstructure:"bloom_filter_size"`
	ErrorRate       float64 `mapstructure:"error_rate"`
	ReadTimeout     string  `mapstructure:"read_timeout"`
	WriteTimeout    string  `mapstructure:"write_timeout"`
}

type Log struct {
	Level      string `mapstructure:"level"`
	LogPath    string `mapstructure:"log_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	Console    bool   `mapstructure:"console"`
	Filename   string `mapstructure:"filename"`
}
