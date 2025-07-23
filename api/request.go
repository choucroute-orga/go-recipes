package api

type IDParam struct {
	ID string `param:"id" validate:"required"`
}
