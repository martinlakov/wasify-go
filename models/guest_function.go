package models

type GuestFunction interface {
	Invoke(args ...any) GuestFunctionResult
}
