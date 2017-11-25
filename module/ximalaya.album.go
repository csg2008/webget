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

// NewXimalayaAlbum 喜马拉雅 FM 专辑声音下载链接生成器
func NewXimalayaAlbum(client *util.Client) schema.Worker {
	return &XimalayaAlbum{
		client: client,
	}
}

// XimalayaAlbum 喜马拉雅 FM 专辑声音下载链接生成器
type XimalayaAlbum struct {
	client *util.Client
}

// Help 显示抓取器帮助
func (s *XimalayaAlbum) Help(detail bool) string {
	var tip string

	if detail {
		tip = "喜马拉雅 FM 专辑声音下载链接生成器"
	} else {
		tip = "喜马拉雅 FM 专辑声音下载链接生成器"
	}

	return tip
}

// Do 提取内容
func (s *XimalayaAlbum) Do(tryModel bool, entry string, fp *os.File) error {
	if "" == entry {
		if tryModel {
			entry = "http://www.ximalaya.com/1000202/album/2667276/"
		} else {
			return errors.New("请输入要抓取的声音专辑入口网址")
		}
	}

	if strings.Index(entry, "?") > 0 || !strings.HasSuffix(entry, "/") {
		return errors.New("声音专辑网址格式不对，正确的格式如：http://www.ximalaya.com/1000202/album/2667276/")
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
func (s *XimalayaAlbum) getItemURL(id string) (string, error) {
	var out = make(map[string]interface{})
	var url = "http://www.ximalaya.com/tracks/" + id + ".json"

	var err = s.client.GetCodec(url, nil, "json", &out)
	if nil == err && nil != out {
		if uri, ok := out["play_path"].(string); ok && "" != uri {
			return uri, nil
		}

		return "", errors.New("声音 " + url + " 的下载网址为空，是收费内容？")
	}

	return "", err
}

// getItemID 读取专辑声音列表
func (s *XimalayaAlbum) getItemID(entry string) ([][]string, error) {
	var idx int64
	var cnt int64
	var err error
	var url string
	var tmp []string
	var doc *goquery.Document
	var ret = make([][]string, 0)
	var reg = regexp.MustCompile(`^\d+$`)

	for {
		if 0 == cnt {
			idx = 1
			url = entry
		} else {
			url = strings.TrimRight(entry, "/") + "?page=" + strconv.FormatInt(idx, 10)
		}

		idx++
		doc, err = s.client.GetDoc(url, nil)
		if nil == err && nil != doc {
			// 提取分页总数
			if 0 == cnt || 0 == idx%5 {
				doc.Find("div#mainbox.mainbox div.mainbox_wrapper div.pagingBar a.pagingBar_page").Each(func(i int, s *goquery.Selection) {
					var p = s.Text()
					if reg.MatchString(p) {
						if num, err := strconv.ParseInt(p, 10, 64); nil == err && num > cnt {
							cnt = num
						}
					}
				})
			}

			// 提取提取声音标题与ID
			doc.Find("div#mainbox.mainbox div.mainbox_wrapper div.album_soundlist div.miniPlayer3 a.title").Each(func(i int, s *goquery.Selection) {
				title := strings.Trim(s.Text(), "\n ")
				href, ok := s.Attr("href")
				if ok && "" != href {
					tmp = strings.SplitN(href, "/", 5)
					if "sound" == tmp[2] {
						ret = append(ret, []string{title, tmp[3]})
					}
				}
			})
		}

		if nil != err || idx > cnt {
			break
		}
	}

	return ret, err
}
