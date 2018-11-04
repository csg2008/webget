package schema

import (
	"fmt"
	"os"

	"github.com/csg2008/webget/util"
)

// Webget 内容抓取器
type Webget struct {
	Debug     bool                    `label:"是否启用调试"`
	Started   int64                   `label:"程序启动时间"`
	Client    *util.Client            `label:"HTTP 客户端"`
	Version   map[string]string       `label:"应用版本信息"`
	Providers map[string]WorkerHandle `label:"数据服务提供者"`
}

// Help 显示帮助信息
func (w *Webget) Help() {
	for k, wh := range w.Providers {
		var wh = wh(w.Client)
		fmt.Fprintln(os.Stderr, "        ", k, " ", wh.Help(false))
	}
}

// Cli 启动内容抓取工作
func (w *Webget) Cli(provider string, entry string, filename string, showHelp bool, tryModel bool) {
	if wh, ok := w.Providers[provider]; ok {
		var worker = wh(w.Client)

		if showHelp {
			worker.Help(showHelp)
		} else {
			if worker.EnableIncrement() {
				w.Client.EnableIncrement(provider)

				defer w.Client.SaveIncrement()
			}

			w.Run(worker, entry, filename, tryModel)
		}
	} else {
		fmt.Fprintln(os.Stderr, "未知的内容抓取模块 [", provider, "] 请使用 -h 参数获取帮助信息")
	}

}

// Web 启动 HTTP 服务
func (w *Webget) Web() {

}

// Run 执行指定模块
func (w *Webget) Run(worker Worker, entry string, filename string, tryModel bool) {
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
		err = worker.Do(tryModel, entry, fp)
	}
	if nil == err {
		fmt.Println("内容抓取结束，谢谢使用")
	} else {
		fmt.Println("抓取失败，", err)
	}
}
