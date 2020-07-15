package flag

import (
	config "cash-server/configs"
	db "cash-server/database"

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
	mod        string = config.GetGlobalConfig().RunMode
)

func init() {

	flag.BoolVar(&h, "h", false, "This help")
	flag.BoolVar(&v, "v", false, "Show version and exit")
	flag.BoolVar(&t, "t", false, "Test mode")
	flag.Usage = usage
}

//StartFlag  初始化Flag
func StartFlag() {
	config.LoadGlobalConfig("")
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

	util.Info("[MODE] " + config.GetGlobalConfig().RunMode)
	if mod == "release" {
		// release mode

	} else {
		// test mode
	}
	util.Info(config.MySQL.DSN(config.GetGlobalConfig().MySQL))
	db.Dbcannot(config.MySQL.DSN(config.GetGlobalConfig().MySQL))
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
