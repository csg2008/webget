package schema

import (
	"os"

	"github.com/csg2008/webget/util"
)

// WorkerHandle 工作处理接口
type WorkerHandle func(client *util.Client) Worker

// Option 抓取选项
type Option struct {
	Cli       bool `label:"是否支持命令行模式"`
	Web       bool `label:"是否支持 HTTP 模式"`
	Task      bool `label:"是否有后任务运行"`
	Increment bool `label:"是否支持增量数据更新"`
}

// Worker 内容抓取工作者
type Worker interface {
	Task() error
	Options() *Option
	Help(detail bool) string
	List() []map[string]string
	Search(keyword string) []map[string]string
	Do(tryModel bool, entry string, rule string, fp *os.File) error
}
