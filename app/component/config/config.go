package config

type App struct {
	Name string `yaml:"name"`
	Mode string `yaml:"mode"`
}

type Openai struct {
	BaseURL  string `yaml:"baseURL"`
	Key      string `yaml:"key"`
	Models    []string `yaml:"models"`
	Lang     string `yaml:"lang"`
	ProxyUrl string `yaml:"proxyUrl"`
}

type LogConfig struct {
	Level      string `yaml:"level"` // info error warn debug
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxSize"`    // 单文件最大尺寸 MB
	MaxAge     int    `yaml:"maxAge"`     // 最多保存天数day
	MaxBackups int    `yaml:"maxBackups"` // 最大备份文件数量
}

type Server struct {
	IP                    string `yaml:"ip"`
	Port                  int    `yaml:"port"`
	MinTime               int    `yaml:"minTime"`
	MaxConnectionIdle     int    `yaml:"maxConnectionIdle"`
	MaxConnectionAgeGrace int    `yaml:"maxConnectionAgeGrace"`
	TimeTick              int    `yaml:"timeTick"`
	Timeout               int    `yaml:"timeout"`
}

// Config is the configuration of the application.
type Config struct {
	App       *App       `yaml:"app"`
	Openai    *Openai    `yaml:"openai"`
	LogConfig *LogConfig `yaml:"log-config"`
	Server    *Server    `yaml:"server"`
}
