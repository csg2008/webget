package module

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/csg2008/webget/schema"
	"github.com/csg2008/webget/util"
)

// NewJD 京东商城数据抓取器
func NewJD(client *util.Client) schema.Worker {
	return &JD{
		client: client,
		option: &schema.Option{Cli: true, Web: true, Task: false, Increment: true},
	}
}

// JDGoods 京东商品
type JDGoods struct {
	CategoryId   string `json:"categoryId" label:"品类ID"`
	Category     string `json:"category" label:"品类名称"`
	WName        string `json:"wname" label:"商品名称"`
	WareId       string `json:"wateId" label:"商品ID"`
	MiaoShaPrice string `json:"categoryId" label:"秒杀价"`
	ImageURL     string `json:"imageurl" label:"图片网址"`
}

type JDSnapshot struct {
	Ts    int64     `json:"ts" label:"同步时间"`
	Goods []JDGoods `json:"goods" label:"商品列表"`
}

// JD 京东商城数据抓取器
type JD struct {
	data   *JDSnapshot
	client *util.Client
	option *schema.Option
}

// Intro 显示抓取器帮助
func (s *JD) Intro(category string) string {
	var tip string

	switch category {
	case "label":
		tip = "京东商城"
	}

	return tip
}

// Options 抓取选项
func (s *JD) Options() *schema.Option {
	return s.option
}

// Task 后台任务
func (s *JD) Task() error {
	var data, err = util.FileGetContents("jd.json")
	if nil == err && nil != data {
		var snapshot = new(JDSnapshot)
		if err = json.Unmarshal(data, snapshot); nil == err {
			var ts = time.Unix(snapshot.Ts, 0).Format("2006-01-02")

			if ts == time.Now().Format("2006-01-02") {
				s.data = snapshot
				s.option.Status = 1
			}
		}
	}

	return err
}

// List 列出已经缓存的资源
func (s *JD) List() []map[string]string {
	return nil
}

// Search 缓存搜索
func (s *JD) Search(keyword string) []map[string]string {
	return nil
}

// Web 模块 web 入口, 返回 true 表示已经准备就绪
func (s *JD) Web(w http.ResponseWriter, req *http.Request, buf *bytes.Buffer) bool {
	return false
}

// Do 提取内容
func (s *JD) Do(tryModel bool, entry string, rule string, fp *os.File) error {
	if "" == entry {
		if tryModel {
			entry = "https://miaosha.jd.com/brandlist.html"
		} else {
			return errors.New("请输入要抓取的京东秒杀入口网址")
		}
	}

	if strings.Index(entry, "?") > 0 || !strings.HasPrefix(entry, "https://miaosha.jd.com/") {
		return errors.New("京东秒杀网址格式不对，正确的格式如：https://miaosha.jd.com/")
	}

	var categoryGoods []map[string]string
	var goodsList = make([]map[string]string, 0, 1000)
	var lists, err = s.getCategory(entry)
	if nil == err && len(lists) > 0 {
		for _, item := range lists {
			if categoryGoods, err = s.getCategoryList(item); nil == err {
				if len(categoryGoods) > 0 {
					goodsList = append(goodsList, categoryGoods...)
				}
			} else {
				fmt.Println(item, err)
			}
		}
	}

	if len(goodsList) > 0 {
		var data = map[string]interface{}{
			"ts":    time.Now().Unix(),
			"goods": goodsList,
		}

		var bj, _ = json.Marshal(data)
		util.FilePutContents("jd.json", bj, false)
	}

	return err
}

// getCategoryList 读取分类列表
func (s *JD) getCategory(entry string) ([][]string, error) {
	var err error
	var data []byte
	var ret [][]string
	var url = "https://ai.jd.com/index_new?app=Seckill&action=pcSeckillCategory&callback=pcSeckillCategory&_=" + strconv.FormatInt(time.Now().Unix(), 10) + "1234"

	if data, _, err = s.client.GetByte(url, nil); nil == err && nil != data {
		if cnt := len(data); cnt > 100 {
			var root = make(map[string]interface{})
			var temp = data[18:(cnt - 2)]
			if err = json.Unmarshal(temp, &root); nil == err {
				if categories, ok := root["categories"].([]interface{}); ok && len(categories) > 0 {
					ret = make([][]string, len(categories))

					for k, v := range categories {
						if category, ok := v.(map[string]interface{}); ok {
							if label, ok := category["categoryName"].(string); ok && "" != label {
								if id, ok := category["cateId"]; ok {
									ret[k] = []string{s.toString(id), label}
								}
							}
						}
					}
				}
			}
		}
	}

	return ret, err
}

// getCategoryList 读取分类列表
func (s *JD) getCategoryList(item []string) ([]map[string]string, error) {
	var err error
	var data []byte
	var ret []map[string]string
	var fields = []string{"wname", "wareId", "miaoShaPrice", "imageurl"}
	var url = "https://ai.jd.com/index_new?app=Seckill&action=pcSeckillCategoryGoods&callback=pcSeckillCategoryGoods&id=" + item[0] + "&_=" + strconv.FormatInt(time.Now().Unix(), 10) + "1234"

	if data, _, err = s.client.GetByte(url, nil); nil == err && nil != data {
		if cnt := len(data); cnt > 100 {
			var root = make(map[string]interface{})
			var temp = data[23:(cnt - 2)]
			if err = json.Unmarshal(temp, &root); nil == err {
				if goodsList, ok := root["goodsList"].([]interface{}); ok && len(goodsList) > 0 {
					ret = make([]map[string]string, len(goodsList))

					for k, v := range goodsList {
						if goods, ok := v.(map[string]interface{}); ok {
							ret[k] = map[string]string{
								"categoryId": item[0],
								"category":   item[1],
							}

							for _, field := range fields {
								if v, ok := goods[field]; ok {
									ret[k][field] = s.toString(v)
								}
							}
						}
					}
				}
			}
		}
	}

	return ret, err
}

func (s *JD) toString(in interface{}) string {
	var ret string

	if v, ok := in.(string); ok {
		ret = v
	} else if v, ok := in.(uint64); ok {
		ret = strconv.FormatUint(v, 10)
	} else if v, ok := in.(int64); ok {
		ret = strconv.FormatInt(v, 10)
	}

	return ret
}
