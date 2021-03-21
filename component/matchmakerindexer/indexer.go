package matchmakerindexer

import (
	"context"
)

type Indexer interface {
	TestFunc(ctx context.Context) (err error)
}
