package module

import (
	"bytes"
	"errors"
	"net/http"
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
		option: &schema.Option{Cli: true, Web: false, Task: false, Increment: true},
	}
}

// XLFM 心理 FM 专辑声音下载链接生成器
type XLFM struct {
	client *util.Client
	option *schema.Option
}

// Intro 显示抓取器帮助
func (s *XLFM) Intro(category string) string {
	var tip string

	switch category {
	case "label":
		tip = "心理 FM"
	}

	return tip
}

// Options 抓取选项
func (s *XLFM) Options() *schema.Option {
	return s.option
}

// Task 后台任务
func (s *XLFM) Task() error {
	return nil
}

// List 列出已经缓存的资源
func (s *XLFM) List() []map[string]string {
	return nil
}

// Search 缓存搜索
func (s *XLFM) Search(keyword string) []map[string]string {
	return nil
}

// Web 模块 web 入口, 返回 true 表示已经准备就绪
func (s *XLFM) Web(w http.ResponseWriter, req *http.Request, buf *bytes.Buffer) bool {
	return false
}

// Do 提取内容
func (s *XLFM) Do(tryModel bool, entry string, rule string, fp *os.File) error {
	if "" == entry {
		entry = "http://fm.xinli001.com/broadcast-list"
	}

	if !strings.HasPrefix(entry, "http://fm.xinli001.com/broadcast-list") {
		return errors.New("声音专辑网址格式不对，正确的格式如：http://fm.xinli001.com/broadcast-list?p=1&page=1")
	}

	var cnt int
	var url string
	var trackID, err = s.getItemID(entry)
	if nil == err && len(trackID) > 0 {
		for _, item := range trackID {
			if url, err = s.getItemURL(item[1]); nil == err && "" != url {
				if err = s.client.Download(url, item[0], true); nil != err {
					cnt++
				}
			}
		}
	}

	if cnt > 0 && nil == err {
		err = errors.New("下载失败了 " + strconv.FormatInt(int64(cnt), 10) + " 个文件")
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
