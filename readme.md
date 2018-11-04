# 说明
模块化的定向内容爬虫

爬虫工具参数说明
~~~ shell
使用说明:  webget  [命令选项] 其它参数
欢迎使用多模块定向内容抓取小工具
项目源码：https://github.com/csg2008/webget

命令选项:
  -e string
    	需要抓取的 URL 入口网址
  -h	显示应用帮助信息并退出
  -l	显示支持的数据抓取模块
  -o string
    	抓取的内容输出文件，为空输出到控制台
  -s	以 HTTP 服务的方式启动
  -t	使用模块默认参数试抓取
  -v	显示应用版本信息并退出
  -w string
    	指定要执行的内容抓取模块名

模块列表:
         xlfm   心理FM 专辑下载器
         lzfm   荔枝FM 专辑下载器
         cnblogs   博客园内容抓取器
         netease_stock   网易股票抓取器
         flysnow   飞雪无情个人博客内容抓取器
         xmly   喜马拉雅 FM 专辑声音下载链接生成器
~~~ 