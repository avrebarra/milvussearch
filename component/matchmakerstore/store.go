package matchmakerstore

import (
	"context"
)

type Store interface {
	TestFunc(ctx context.Context) (err error)
}
