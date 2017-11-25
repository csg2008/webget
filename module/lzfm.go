package module

import (
	"errors"
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

// Help 显示抓取器帮助
func (s *LZFM) Help(detail bool) string {
	var tip string

	if detail {
		tip = "荔枝FM 专辑下载器"
	} else {
		tip = "荔枝FM 专辑下载器"
	}

	return tip
}

// Do 提取内容
func (s *LZFM) Do(tryModel bool, entry string, fp *os.File) error {
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

	var trackID, err = s.getItemID(entry)
	if nil == err && len(trackID) > 0 {
		for _, item := range trackID {
			s.client.Download(item[1], item[0], true)
		}
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
			if 0 == cnt {
				doc.Find("div.wrap div.frame div.page.right a").Each(func(i int, s *goquery.Selection) {
					var p = s.Text()
					if reg.MatchString(p) {
						if num, err := strconv.ParseInt(p, 10, 64); nil == err {
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