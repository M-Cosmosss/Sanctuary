package servicecenter

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var center ServiceCenter

func Run(c ServiceCenter) {
	center = c
	g := gin.Default()
	g.POST("/register", register)
	log.Println(g.Run(":9091"))
}

func register(c *gin.Context) {
	n := &ServiceNode{}
	if err := c.ShouldBind(n); err == nil {
		if err = center.AddServiceNode(n); err == nil {
			fmt.Printf("add node:%s success.url:%s", n.ServiceName, n.Url)
		} else {
			c.String(http.StatusBadRequest, err.Error())
		}
	} else {
		c.String(http.StatusBadRequest, "bad parameter")
	}

}
