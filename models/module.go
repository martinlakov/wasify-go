package models

import "context"

type Module interface {
	Memory() Memory
	Close(ctx context.Context) error
	GuestFunction(ctx context.Context, name string) GuestFunction
}
