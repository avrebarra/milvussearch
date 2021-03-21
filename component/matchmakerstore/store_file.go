package matchmakerstore

import (
	"context"

	"github.com/go-playground/validator"
)

type ConfigFile struct {
	Dependency string `validate:"required"`
}

type File struct {
	config ConfigFile
}

func NewFile(cfg ConfigFile) (Store, error) {
	if err := validator.New().Struct(cfg); err != nil {
		return nil, err
	}
	return &File{config: cfg}, nil
}

func (e *File) TestFunc(ctx context.Context) (err error) {
	return
}
