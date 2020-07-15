package config

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
	RunMode string `toml:"run_mode"`
	Swagger bool   `toml:"swagger"`
	HTTPS   bool   `toml:"https"`
	MySQL   MySQL  `toml:"mysql"`
}

//Redis 配置參數
type Redis struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
}

//MySQL 配置參數
type MySQL struct {
	Host       string `toml:"host"`
	Port       int    `toml:"port"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	DBName     string `toml:"db_name"`
	Parameters string `toml:"parameters"`
}

//DSN mysql-DSN 數據庫連接串
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}
