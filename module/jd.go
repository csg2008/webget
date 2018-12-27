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
	"sync"
	"time"

	"github.com/csg2008/webget/schema"
	"github.com/csg2008/webget/util"
)

// NewJD 京东商城数据抓取器
func NewJD(client *util.Client) schema.Worker {
	return &JD{
		client: client,
		option: &schema.Option{Cli: true, Web: true, Task: true, AutoStart: true, Increment: false, Mux: new(sync.RWMutex)},
	}
}

// JDGoods 京东商品
type JDGoods struct {
	CategoryId   string `json:"categoryId" label:"品类ID"`
	Category     string `json:"category" label:"品类名称"`
	WName        string `json:"wname" label:"商品名称"`
	WareId       string `json:"wareId" label:"商品ID"`
	MiaoShaPrice string `json:"miaoShaPrice" label:"秒杀价"`
	ImageURL     string `json:"imageurl" label:"图片网址"`
}

// Search 商品搜索
func (g *JDGoods) Search(keyword ...string) bool {
	var flag bool
	var word string
	var idx int
	var cnt = len(keyword)

	if cnt > 0 {
		for _, word = range keyword {
			if -1 != strings.Index(g.WName, word) {
				idx++
			}
		}

		if idx == cnt {
			flag = true
		}
	}

	return flag
}

// JDSnapshot 缓存数据快照
type JDSnapshot struct {
	Ts    int64     `json:"ts" label:"同步时间"`
	Goods []JDGoods `json:"goods" label:"商品列表"`
}

// JD 京东商城数据抓取器
type JD struct {
	Status int            `json:"status" label:"程序状态, 0 未初始化 1 正常 2 下载数据中"`
	data   *JDSnapshot    `json:"data" label:"秒杀数据快照"`
	client *util.Client   `json:"client" label:"http 客户端"`
	option *schema.Option `json:"option" label:"程序运行时选项"`
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
	var ts string
	var flag bool
	var cur = time.Now().Format("2006010215")

	s.option.Mux.RLock()
	if nil != s.data {
		ts = time.Unix(s.data.Ts, 0).Format("2006010215")
		if ts == cur && schema.WorkerStatusComplete == s.option.Status {
			flag = true
		}
	}

	s.option.Mux.RUnlock()

	if flag {
		return nil
	}

	s.option.Mux.Lock()

	var data, err = util.FileGetContents("jd.json")
	if nil == err && nil != data {
		var snapshot = new(JDSnapshot)
		if err = json.Unmarshal(data, snapshot); nil == err {
			ts = time.Unix(snapshot.Ts, 0).Format("2006010215")
			if ts == cur {
				flag = true
				s.data = snapshot
				s.option.Status = schema.WorkerStatusComplete
			}
		}
	}

	if !flag {
		s.Status = 2
		err = s.Do(true, "", "", nil)
	} else {
		if len(s.data.Goods) > 0 {
			s.Status = 1
		} else {
			s.Status = 0
		}
	}

	s.option.Mux.Unlock()

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

// SearchWeb 返回 HTML 形式的搜索结果
func (s *JD) SearchWeb(keyword string) string {
	var goods JDGoods
	var buffer = new(strings.Builder)
	var keys = strings.Split(keyword, " ")

	s.option.Mux.RLock()

	if nil != s.data && len(s.data.Goods) > 0 {
		for _, goods = range s.data.Goods {
			if goods.Search(keys...) {
				buffer.WriteString("<div class = 'goods' style='height: 100px; width: 700px;overflow: hidden;'><a target='_blank' href='https://item.jd.com/")
				buffer.WriteString(goods.WareId)
				buffer.WriteString(".html'><img style='float: left;display: block;width:100px;height:100px;overflow:hidden;' src='")
				buffer.WriteString(goods.ImageURL)
				buffer.WriteString("'><div style='float: left;margin-left: 10px;height: 100px; width: 550px; overflow: hidden;'>￥ ")
				buffer.WriteString(goods.MiaoShaPrice)
				buffer.WriteString("<br />")
				buffer.WriteString(goods.WName)
				buffer.WriteString("</div><div style='clean:both;'></div></a></div>")
			}
		}
	}

	s.option.Mux.RUnlock()

	return buffer.String()
}

// Web 模块 web 入口, 返回 true 表示已经准备就绪
func (s *JD) Web(w http.ResponseWriter, req *http.Request, buf *bytes.Buffer) bool {
	var status bool
	var q = req.FormValue("q")

	if 1 == s.Status {
		buf.WriteString("<script>function checkform(){var qs=document.getElementById('q');if (qs.value.length<1){alert('请输入要搜索的关键词，多个关键词之间用空格隔开');return false;}}</script>")
		buf.WriteString("<form method='post' action='")
		buf.WriteString(req.URL.String())
		buf.WriteString("'><input type='text' id='q' name='q' value='")
		buf.WriteString(q)
		buf.WriteString("' style='width:450px;' />")
		buf.WriteString("<input type='submit' value='search' onClick='javascript:checkform();' />")
		buf.WriteString("</form>")
	} else {
		buf.WriteString("正在更新数据，请稍候……")
		buf.WriteString("<script>setInterval(location.reload(), 1000);</script>")
	}

	if "" == q {
		buf.WriteString("请输入要搜索的内容")
	} else {
		var result = s.SearchWeb(q)
		if "" != result {
			status = true
			buf.WriteString(result)
		}
	}

	return status
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
	} else if v, ok := in.(float64); ok {
		ret = strconv.FormatFloat(v, 'f', 0, 64)
	} else if v, ok := in.(float32); ok {
		ret = strconv.FormatFloat(float64(v), 'f', 0, 32)
	}

	return ret
}
