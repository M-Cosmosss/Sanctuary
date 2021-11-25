package form

type NewServiceGroup struct {
	Name   string   `validate:"required,lt=255"`
	Plugin []string `validate:"required,lt=255"`
}
