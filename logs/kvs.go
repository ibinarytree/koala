package logs

import (
	"context"
	"sync"
)

var (
	initFields sync.Once
)

type LogField struct {
	fields sync.Map
}

type kvsIdKey struct{}

func WithFieldContext(ctx context.Context) context.Context {

	fields := &LogField{}
	return context.WithValue(ctx, kvsIdKey{}, fields)
}

func AddField(ctx context.Context, key string, val interface{}) {

	field := getFields(ctx)
	if field == nil {
		return
	}

	field.fields.Store(key, val)
}

func getFields(ctx context.Context) *LogField {
	field, ok := ctx.Value(kvsIdKey{}).(*LogField)
	if !ok {
		return nil
	}
	return field
}
