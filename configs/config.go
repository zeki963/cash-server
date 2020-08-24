package configs

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

var (
	global *Config
	path   = "./configs/config.toml"
)

//func main(){
//	println(" < - App start - > ")
//	LoadGlobalConfig("./config.toml")
//	fmt.Println(global.MySQL.DSN())
//}

//LoadGlobalConfig 加載全局配置
func LoadGlobalConfig(fpath string) error {
	if fpath == "" {
		fpath = path
	}
	c, err := ParseConfig(fpath)
	if err != nil {
		return err
	}
	global = c
	return nil
}

//GetGlobalConfig 獲取全局配置
func GetGlobalConfig() *Config {
	if global == nil {
		return &Config{}
	}
	return global
}

//ParseConfig 解析配置文件
func ParseConfig(fpath string) (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(fpath, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

//Config 配置參數
type Config struct {
	RunMode  string   `toml:"runMode"`
	Swagger  bool     `toml:"swagger"`
	HTTP     HTTP     `toml:"http"`
	Logconf  Logconf  `toml:"logconf"`
	MySQL    MySQL    `toml:"mysql"`
	Mycard   Mycard   `toml:"mycard"`
	Postgres Postgres `toml:"Postgres"`
	Casino   Casino   `toml:"Casino"`
}

//HTTP 網路接口
type HTTP struct {
	HTTPS    bool `toml:"https"`
	HTTPPort int  `toml:"httpport"`
}

//Redis 配置參數
type Redis struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}

//Logconf 配置參數
type Logconf struct {
	Logger       bool `toml:"logger"`
	LoggerToFile bool `toml:"loggerToFile"`
	LoggerToDB   bool `toml:"loggerToDB"`
}

//MySQL 配置參數
type MySQL struct {
	Host       string `toml:"host"`
	Port       int    `toml:"port"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	DBName     string `toml:"dbName"`
	Parameters string `toml:"parameters"`
}

//Postgres 配置參數
type Postgres struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbName"`
}

//DSN Postgres-DSN 數據庫連接串
func (a Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		a.Host, a.Port, a.User, a.DBName, a.Password)
}

//DSN mysql-DSN 數據庫連接串
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

//Mycard 配置參數
type Mycard struct {
	FacServiceID string `toml:"FacServiceID"`
}

//Casino 配置參數
type Casino struct {
	Alphaip string `toml:"alphaip"`
	Betaip  string `toml:"betaip"`
	Proip   string `toml:"proip"`
}
