package helper

import "context"

func GetFromContext[T any](ctx context.Context, key string) (T, bool) {
	var zero T

	val := ctx.Value(key)
	if val == nil {
		return zero, false
	}

	tVal, ok := val.(T)
	return tVal, ok
}
