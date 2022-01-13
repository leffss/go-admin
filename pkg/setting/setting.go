package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var StartTime time.Time

type App struct {
	JwtSecret string
	JwtExpireTime time.Duration
	JwtRefreshTime int64
	JwtHeaderKey string
	JwtHeaderPrefixKey string
	JwtWsHeaderKey string
	JwtWhiteLists []string
	PageSize  int
	MaxPageSize int
	PrefixUrl string
	RuntimeRootPath string
	LogSavePath string
	LogServerName string
	LogRequestName string
	LogFileExt  string
	LogLevel string
	LogMaxSize int
	LogMaxBackups int
	LogMaxAge int
	LogCompress bool
	LogConsoleOutput bool
	LogFileOutput bool
	AesKey string
}

var AppSetting = &App{}

func GetAppSetting() *App {
	return AppSetting
}

type Server struct {
	RunMode      string
	GinLog bool
	Ginpprof bool
	Swagger bool
	Port     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	CloseTimeout time.Duration
	MaxHeaderBytes int
	SSL bool
	SSLCa string
	SSLKey string
}

var ServerSetting = &Server{}

func GetServerSetting() *Server {
	return ServerSetting
}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Port		string
	Db        string
	MaxIdleConns int
	MaxOpenConns int
	ConnMaxLifetime time.Duration
	TimeZone string
	DdFile string
}

var DatabaseSetting = &Database{}

func GetDatabaseSetting() *Database {
	return DatabaseSetting
}

type Redis struct {
	Host        []string
	MasterName  string
	Password    string
	MaxIdle     int
	MaxActive   int
	PoolSize    int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

func GetRedisSetting() *Redis {
	return RedisSetting
}

type Tasks struct {
	Group        string
	StatusExpires	int
	ResultExpires    int
	Debug     bool
	NumWorkers   int
}

var TasksSetting = &Tasks{}

func GetTasksSetting() *Tasks {
	return TasksSetting
}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	StartTime = time.Now()
	var err error
	//cfg, err = ini.Load("conf/app.ini")
	cfg, err = ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,	// 跳过无法识别的数据行
	}, "conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("tasks", TasksSetting)

	if ServerSetting.RunMode != "release" && ServerSetting.RunMode != "debug" {
		ServerSetting.RunMode = "release"
	}

	AppSetting.JwtExpireTime = AppSetting.JwtExpireTime * time.Second
	//AppSetting.JwtRefreshTime = AppSetting.JwtRefreshTime * time.Second
	//AppSetting.AesKey = fmt.Sprintf("DS@#p4e)%s", AppSetting.AesKey)
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	ServerSetting.CloseTimeout = ServerSetting.CloseTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
	DatabaseSetting.ConnMaxLifetime = DatabaseSetting.ConnMaxLifetime * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

func init()  {
	Setup()
}
