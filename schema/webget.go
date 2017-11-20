package schema

import (
	"flag"
	"fmt"
	"os"
)

var tryModel = flag.Bool("t", false, "使用模块默认参数试抓取")

// NewWebget 新建内容抓取器
func NewWebget(worker Worker, ver map[string]string) *Webget {
	return &Webget{
		worker:  worker,
		version: ver,
	}
}

// Webget 内容抓取器
type Webget struct {
	started int64             `label:"程序启动时间"`
	worker  Worker            `label:"内容抓取工具模块"`
	version map[string]string `label:"应用版本信息"`
}

// Startup 启动内容抓取工作
func (w *Webget) Startup(showDetailHelp bool, entry string, filename string) {
	if showDetailHelp {
		w.worker.Help(showDetailHelp)
	} else {
		var err error
		var fp *os.File
		if "" == filename {
			fp = os.Stdout
		} else {
			fp, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
			if nil == err {
				defer fp.Close()
			}
		}
		if nil == err {
			err = w.worker.Do(*tryModel, entry, fp)
		}
		if nil == err {
			fmt.Println("内容抓取结束，谢谢使用")
		} else {
			fmt.Println("抓取失败，", err)
		}
	}
}
