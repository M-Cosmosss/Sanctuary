package form

type NewRoute struct {
	Method    string `validate:"required,lt=255"`
	Path      string `validate:"required,lt=255"`
	GroupID   uint   `validate:"required"`
	ServiceID uint   `validate:"required"`
}

type NewRouteGroup struct {
	Name   string   `validate:"required,lt=255"`
	Path   string   `validate:"required,lt=255"`
	Plugin []string `validate:"lt=255"`
}

type GetRouteOption struct {
	OrderBy  string
	Page     int
	PageSize int
}
