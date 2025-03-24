package utils

import (
	"context"
	"math"
	"strconv"

	"gorm.io/gorm"
)

type ctxKey string

const (
	pageKey       ctxKey = "page"
	pageSizeKey   ctxKey = "pageSize"
	totalKey      ctxKey = "totalItems"
	totalPagesKey ctxKey = "totalPages"
)

func Paginate(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := getPage(ctx)
		pageSize := getPageSize(ctx)

		clone := db.Session(&gorm.Session{})
		var total int64
		clone.Count(&total)

		totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

		// Contexto padrão não tem Set, mas você pode usar outro mecanismo aqui.
		// Como não pode quebrar retrocompatibilidade, apenas guardamos esses valores em context (se necessário)
		ctx = context.WithValue(ctx, totalKey, total)
		ctx = context.WithValue(ctx, totalPagesKey, totalPages)
		ctx = context.WithValue(ctx, pageKey, page)
		ctx = context.WithValue(ctx, pageSizeKey, pageSize)

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func getPage(ctx context.Context) int {
	raw := getQueryValue(ctx, "page", "1")
	p, err := strconv.Atoi(raw)
	if err != nil || p <= 0 {
		return 1
	}
	return p
}

func getPageSize(ctx context.Context) int {
	raw := getQueryValue(ctx, "page_size", "100")
	ps, err := strconv.Atoi(raw)
	if err != nil {
		ps = 100
	}
	switch {
	case ps > 100:
		ps = 100
	case ps <= 0:
		ps = 10
	}
	return ps
}

func getQueryValue(ctx context.Context, key string, defaultVal string) string {
	if v := ctx.Value(ctxKey(key)); v != nil {
		if str, ok := v.(string); ok {
			return str
		}
	}
	return defaultVal
}
