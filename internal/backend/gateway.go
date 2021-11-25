package backend

type GatewayHandler struct{}

func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{}
}

type Gateway interface {
	List() error
	NewRoute() error
	UpdateRoute() error
	DeleteRoute() error
	//NewRoute() error
	NewGroup() error
	UpdateGroup() error
	DeleteGroup() error
}

func (g *GatewayHandler) List(ctx Context) error {
	return ctx.Success()
}

func (g *GatewayHandler) NewRoute() error {
	panic("implement me")
}

func (g *GatewayHandler) UpdateRoute() error {
	panic("implement me")
}

func (g *GatewayHandler) DeleteRoute() error {
	panic("implement me")
}

func (g *GatewayHandler) NewGroup() error {
	panic("implement me")
}

func (g *GatewayHandler) UpdateGroup() error {
	panic("implement me")
}

func (g *GatewayHandler) DeleteGroup() error {
	panic("implement me")
}
