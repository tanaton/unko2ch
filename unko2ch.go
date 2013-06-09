package main

import (
	"log"
	"net/http"
	"regexp"
	"time"
)

type UnkoHandle struct {
	fs http.Handler
}

type Session struct {
	w http.ResponseWriter
	r *http.Request
}

var g_reg_thread *regexp.Regexp = regexp.MustCompile("(\\w+\\/\\d{9,10}(?:\\/(?:[\\-,\\d]+|l\\d+)?)?)")

func main() {
	server := &http.Server{
		Addr: ":80",
		Handler: &UnkoHandle{
			fs: http.FileServer(http.Dir("./public_html")),
		},
		ReadTimeout:    time.Duration(10) * time.Second,
		WriteTimeout:   time.Duration(10) * time.Second,
		MaxHeaderBytes: 1024 * 100,
	}
	log.Printf("listen start %s\n", server.Addr)
	// サーバ起動
	log.Fatal(server.ListenAndServe())
}

func (h *UnkoHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ses := &Session{
		w: w,
		r: r,
	}
	if ses.r.Method == "GET" {
		l := len(ses.r.URL.Path)
		if l > 15 && ses.r.URL.Path[:15] == "/test/read.cgi/" {
			// 転送処理
			if str := g_reg_thread.FindString(ses.r.URL.Path[15:]); str != "" {
				// 301
				ses.w.Header().Set("Location", "http://unkar.org/r/"+str)
				ses.statusCode(http.StatusMovedPermanently)
			} else {
				// 404
				ses.statusCode(http.StatusNotFound)
			}
		} else {
			// 後はファイルサーバーさんに任せる
			h.fs.ServeHTTP(ses.w, ses.r)
			ses.statusLog(0)
		}
	} else {
		// 501
		ses.statusCode(http.StatusNotImplemented)
	}
	return
}

func (ses *Session) statusCode(code int) {
	// ステータスコード出力
	ses.w.WriteHeader(code)
	ses.statusLog(code)
}

func (ses *Session) statusLog(code int) {
	log.Printf("%s - \"%s %s %s\" code:%d", ses.r.RemoteAddr, ses.r.Method, ses.r.RequestURI, ses.r.Proto, code)
}
