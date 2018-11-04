package schema

import (
	"os"

	"github.com/csg2008/webget/util"
)

// WorkerHandle 工作处理接口
type WorkerHandle func(client *util.Client) Worker

// Worker 内容抓取工作者
type Worker interface {
	EnableIncrement() bool
	Help(detail bool) string
	List() []map[string]string
	Search(keyword string) []map[string]string
	Do(tryModel bool, entry string, fp *os.File) error
}
