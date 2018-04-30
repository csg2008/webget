package module

import (
	"fmt"
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

// EnableIncrement 是否支持增量下载
func (s *NeteaseStock) EnableIncrement() bool {
	return true
}

func (s *NeteaseStock) Help(detail bool) string {
	var tip string

	if detail {
		tip = ""
	} else {
		tip = "网易股票抓取器"
	}

	return tip
}

func (s *NeteaseStock) Do(tryModel bool, entry string, fp *os.File) error {
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
