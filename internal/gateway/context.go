package gateway

import (
	"math"
	"net/http"
)

type Handler func(ctx *Context)

// HandlersChain defines a HandlerFunc array.
type HandlersChain []Handler

var abortIndex int = math.MaxInt32 / 2

type Context struct {
	ErrMsg         string
	Handlers       HandlersChain
	mIndex         int
	ResponseWriter http.ResponseWriter
	Req            *http.Request
	RouteGroupID   uint
	ServiceID      uint
}

func newContext(w http.ResponseWriter, r *http.Request, h HandlersChain, route *Route) *Context {
	return &Context{
		ErrMsg:         "",
		Handlers:       h,
		mIndex:         -1,
		ResponseWriter: w,
		Req:            r,
		RouteGroupID:   route.RouteGroupID,
		ServiceID:      route.ServiceID,
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
