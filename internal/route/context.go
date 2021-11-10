package route

import (
	"math"
	"net/http"
)

// HandlersChain defines a HandlerFunc array.
type HandlersChain []Handle

var abortIndex int = math.MaxInt32 / 2

type Context struct {
	ErrMsg         string
	Handlers       HandlersChain
	mIndex         int
	ResponseWriter http.ResponseWriter
	Req            *http.Request
}

func newContext(w http.ResponseWriter,r *http.Request,h HandlersChain) *Context {
	return &Context{
		ErrMsg:         "",
		Handlers:       h,
		mIndex:         -1,
		ResponseWriter: w,
		Req:            r,
	}
}

func (c *Context) Next() {
	c.mIndex++
	for c.mIndex < len(c.Handlers) {
		c.Handlers[c.mIndex](c)
		c.mIndex++
	}
}

func (c *Context) Abort() {
	c.mIndex = abortIndex
}