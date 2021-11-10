package proxy

import (
	"io"
	"net/http"
	"sanctuary/internal/route"
)

func ReverseProxy(url string,c *route.Context) {
	req,w:=c.Req,c.ResponseWriter
	newReq,err:=http.NewRequest(req.Method,url,req.Body)
	if err!=nil{
		panic(err)
	}
	client:=http.Client{}
	resp,err:=client.Do(newReq)
	if err!=nil{
		panic(err)
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w,resp.Body)
}