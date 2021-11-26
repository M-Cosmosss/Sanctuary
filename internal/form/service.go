package form

type NewService struct {
	Name    string `validate:"required,lt=255"`
	Health  []string
	GroupID uint `validate:"required"`
}

type NewServiceGroup struct {
	Name   string   `validate:"required,lt=255"`
	Plugin []string `validate:"required,lt=255"`
}

type NewServiceNode struct {
	Url       string `validate:"required,lt=255"`
	ServiceID uint   `validate:"required"`
}

type GetServiceOption struct {
	OrderBy  string
	Page     int
	PageSize int
}
