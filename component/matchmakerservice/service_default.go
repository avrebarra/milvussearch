package matchmakerservice

import (
	"context"

	"github.com/avrebarra/milvus-dating/component/matchmakerindexer"
	"github.com/avrebarra/milvus-dating/component/matchmakerstore"
	"github.com/go-playground/validator"
)

type Config struct {
	Index matchmakerindexer.Indexer `validate:"required"`
	Store matchmakerstore.Store     `validate:"required"`
}

type Default struct {
	config Config
}

func New(cfg Config) (Service, error) {
	if err := validator.New().Struct(cfg); err != nil {
		return nil, err
	}
	return &Default{config: cfg}, nil
}

func (e *Default) TestFunc(ctx context.Context) (err error) {
	return
}
