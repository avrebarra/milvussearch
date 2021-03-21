package matchmakerservice

import (
	"context"
)

type Service interface {
	TestFunc(ctx context.Context) (err error)
}
