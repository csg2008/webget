package schema

import (
	"bytes"
	"net/http"
	"os"

	"github.com/csg2008/webget/util"
)

// WorkerHandle 工作处理接口
type WorkerHandle func(client *util.Client) Worker

// Option 抓取选项
type Option struct {
	Cli          bool     `label:"是否支持命令行模式"`
	Web          bool     `label:"是否支持 HTTP 模式"`
	Task         bool     `label:"是否有后任务运行"`
	Increment    bool     `label:"是否支持增量数据更新"`
	AutoStart    bool     `label:"是否自动启动攫取"`
	Status       uint     `label:"运行状态: 0 未开始 1 成功结束 2 需要授权，其它代码由应用自定义"`
	NotifyDomain []string `label:"接收通知的域名"`
}

// Worker 内容抓取工作者
type Worker interface {
	Task() error
	Options() *Option
	Intro(category string) string
	List() []map[string]string
	Search(keyword string) []map[string]string
	Web(w http.ResponseWriter, req *http.Request, buf *bytes.Buffer) bool
	Do(tryModel bool, entry string, rule string, fp *os.File) error
}
