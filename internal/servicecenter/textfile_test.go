package servicecenter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var c ServiceCenter

func init() {
	c = NewCenterOnTextFile(&TextFileOption{Path: "./testServices.json"})
}

func TestSync(t *testing.T) {
	var err error
	go c.Run()
	err = c.RegisterServiceGroup(&ServiceGroup{
		Name:     "Group1",
		Services: make(ServicesMap),
	})
	assert.Nil(t, err)
	err = c.RegisterService(&Service{
		Name:      "service1",
		GroupName: "Group1",
	})
	assert.Nil(t, err)
	err = c.AddServiceNode(&ServiceNode{
		Url:              "testurl.com:1234",
		ServiceName:      "service1",
		ServiceGroupName: "Group1",
	})
	assert.Nil(t, err)
	time.Sleep(time.Second)
	c.PrintAllRMapGroups()
	err = c.RegisterServiceGroup(&ServiceGroup{
		Name:     "Group2",
		Services: make(ServicesMap),
	})
	assert.Nil(t, err)
	err = c.RegisterService(&Service{
		Name:      "service2",
		GroupName: "Group2",
	})
	assert.Nil(t, err)
	err = c.AddServiceNode(&ServiceNode{
		Url:              "testurl2.com:1234",
		ServiceName:      "service2",
		ServiceGroupName: "Group2",
	})
	assert.Nil(t, err)
	time.Sleep(time.Second)
	c.PrintAllRMapGroups()
}
