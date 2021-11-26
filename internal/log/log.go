package log

import (
	"io"
	"sanctuary/internal/backend"
)

type Logger struct {
	out io.Writer
}

func Log(c *backend.Context) {

	//req,w:=c.Req,c.ResponseWriter
	//l := &db.RequestLog{
	//	IP:      req.RemoteAddr,
	//	Method:  req.Method,
	//	Path:    req.URL.String(),
	//	Success: true,
	//	ErrMsg:  "",
	//}
	////db.Mysql.Create(l)
	//println(l)
	////req.Header.Set("Host","")
	////req.Write()
	//fmt.Println("log done")
	c.Next()
}
