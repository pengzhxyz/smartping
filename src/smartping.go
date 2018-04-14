package main

import (
	"smartping/src/funcs"
	"smartping/src/g"
	"smartping/src/http"
	"flag"
	"fmt"
	"github.com/gy-games-libs/cron"
	"os"
	"runtime"
	"sync"
)

// Init config
var Version = "0.5.0"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	version := flag.Bool("v", false, "show version")
	flag.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	g.ParseConfig(Version)

	for _, target := range g.Cfg.Targets {
		go funcs.CreatePingTable(target)
	}
	c := cron.New()
	c.AddFunc("*/60 * * * * *", func() {
		var wg sync.WaitGroup
		for _, target := range g.Cfg.Targets {
			if target.Addr != g.Cfg.Ip {
				wg.Add(1)
				go funcs.StartPing(target, &wg)
			}
		}
		wg.Wait()
		go funcs.StartAlert()
	}, "ping")
	c.AddFunc("0 0 0 * * *", func() {
		go funcs.ClearAlertTable()
		go funcs.ClearPingTable()
	}, "mtc")
	c.Start()
	http.StartHttp()
}
