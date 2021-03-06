package util

import (
	"runtime/debug"
	"github.com/cihub/seelog"
	"os"
	"github.com/pquerna/ffjson/ffjson"
)

func PanicExit() {
	if err := recover(); err != nil {
		seelog.Criticalf("%v\n%s\n======================================================\n", err, debug.Stack())
		os.Exit(1)
	}
}

func JsonEncode(data interface{}) []byte {
	res, err := ffjson.Marshal(data)
	if (err != nil) {
		seelog.Errorf("json_encode error: %v", err.Error())
	}
	return res
}
