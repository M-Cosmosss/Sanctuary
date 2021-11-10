package route

import "net/http"

type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	Route    *Router
	//root     bool
}

// Use adds middleware to the group, see example code in GitHub.
func (group *RouterGroup) Use(middleware ...Handle) {
	group.Handlers = append(group.Handlers, middleware...)
	return
}

// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (group *RouterGroup) Group(relativePath string, handlers ...Handle) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		basePath: group.calculateAbsolutePath(relativePath),
		Route:    group.Route,
	}
}

func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.basePath, relativePath)
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (group *RouterGroup) GET(relativePath string, handlers ...Handle) error {
	if err := group.handle(http.MethodGet, relativePath, handlers); err != nil {
		return err
	}
	return nil
}

func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) error {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers)
	if err := group.Route.Handle(httpMethod, absolutePath, handlers); err != nil {
		return err
	}
	return nil
}
