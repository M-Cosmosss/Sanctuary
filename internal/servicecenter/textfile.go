package servicecenter

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	TaskAddServiceNode = iota
	TaskRegisterService
	TaskRegisterServiceGroup
	TaskSyncMap
)

type Task struct {
	TaskType  int
	Parameter interface{}
	Ctx       context.Context
	ErrorCh   chan error
}

type CenterOnTextFile struct {
	Path              string
	RServiceGroupsMap ServicesGroupMap
	WServiceGroupsMap ServicesGroupMap
	TaskQueue         chan *Task
	IsSync            bool
}

type TextFileOption struct {
	Path string
}

func NewCenterOnTextFile(op *TextFileOption) ServiceCenter {
	return &CenterOnTextFile{
		Path:              op.Path,
		RServiceGroupsMap: make(ServicesGroupMap),
		WServiceGroupsMap: make(ServicesGroupMap),
		TaskQueue:         make(chan *Task, 100),
	}
}

func (c *CenterOnTextFile) Run() error {
	go c.worker()
	for {
		select {
		case <-time.After(time.Millisecond * 300):
			if c.IsSync == false {
				c.TaskQueue <- &Task{
					TaskType:  TaskSyncMap,
					Parameter: nil,
					Ctx:       nil,
					ErrorCh:   nil,
				}
			}
		}
	}
	return errors.New("")
}

func (c *CenterOnTextFile) GetServiceNodes(s *Service) ([]string, error) {
	return nil, nil
}

func (c *CenterOnTextFile) Load() error {
	var err error
	b, err := os.ReadFile(c.Path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &c.WServiceGroupsMap)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &c.RServiceGroupsMap)
	if err != nil {
		return err
	}
	return nil
}

func (c *CenterOnTextFile) Store() error {
	var err error
	fd, err := os.Create(c.Path)
	if err != nil {
		return err
	}
	b, err := json.Marshal(c.RServiceGroupsMap)
	if err != nil {
		return err
	}
	_, err = fd.Write(b)
	return err
}

func (c *CenterOnTextFile) RegisterService(s *Service) error {
	t := &Task{
		TaskType:  TaskRegisterService,
		Parameter: s,
		ErrorCh:   make(chan error),
	}
	defer close(t.ErrorCh)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	t.Ctx = ctx
	defer cancel()
	c.TaskQueue <- t
	select {
	case err := <-t.ErrorCh:
		return err
	case <-ctx.Done():
		return errors.New("time out:service busy")
	}
}

func (c *CenterOnTextFile) RegisterServiceGroup(s *ServiceGroup) error {
	t := &Task{
		TaskType:  TaskRegisterServiceGroup,
		Parameter: s,
		ErrorCh:   make(chan error),
	}
	defer close(t.ErrorCh)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	t.Ctx = ctx
	defer cancel()
	c.TaskQueue <- t
	select {
	case err := <-t.ErrorCh:
		return err
	case <-ctx.Done():
		return errors.New("time out:service busy")
	}
}

func (c *CenterOnTextFile) AddServiceNode(n *ServiceNode) error {
	t := &Task{
		TaskType:  TaskAddServiceNode,
		Parameter: n,
		ErrorCh:   make(chan error),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	t.Ctx = ctx
	defer close(t.ErrorCh)
	c.TaskQueue <- t
	select {
	case err := <-t.ErrorCh:
		return err
	case <-ctx.Done():
		return errors.New("time out:service busy")
	}
}

func (c *CenterOnTextFile) addServiceNodeWorker(ctx context.Context, n *ServiceNode, ch chan error) {
	var g *ServiceGroup
	var s *Service
	var ok bool
	select {
	case <-ctx.Done():
		return
	default:
		if g, ok = c.WServiceGroupsMap[n.ServiceGroupName]; !ok {
			ch <- errors.New("unregister service group")
			return
		}
		if s, ok = g.Services[n.ServiceName]; !ok {
			ch <- errors.New("unregister service")
			return
		}
		s.Nodes = append(s.Nodes, n.Url)
		log.Printf("addServiceNode success:%s.%s %s", n.ServiceGroupName, n.ServiceName, n.Url)
		ch <- nil
	}
}

func (c *CenterOnTextFile) registerServiceGroup(ctx context.Context, s *ServiceGroup, ch chan error) {
	select {
	case <-ctx.Done():
		return
	default:
		if _, ok := c.WServiceGroupsMap[s.Name]; ok {
			ch <- errors.New("service group already exists")
			return
		}
		c.WServiceGroupsMap[s.Name] = &ServiceGroup{
			Name:     s.Name,
			Services: make(ServicesMap),
		}
		log.Printf("registerServiceGroup success:%s", s.Name)
		ch <- nil
	}
}

func (c *CenterOnTextFile) registerService(ctx context.Context, s *Service, ch chan error) {
	var g *ServiceGroup
	var ok bool
	select {
	case <-ctx.Done():
		return
	default:
		if g, ok = c.WServiceGroupsMap[s.GroupName]; !ok {
			ch <- errors.New("service group does not exist")
			return
		}
		if _, ok := g.Services[s.Name]; ok {
			ch <- errors.New("service in this group already exists")
			return
		}
		g.Services[s.Name] = &Service{
			Name:      s.Name,
			GroupName: s.GroupName,
			Nodes:     make([]string, 0),
		}
		log.Printf("registerService success:%s", s.Name)
		ch <- nil
	}
}

func (c *CenterOnTextFile) PrintAllRMapGroups() {
	for _, v := range c.RServiceGroupsMap {
		log.Printf("GroupName:%s\n", v.Name)
		for _, s := range v.Services {
			log.Printf("    ServiceName:%s\n", s.Name)
			for i, n := range s.Nodes {
				log.Printf("    Node %d:%s\n", i, n)
			}
		}
	}
}

func (c *CenterOnTextFile) worker() error {
	for {
		t := <-c.TaskQueue
		c.IsSync = false
		switch v := t.Parameter.(type) {
		case *ServiceNode:
			c.addServiceNodeWorker(t.Ctx, v, t.ErrorCh)
		case *Service:
			c.registerService(t.Ctx, v, t.ErrorCh)
		case *ServiceGroup:
			c.registerServiceGroup(t.Ctx, v, t.ErrorCh)
		case nil:
			c.syncRWMap()
		default:
			switch t.TaskType {
			case TaskSyncMap:

			}
			t.ErrorCh <- errors.New("type of task parameter incorrect")
		}

	}
}

func (c *CenterOnTextFile) syncRWMap() {
	c.RServiceGroupsMap, c.WServiceGroupsMap = c.WServiceGroupsMap, c.RServiceGroupsMap
	b, err := json.Marshal(c.RServiceGroupsMap)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &c.WServiceGroupsMap)
	if err != nil {
		panic(err)
	}
	c.IsSync = true
}
