package models

import "context"

type Runtime interface {
	Close(ctx context.Context) error
	Create(config *ModuleConfig) (Module, error)
}
