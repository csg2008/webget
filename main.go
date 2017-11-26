package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/csg2008/webget/module"
	"github.com/csg2008/webget/schema"
	"github.com/csg2008/webget/util"
)

var version = map[string]interface{}{}
var appName string
var buildRev string
var buildDate string
var buildVersion string

func main() {
	var webget *schema.Webget
	var entry = flag.String("e", "", "需要抓取的 URL 入口网址")
	var output = flag.String("o", "", "抓取的内容输出文件，为空输出到控制台")
	var worker = flag.String("w", "", "指定要执行的内容抓取模块名")
	var showVer = flag.Bool("v", false, "显示应用版本信息并退出")
	var showHelp = flag.Bool("h", false, "显示应用帮助信息并退出")
	var version = map[string]string{
		"rev":     buildRev,
		"date":    buildDate,
		"version": buildVersion,
	}
	var provider = map[string]schema.WorkerHandle{
		"cnblogs":       module.NewCnblogs,
		"netease_stock": module.NewNeteaseStock,
		"flysnow":       module.NewFlysnow,
		"xmly":          module.NewXimalayaAlbum,
		"xlfm":          module.NewXLFM,
		"lzfm":          module.NewLZFM,
	}

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "使用说明: ", strings.TrimRight(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0])), " [命令选项] 其它参数")
		fmt.Fprintln(os.Stderr, "欢迎使用多模块定向内容抓取小工具")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "命令选项:")

		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVer {
		fmt.Println(appName + " " + "Ver: " + buildVersion + " build: " + buildDate + " Rev:" + buildRev)
		return
	}

	var client = util.NewClient()
	if "" == *worker {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "模块列表:")
		for k, w := range provider {
			var wh = w(client)
			fmt.Fprintln(os.Stderr, "    ", k, " ", wh.Help(false))
		}
	} else {
		if w, ok := provider[*worker]; ok {
			var wh = w(client)
			webget = schema.NewWebget(wh, version)
			webget.Startup(*showHelp, *entry, *output)
		} else {
			fmt.Fprintln(os.Stderr, "未知的内容抓取模块 [", *worker, "] 请使用 -h 参数获取帮助信息")
		}
	}
}
