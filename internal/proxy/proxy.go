package proxy

import (
	"io"
	"net/http"
)

func ReverseProxy(w http.ResponseWriter, req *http.Request, dest string) {
	newReq, err := http.NewRequest(req.Method, dest, req.Body)
	if err != nil {
		panic(err)
	}
	client := http.Client{}
	resp, err := client.Do(newReq)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
