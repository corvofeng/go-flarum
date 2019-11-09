package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	lean "github.com/johnzeng/leancloud-go-sdk"
	"github.com/op/go-logging"
)

var loggerName = "yiqi"

var (
	// 本地变量
	logger = logging.MustGetLogger(loggerName)
	format = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05} ▶ %{level:.1s} [%{shortfile}] %{message}%{color:reset}`,
	)
	leanClient = lean.NewClient(
		"k1CmAIbaqtreQSJHWEDcEeNS-gzGzoHsz",
		"HheYmnCuNLtIIRFvTRMuNHqk",
		"LYqVOTAztu4yiOLAh4yYKAm6")
)

// InitLogger 初始化Logger, 未来可以使用文件进行初始化.
func InitLogger() {

	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	logging.SetLevel(logging.DEBUG, loggerName)

	// switch *config.Loglvl {
	// case "INFO":
	// 	logging.SetLevel(logging.INFO, loggerName)
	// 	break

	// case "DEBUG":
	// 	break
	// default:
	// 	logger.Notice("Unkonw level flags")
	// }
}

//CheckError error check, 检查到error, 返回true
func CheckError(err error, info string) (res bool) {
	if err != nil {
		// Get the real file and line number
		_, fn, line, _ := runtime.Caller(1)
		fn = filepath.Base(fn)
		// pc, fn, line, _ := runtime.Caller(1)
		// logger.Errorf(
		// 	"[error] in %s[%s:%d] %v",
		// 	runtime.FuncForPC(pc).Name(), fn, line, err)

		logger.Errorf("[%s:%d] %s %v", fn, line, info, err)
		return true
	}
	return false
}

// GetLogger 全局使用的Logger
func GetLogger() *logging.Logger {
	return logger
}

func printJSON(m map[string]interface{}) {
	pp := func(args ...interface{}) {
		logger.Debug(args)
	}

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			pp(k, "is string", vv)
		case float64:
			pp(k, "is float", int64(vv))
		case int:
			pp(k, "is int", vv)
		case []interface{}:
			pp(k, "is an array:")
			for i, u := range vv {
				pp(i, u)
			}
		case nil:
			pp(k, "is nil", "null")
		case map[string]interface{}:
			pp(k, "is an map:")
			printJSON(vv)
		default:
			pp(k, "is of a type I don't know how to handle ", fmt.Sprintf("%T", v))
		}
	}
}
