package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/csg2008/webget/module"
	"github.com/csg2008/webget/schema"
	"github.com/csg2008/webget/util"
)

var appName string
var buildRev string
var buildDate string
var buildVersion string

func main() {
	var entry = flag.String("e", "", "需要抓取的 URL 入口网址")
	var output = flag.String("o", "", "抓取的内容输出文件，为空输出到控制台")
	var worker = flag.String("w", "", "指定要执行的内容抓取模块名")
	var server = flag.Bool("s", false, "以 HTTP 服务的方式启动")
	var debug = flag.Bool("d", false, "以调试模式运行")
	var tryModel = flag.Bool("t", false, "使用模块默认参数试抓取")
	var showList = flag.Bool("l", false, "显示支持的数据抓取模块")
	var showHelp = flag.Bool("h", false, "显示应用帮助信息并退出")
	var showVersion = flag.Bool("v", false, "显示应用版本信息并退出")

	var webget = &schema.Webget{
		Debug:   *debug,
		Started: time.Now().Unix(),
		Client:  util.NewClient(),
		Version: map[string]string{
			"rev":     buildRev,
			"date":    buildDate,
			"version": buildVersion,
		},
		Providers: map[string]schema.WorkerHandle{
			"cnblogs":       module.NewCnblogs,
			"netease_stock": module.NewNeteaseStock,
			"flysnow":       module.NewFlysnow,
			"xmly":          module.NewXimalayaAlbum,
			"xlfm":          module.NewXLFM,
			"lzfm":          module.NewLZFM,
		},
	}

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "使用说明: ", strings.TrimRight(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0])), " [命令选项] 其它参数")
		fmt.Fprintln(os.Stderr, "欢迎使用多模块定向内容抓取小工具")
		fmt.Fprintln(os.Stderr, "项目源码：https://github.com/csg2008/webget")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "命令选项:")

		flag.PrintDefaults()

		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "模块列表:")

		webget.Help()
	}

	flag.Parse()

	if *showList {
		fmt.Fprintln(os.Stderr, "模块列表:")

		webget.Help()

		return
	}
	if *showVersion {
		fmt.Println(appName + " " + "Ver: " + buildVersion + " build: " + buildDate + " Rev:" + buildRev)
		return
	}

	if *server {
		webget.Web()
	} else if "" == *worker {
		flag.Usage()
	} else {
		webget.Cli(*worker, *entry, *output, *showHelp, *tryModel)
	}
}
