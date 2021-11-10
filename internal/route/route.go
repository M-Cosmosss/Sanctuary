package route

import (
	"log"
	"net/http"
	"sanctuary/internal/db"
)

type route struct {
	//http.se
}

func (r *route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	go Neww(w,req)
}

func Neww(w http.ResponseWriter, req *http.Request){
	l := &db.RequestLog{
		IP:      req.RemoteAddr,
		Method:  req.Method,
		Path:    req.URL.String(),
		Success: true,
		ErrMsg:  "",
	}
	db.Mysql.Create(l)
	w.Write([]byte("done"))
}

func Run() {
	log.Println("start")
	http.ListenAndServe("localhost:6789",&route{})
}
