package pagination

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
		page := ctx.Value(pageKey).(int)
		pageSize := ctx.Value(pageSizeKey).(int)

		clone := db.Session(&gorm.Session{})
		var total int64
		clone.Count(&total)

		totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

		if setter, ok := ctx.(interface {
			Set(string, any)
		}); ok {
			setter.Set("totalItems", total)
			setter.Set("totalPages", totalPages)
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func AddPaginationToContext(ctx context.Context, pageStr, pageSizeStr string) context.Context {
	page := parsePage(pageStr)
	pageSize := parsePageSize(pageSizeStr)

	ctx = context.WithValue(ctx, pageKey, page)
	ctx = context.WithValue(ctx, pageSizeKey, pageSize)
	return ctx
}

func parsePage(value string) int {
	p, err := strconv.Atoi(value)
	if err != nil || p <= 0 {
		return 1
	}
	return p
}

func parsePageSize(value string) int {
	ps, err := strconv.Atoi(value)
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
