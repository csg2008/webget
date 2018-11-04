package module

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/csg2008/webget/schema"
	"github.com/csg2008/webget/util"
)

// NewFlysnow 新建飞雪无情博客内容抓取器
func NewFlysnow(client *util.Client) schema.Worker {
	return &Flysnow{
		client: client,
	}
}

// Flysnow 飞雪无情博客内容抓取器
type Flysnow struct {
	client *util.Client
}

// EnableIncrement 是否支持增量下载
func (s *Flysnow) EnableIncrement() bool {
	return true
}

// Help 显示抓取器帮助
func (s *Flysnow) Help(detail bool) string {
	var tip string

	if detail {
		tip = "飞雪无情个人博客内容抓取器"
	} else {
		tip = "飞雪无情个人博客内容抓取器"
	}

	return tip
}

// List 列出已经缓存的资源
func (s *Flysnow) List() []map[string]string {
	return nil
}

// Search 缓存搜索
func (s *Flysnow) Search(keyword string) []map[string]string {
	return nil
}

// Do 提取内容
func (s *Flysnow) Do(tryModel bool, entry string, fp *os.File) error {
	if "" == entry {
		if tryModel {
			entry = "http://www.flysnow.org/categories/Golang/"
		} else {
			return errors.New("请输入要抓取的博客专栏入口网址")
		}
	}

	if strings.Index(entry, "?") > 0 || !strings.HasSuffix(entry, "/") {
		return errors.New("博客专栏网址格式不对，正确的格式如：http://www.flysnow.org/categories/Golang/")
	}

	var title = s.getTitleURL(entry)
	var count = len(title)

	if count < 1 {
		return errors.New("专栏文章列表为空，请检查输入的网址是否正常？")
	}

	fmt.Fprintln(fp, "<html lang=\"zh-cn\"><head><meta charset=\"utf-8\"/><title>飞雪无情博客专栏</title></head><body>")
	for idx := len(title) - 1; idx >= 0; idx-- {
		doc, err := s.client.GetDoc(title[idx], nil)
		if err != nil {
			continue
		}

		doc.Find("div.content_container div.post").Each(func(i int, s *goquery.Selection) {
			var content, err = s.Html()
			if nil == err {
				content = strings.Replace(content, "href=\"/", "href=\"http://www.flysnow.org/", -1)
				fmt.Fprintln(fp, "<div class='post'>", content, "</div>")
			}
		})
	}

	fmt.Fprintln(fp, "</body></html>")

	return nil
}

// 读取 GO 分离文章 URL 链接
func (s *Flysnow) getTitleURL(entry string) []string {
	var idx int64 = 1
	var cnt int64
	var url string
	var ret = make([]string, 0)

	for {
		if 1 == idx {
			url = entry
		} else {
			url = fmt.Sprintf(entry+"page/%d/", idx)
		}

		doc, err := s.client.GetDoc(url, nil)
		if err != nil {
			continue
		}

		// 提取分页页码
		if 0 == cnt || 0 == idx%5 {
			doc.Find("div.content_container nav.page-navigator a.page-number").Each(func(i int, s *goquery.Selection) {
				var page, err = strconv.ParseInt(s.Text(), 10, 64)
				if nil == err && page > cnt {
					cnt = page
				}
			})
		}

		// 提取文章链接
		doc.Find("div.content_container div.post div.post-archive ul.listing li a").Each(func(i int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if ok && "" != href {
				ret = append(ret, "http://www.flysnow.org"+href)
			}
		})

		idx++

		if idx > cnt {
			break
		}
	}

	return ret
}
