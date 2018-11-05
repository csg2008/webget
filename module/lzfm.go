package module

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/csg2008/webget/schema"
	"github.com/csg2008/webget/util"
)

// NewXLFM 荔枝FM下载链接生成器
func NewLZFM(client *util.Client) schema.Worker {
	return &LZFM{
		client: client,
	}
}

// LZFM 荔枝 FM 专辑声音下载链接生成器
type LZFM struct {
	client *util.Client
}

// Intro 显示抓取器帮助
func (s *LZFM) Intro(category string) string {
	var tip string

	switch category {
	case "label":
		tip = "荔枝 FM"
	}

	return tip
}

// Options 抓取选项
func (s *LZFM) Options() *schema.Option {
	return &schema.Option{Cli: true, Web: true, Task: false, Increment: true}
}

// Task 后台任务
func (s *LZFM) Task() error {
	return nil
}

// List 列出已经缓存的资源
func (s *LZFM) List() []map[string]string {
	return nil
}

// Search 缓存搜索
func (s *LZFM) Search(keyword string) []map[string]string {
	return nil
}

// Web 模块 web 入口
func (s *LZFM) Web(w http.ResponseWriter, req *http.Request, buf *bytes.Buffer) {

}

// Do 提取内容
func (s *LZFM) Do(tryModel bool, entry string, rule string, fp *os.File) error {
	if "" == entry {
		if tryModel {
			entry = "http://www.lizhi.fm/user/2544758401649219116"
		} else {
			return errors.New("请输入要抓取的声音专辑入口网址")
		}
	}

	if strings.Index(entry, "?") > 0 || !strings.HasPrefix(entry, "http://www.lizhi.fm/user/") {
		return errors.New("声音专辑网址格式不对，正确的格式如：http://www.lizhi.fm/user/2544758401649219116")
	}

	var cnt int
	var trackID, err = s.getItemID(entry)
	if nil == err && len(trackID) > 0 {
		for _, item := range trackID {
			if err = s.client.Download(item[1], item[0], true); nil != err {

			}
		}
	}

	if cnt > 0 && nil == err {
		err = errors.New("下载失败了 " + strconv.FormatInt(int64(cnt), 10) + " 个文件")
	}

	return err
}

// getItemID 读取专辑声音列表
func (s *LZFM) getItemID(entry string) ([][]string, error) {
	var idx int64
	var cnt int64
	var err error
	var url string
	var doc *goquery.Document
	var ret = make([][]string, 0)
	var reg = regexp.MustCompile(`^\d+$`)

	for {
		if 0 == cnt {
			idx = 1
			url = entry
		} else {
			url = entry + "/p/" + strconv.FormatInt(idx, 10) + ".html"
		}

		idx++
		doc, err = s.client.GetDoc(url, nil)
		if nil == err && nil != doc {
			// 提取分页总数
			if 0 == cnt || 0 == idx%5 {
				doc.Find("div.wrap div.frame div.page.right a").Each(func(i int, s *goquery.Selection) {
					var p = s.Text()
					if reg.MatchString(p) {
						if num, err := strconv.ParseInt(p, 10, 64); nil == err && num > cnt {
							cnt = num
						}
					}
				})
			}

			// 提取提取声音标题与ID
			doc.Find("div.wrap div.frame ul.audioList li a.js-play-data.audio-list-item").Each(func(i int, s *goquery.Selection) {
				title, _ := s.Attr("title")
				url, _ := s.Attr("data-url")
				if "" != title && "" != url {
					ret = append(ret, []string{title, url})
				}
			})
		}

		if nil != err || idx > cnt {
			break
		}
	}

	return ret, err
}
