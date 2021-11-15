package proxy

//func ReverseProxy(c *route.Context) {
//	url=c.url
//	req,w:=c.Req,c.ResponseWriter
//	newReq,err:=http.NewRequest(req.Method,url,req.Body)
//	if err!=nil{
//		panic(err)
//	}
//	client:=http.Client{}
//	resp,err:=client.Do(newReq)
//	if err!=nil{
//		panic(err)
//	}
//	defer resp.Body.Close()
//	w.WriteHeader(resp.StatusCode)
//	io.Copy(w,resp.Body)
//}
