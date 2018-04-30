package module

import (
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/csg2008/webget/schema"
	"github.com/csg2008/webget/util"
)

// NewXLFM 心理FM下载链接生成器
func NewXLFM(client *util.Client) schema.Worker {
	return &XLFM{
		client: client,
	}
}

// XLFM 心理 FM 专辑声音下载链接生成器
type XLFM struct {
	client *util.Client
}

// EnableIncrement 是否支持增量下载
func (s *XLFM) EnableIncrement() bool {
	return true
}

// Help 显示抓取器帮助
func (s *XLFM) Help(detail bool) string {
	var tip string

	if detail {
		tip = "心理FM 专辑下载器"
	} else {
		tip = "心理FM 专辑下载器"
	}

	return tip
}

// Do 提取内容
func (s *XLFM) Do(tryModel bool, entry string, fp *os.File) error {
	if "" == entry {
		entry = "http://fm.xinli001.com/broadcast-list"
	}

	if !strings.HasPrefix(entry, "http://fm.xinli001.com/broadcast-list") {
		return errors.New("声音专辑网址格式不对，正确的格式如：http://fm.xinli001.com/broadcast-list?p=1&page=1")
	}

	var url string
	var trackID, err = s.getItemID(entry)
	if nil == err && len(trackID) > 0 {
		for _, item := range trackID {
			if url, err = s.getItemURL(item[1]); nil == err && "" != url {
				s.client.Download(url, item[0], true)
			}
		}
	}

	return err
}

// getItemURL 获取声音 ID 对应的 URL
func (s *XLFM) getItemURL(id string) (string, error) {
	var out = make(map[string]interface{})
	var url = "http://fm.xinli001.com/broadcast?pk=" + id

	var err = s.client.GetCodec(url, nil, "json", &out)
	if nil == err && nil != out {
		if data, ok := out["data"].(map[string]interface{}); ok && nil != data {
			if media, ok := data["url"].(string); ok && "" != media {
				return media, nil
			}
		}
	}

	return "", err
}

// getItemID 读取专辑声音列表
func (s *XLFM) getItemID(entry string) ([][]string, error) {
	var idx int64
	var doc *goquery.Document
	var ret = make([][]string, 0)

	var raw, err = url.Parse(entry)
	if nil != err {
		return nil, err
	}

	if "" != raw.Query().Get("page") {
		idx, err = strconv.ParseInt(raw.Query().Get("page"), 10, 64)
		if nil != err {
			return nil, err
		}
	}

	var pos = strings.IndexByte(entry, '?')
	if pos > 0 {
		entry = entry[:pos]
	}
	if 0 == idx {
		idx = 1
	}

	for {
		var cur = entry + "?page=" + strconv.FormatInt(idx, 10) + "&p=" + strconv.FormatInt(idx, 10)

		// 提取提取声音标题与ID
		doc, err = s.client.GetDoc(cur, nil)
		if nil == err && nil != doc && doc.Length() > 0 {
			var cnt = 0
			doc.Find("li a.broadcast_title").Each(func(i int, s *goquery.Selection) {
				title := strings.Trim(s.Text(), "\n ")
				href, ok := s.Attr("href")
				if ok && "" != href {
					cnt++
					ret = append(ret, []string{title, strings.Trim(href, "/")})
				}
			})

			idx++
			if 0 == cnt {
				break
			}
		} else {
			break
		}
	}

	return ret, err
}
