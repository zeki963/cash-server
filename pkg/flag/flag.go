package flag

import (
	"cash-server/configs"
	"cash-server/db"

	"cash-server/pkg/util"
	"flag"
	"fmt"
	"os"
)

var (
	h          bool
	v          bool
	t          bool
	s          string
	appversion = "v1"
	modelFile  string
	mod        string = configs.GetGlobalConfig().RunMode
)

func init() {
	flag.BoolVar(&h, "h", false, "This help")
	flag.BoolVar(&v, "v", false, "Show version and exit")
	flag.BoolVar(&t, "t", false, "Test mode")
	flag.Usage = usage
}

//StartFlag  初始化Flag
func StartFlag() {
	configs.LoadGlobalConfig("")
	flag.Parse()
	if h {
		flag.Usage()
	}
	if v {
		util.Info(appversion)
	}
	if t {
		mod = "test"
	}

	util.Info("[MODE] " + configs.GetGlobalConfig().RunMode)
	if mod == "release" {
		// release mode

	} else {
		// test mode
	}
	util.Info(" < - MyDB INIT - >")
	if err := db.Dbcannot(configs.MySQL.DSN(configs.GetGlobalConfig().MySQL)); err != nil {
		util.Error(err.Error())
	} else {
		util.Info(configs.GetGlobalConfig().MySQL.Host)
	}

}

//usage 說明文件
func usage() {
	fmt.Fprintf(os.Stderr, `App version: `+appversion+`
Usage: App [-hvVtTq] [-s signal]  [-m model]

Options:
`)
	flag.PrintDefaults()
	os.Exit(0)
}
