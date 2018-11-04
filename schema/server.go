package schema

import (
	"net/http"
	"time"
)

// Server http base spider web ui
type Server struct {
	server *http.Server
}

// NewServer new a http server
func NewServer(addr string) *Server {
	var ns = &Server{
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
	var mux = http.NewServeMux()
	mux.HandleFunc("/", s.home)

	s.server.Handler = mux

	return s.server.ListenAndServe()
}

// Stop 停止服务器
func (s *Server) Stop() error {
	return s.server.Shutdown(nil)
}

// home 爬虫首页
func (s *Server) home(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("webcome"))
}
