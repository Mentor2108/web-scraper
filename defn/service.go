package defn

import "context"

type Service interface {
	Create(ctx context.Context, data ...interface{}) (map[string]interface{}, error)
}
