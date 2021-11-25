package backend

import (
	"log"
	"net/http"

	"github.com/flamego/flamego"
	jsoniter "github.com/json-iterator/go"
)

type Context struct {
	flamego.Context
}

func (c *Context) Success(data ...interface{}) error {
	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)

	var d interface{}
	if len(data) == 1 {
		d = data[0]
	} else {
		d = ""
	}

	err := jsoniter.NewEncoder(c.ResponseWriter()).Encode(
		map[string]interface{}{
			"error": 0,
			"data":  d,
		},
	)
	if err != nil {
		log.Printf("Failed to encode: %v", err)
	}
	return nil
}

// Contexter initializes a classic context for a request.
func Contexter() flamego.Handler {
	return func(ctx flamego.Context) {
		c := Context{
			Context: ctx,
		}

		c.Map(c)
	}
}

func (c *Context) ServerError() error {
	return c.Error(http.StatusInternalServerError*100, "Internal server error")
}

func (c *Context) Error(errorCode uint, message string) error {
	statusCode := int(errorCode)

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(statusCode)

	err := jsoniter.NewEncoder(c.ResponseWriter()).Encode(
		map[string]interface{}{
			"error": errorCode,
			"msg":   message,
		},
	)
	if err != nil {
		log.Printf("Failed to encode: %v\n", err)
	}
	return nil
}
