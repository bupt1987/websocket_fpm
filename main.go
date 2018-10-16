package main

import (
	"net/http"
	"flag"
	"fmt"
	"os"
	"./connector"
	"./util"
	"github.com/cihub/seelog"
)

func main() {
	addr := flag.String("addr", ":9090", "http service address")
	ws := flag.String("ws_path", "/ws/", "http service address")
	docroot := flag.String("docroot", "/Users/junjie/www/rob/rob-server/www/", "php docroot")
	api := flag.String("api", "/api/index.php", "php docroot")
	socket := flag.String("socket", "/usr/local/var/run/php7.1-fpm.sock", "fpm listen socket")
	sLoggerConfig := flag.String("log", "./logger.xml", "log config file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	//init logger
	newLogger, err := seelog.LoggerFromConfigAsFile(*sLoggerConfig)
	if err != nil {
		panic(err)
	}
	defer util.PanicExit()

	seelog.ReplaceLogger(newLogger)
	defer seelog.Flush()

	connector.SetFcgi(*socket, *docroot, *api)

	//websocket 连接的客户端集合
	hub := connector.NewWsGroup()
	hub.Run()

	// websocket listen
	http.HandleFunc(*ws, func(w http.ResponseWriter, r *http.Request) {
		connector.ServeWs(hub, w, r)
	})

	seelog.Infof("ws start listen path => %v, address => %v", *ws, *addr)
	seelog.Error(http.ListenAndServe(*addr, nil))
}
