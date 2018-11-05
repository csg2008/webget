package schema

import (
	"bytes"
	"net/http"
	"time"
)

// Server http base spider web ui
type Server struct {
	webget *Webget
	server *http.Server
}

// NewServer new a http server
func NewServer(webget *Webget, addr string) *Server {
	var ns = &Server{
		webget: webget,
		server: &http.Server{
			Addr:           addr,
			Handler:        http.DefaultServeMux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}

	return ns
}

// Start 启动服务器
func (s *Server) Start() error {
	var worker Worker
	var option *Option
	var mux = http.NewServeMux()

	mux.HandleFunc("/", s.home)
	mux.HandleFunc("/notify.html", s.notify)
	mux.HandleFunc("/module.html", s.module)

	s.server.Handler = mux

	for k, v := range s.webget.Providers {
		worker = v(s.webget.Client)
		option = worker.Options()

		if option.Web {
			s.webget.Workers[k] = worker

			if option.Task && option.AutoStart {
				go s.webget.Workers[k].Task()
			}
		}
	}

	return s.server.ListenAndServe()
}

// Stop 停止服务器
func (s *Server) Stop() error {
	return s.server.Shutdown(nil)
}

// home 爬虫首页
func (s *Server) home(w http.ResponseWriter, req *http.Request) {
	var html = new(bytes.Buffer)

	html.WriteString("<html><head><title>welcome use webget</title></head><body>")

	for k, worker := range s.webget.Workers {
		if worker.Options().Web {
			html.WriteString("<div><a href='/module.html?m=")
			html.WriteString(k)
			html.WriteString("'>")
			html.WriteString(worker.Intro("label"))
			html.WriteString("</a></div>")
		}
	}

	html.WriteString("</body></html>")

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html.Bytes())
}

// notify 接收消息通知
func (s *Server) notify(w http.ResponseWriter, req *http.Request) {

}

// module 处理模块请求
func (s *Server) module(w http.ResponseWriter, req *http.Request) {
	var html = new(bytes.Buffer)
	var module = req.URL.Query().Get("m")

	html.WriteString("<html><head><title>welcome use webget</title></head><body>")

	if worker, ok := s.webget.Workers[module]; ok {
		var option = worker.Options()

		if !worker.Web(w, req, html) && 0 == option.Status {
			go worker.Task()
		}
	} else {
		html.WriteString("module [")
		html.WriteString(module)
		html.WriteString("] not exists or not support web ui")
	}

	html.WriteString("</body></html>")

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(html.Bytes())
}
