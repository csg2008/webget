package module

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/csg2008/webget/schema"
	"github.com/csg2008/webget/util"
)

func NewNeteaseStock(client *util.Client) schema.Worker {
	return &NeteaseStock{
		client: client,
	}
}

type NeteaseStock struct {
	client *util.Client
}

// Intro 显示抓取器帮助
func (s *NeteaseStock) Intro(category string) string {
	var tip string

	switch category {
	case "label":
		tip = "网易股票"
	}

	return tip
}

// Options 抓取选项
func (s *NeteaseStock) Options() *schema.Option {
	return &schema.Option{Cli: true, Web: true, Task: false, Increment: true}
}

// Task 后台任务
func (s *NeteaseStock) Task() error {
	return nil
}

// List 列出已经缓存的资源
func (s *NeteaseStock) List() []map[string]string {
	return nil
}

// Search 缓存搜索
func (s *NeteaseStock) Search(keyword string) []map[string]string {
	return nil
}

// Web 模块 web 入口
func (s *NeteaseStock) Web(w http.ResponseWriter, req *http.Request, buf *bytes.Buffer) {

}

func (s *NeteaseStock) Do(tryModel bool, entry string, rule string, fp *os.File) error {
	var url = "http://quotes.money.163.com/trade/lsjysj_300104.html"
	var doc, err = s.client.GetDoc(url, nil)
	if nil == err {
		var html []byte
		var code, _ = doc.Find("a.add_btn").Attr("data-code")
		var start, _ = doc.Find("input[name='date_start_type']").Eq(1).Attr("value")
		var end, _ = doc.Find("input[name='date_end_type']").Eq(1).Attr("value")
		var downUrl = fmt.Sprintf("http://quotes.money.163.com/service/chddata.html?code=%s&start=%s&end=%s&fields=TCLOSE;HIGH;LOW;TOPEN;LCLOSE;CHG;PCHG;TURNOVER;VOTURNOVER;VATURNOVER;TCAP;MCAP", code, strings.Replace(start, "-", "", -1), strings.Replace(end, "-", "", -1))
		html, _, err = s.client.GetByte(downUrl, nil)
		if nil == err {
			var file = code + ".csv"
			util.FilePutContents(file, html, false)
		} else {
			fmt.Println("err:", err)
		}
		fmt.Println("down url:", downUrl)
	}

	return nil
}
