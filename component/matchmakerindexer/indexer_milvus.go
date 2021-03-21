package matchmakerindexer

import (
	"context"

	"github.com/go-playground/validator"
)

type ConfigMilvus struct {
	Dependency string `validate:"required"`
}

type Milvus struct {
	config ConfigMilvus
}

func NewMilvus(cfg ConfigMilvus) (Indexer, error) {
	if err := validator.New().Struct(cfg); err != nil {
		return nil, err
	}
	return &Milvus{config: cfg}, nil
}

func (e *Milvus) TestFunc(ctx context.Context) (err error) {
	return
}
