package schema

import (
	"os"

	"github.com/csg2008/webget/util"
)

// WorkerHandle 工作处理接口
type WorkerHandle func(client *util.Client) Worker

// Worker 内容抓取工作者
type Worker interface {
	Help(detail bool) string
	Do(tryModel bool, entry string, fp *os.File) error
}
