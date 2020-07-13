package flag

import (
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

	modelFile string
	testMod   bool = false
)

func init() {
	flag.BoolVar(&h, "h", false, "This help")
	flag.BoolVar(&v, "v", false, "Show version and exit")
	flag.BoolVar(&t, "t", false, "Test mode")
	flag.Usage = usage
}

//StartFlag  初始化Flag
func StartFlag() {
	flag.Parse()
	if h {
		flag.Usage()
	}
	if v {
		util.Info(appversion)
	}
	if t {
		testMod = true
	}

	if testMod == false {
		// db.Dbcannot(dblink)
	} else {
		util.Info(" Test Mode ")

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

//SettestMod 設定
func SettestMod(tf bool) {
	testMod = tf
}
