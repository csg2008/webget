package util

import (
	"io/ioutil"
	"strings"
)

// SafeFileName replace all illegal chars to a underline char
func SafeFileName(fileName string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(`/\:*?"><|`, r) != -1 {
			return '_'
		}
		return r
	}, fileName)
}

// GetDirFiles 获取指定文件夹文件列表
func GetDirFiles(dirPath string, stripExt bool) []string {
	var ret []string

	if files, err := ioutil.ReadDir(dirPath); nil == err && nil != files {
		ret = make([]string, 0, len(files))
		for _, v := range files {
			if !v.IsDir() {
				if stripExt {
					var tmp = v.Name()
					var idx = strings.LastIndexByte(v.Name(), '.')
					if idx > 0 {
						ret = append(ret, strings.Trim(tmp[0:idx], " "))
					} else {
						ret = append(ret, v.Name())
					}
				} else {
					ret = append(ret, v.Name())
				}
			}
		}
	}

	return ret
}
