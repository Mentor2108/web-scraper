package defn

import "context"

type Repository interface {
	Create(ctx context.Context, data ...interface{}) (map[string]interface{}, error)
}
