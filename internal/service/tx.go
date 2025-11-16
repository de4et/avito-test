package service

import "context"

type transactor interface {
	WithTx(context.Context, func(ctx context.Context) error) error
}
