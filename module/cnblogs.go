package module

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/csg2008/webget/schema"
	"github.com/csg2008/webget/util"
)

// NewCnblogs 创建博客园内容抓取器实例
func NewCnblogs(client *util.Client) schema.Worker {
	return &Cnblogs{
		client: client,
	}
}

// Cnblogs 博客园内容抓取器
type Cnblogs struct {
	client *util.Client
}

// Intro 输出帮忙内容
func (s *Cnblogs) Intro(category string) string {
	var tip string

	switch category {
	case "label":
		tip = "博客园"
	}

	return tip
}

// Options 抓取选项
func (s *Cnblogs) Options() *schema.Option {
	return &schema.Option{Cli: true, Web: true, Task: false, Increment: true}
}

// Task 后台任务
func (s *Cnblogs) Task() error {
	return nil
}

// List 列出已经缓存的资源
func (s *Cnblogs) List() []map[string]string {
	return nil
}

// Search 缓存搜索
func (s *Cnblogs) Search(keyword string) []map[string]string {
	return nil
}

// Web 模块 web 入口
func (s *Cnblogs) Web(w http.ResponseWriter, req *http.Request, buf *bytes.Buffer) {

}

// Do 执行内容抓取
func (s *Cnblogs) Do(tryModel bool, entry string, rule string, fp *os.File) error {
	if "" == entry {
		if tryModel {
			entry = "http://www.cnblogs.com/coderfenghc/default.html"
		} else {
			return errors.New("请输入要抓取的博客专栏入口网址")
		}
	}

	if strings.Index(entry, "?") > 0 || !strings.HasSuffix(entry, ".html") {
		return errors.New("博客专栏网址格式不对，正确的格式如：http://www.cnblogs.com/coderfenghc/default.html")
	}

	var title = s.getTitleURL(entry)
	var count = len(title)

	if count < 1 {
		return errors.New("专栏文章列表为空，请检查输入的网址是否正常？")
	}

	fmt.Fprintln(fp, "<html><head><title>博客园专栏</title></head><body>")
	for idx := len(title) - 1; idx >= 0; idx-- {
		doc, err := s.client.GetDoc(title[idx], nil)
		if err != nil {
			continue
		}

		doc.Find("html body div#main div.post").Each(func(i int, s *goquery.Selection) {
			var title, _ = s.Find(".postTitle a").Html()
			var content, err = s.Find("div#cnblogs_post_body").Html()
			if nil == err {
				fmt.Fprintln(fp, "<div class='post'><h2>", title, "</h2><div>", content, "</div></div>")
			}
		})
	}

	fmt.Fprintln(fp, "</body></html>")

	return nil
}

// 读取 GO 分离文章 URL 链接
func (s *Cnblogs) getTitleURL(entry string) []string {
	var idx int64 = 1
	var url string
	var ret = make([]string, 0)

	for {
		if 1 == idx {
			url = entry
		} else {
			url = fmt.Sprintf(entry+"?page=%d", idx)
		}

		var flag bool
		var doc, err = s.client.GetDoc(url, nil)
		if err != nil {
			continue
		}

		// 提取文章链接
		doc.Find("html body div#main .postTitle a").Each(func(i int, s *goquery.Selection) {
			if href, ok := s.Attr("href"); ok && "" != href {
				ret = append(ret, href)
			}
		})

		// 处理分页
		doc.Find("div#main div.topicListFooter").Each(func(i int, s *goquery.Selection) {
			if s.Find("div.pager a").Length() > 0 || s.Find("div#nav_next_page a").Length() > 0 {
				flag = true
			}
		})
		if !flag {
			break
		}

		idx++
	}

	return ret
}
