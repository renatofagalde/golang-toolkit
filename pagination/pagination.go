package pagination

import (
	"context"
	"math"
	"strconv"

	"gorm.io/gorm"
)

type ctxKey string

const (
	ctxKeyPage       ctxKey = "page"
	ctxKeyPageSize   ctxKey = "page_size"
	ctxKeyTotalItems ctxKey = "total_items"
	ctxKeyTotalPages ctxKey = "total_pages"
)

// Pagination é uma estrutura genérica que representa o resultado da paginação.
type Pagination[T any] struct {
	TotalPages int   `json:"total_pages"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	Contents   []T   `json:"contents"`
}

// Paginate aplica paginação no GORM, usando informações do context.Context.
func Paginate(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := getPage(ctx)
		pageSize := getPageSize(ctx)

		// Calcula o total de registros
		var total int64
		db.Session(&gorm.Session{}).Count(&total)

		totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
		offset := (page - 1) * pageSize

		// Insere os dados no contexto
		ctx = context.WithValue(ctx, ctxKeyPage, page)
		ctx = context.WithValue(ctx, ctxKeyPageSize, pageSize)
		ctx = context.WithValue(ctx, ctxKeyTotalItems, total)
		ctx = context.WithValue(ctx, ctxKeyTotalPages, totalPages)

		return db.Offset(offset).Limit(pageSize)
	}
}

// getPage extrai o valor de "page" do contexto (como string), ou retorna 1 por padrão.
func getPage(ctx context.Context) int {
	if v := ctx.Value(ctxKeyPage); v != nil {
		if str, ok := v.(string); ok {
			if val, err := strconv.Atoi(str); err == nil && val > 0 {
				return val
			}
		} else if val, ok := v.(int); ok && val > 0 {
			return val
		}
	}
	return 1
}

// getPageSize extrai o valor de "page_size" do contexto (como string), com validações.
func getPageSize(ctx context.Context) int {
	if v := ctx.Value(ctxKeyPageSize); v != nil {
		if str, ok := v.(string); ok {
			if val, err := strconv.Atoi(str); err == nil {
				return normalizePageSize(val)
			}
		} else if val, ok := v.(int); ok {
			return normalizePageSize(val)
		}
	}
	return 100
}

// normalizePageSize limita o page_size entre 10 e 100.
func normalizePageSize(ps int) int {
	switch {
	case ps > 100:
		return 100
	case ps <= 0:
		return 10
	default:
		return ps
	}
}

// Funções utilitárias para recuperar valores do context pós-paginação

func GetPage(ctx context.Context) int {
	if v := ctx.Value(ctxKeyPage); v != nil {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 1
}

func GetPageSize(ctx context.Context) int {
	if v := ctx.Value(ctxKeyPageSize); v != nil {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 100
}

func GetTotalPages(ctx context.Context) int {
	if v := ctx.Value(ctxKeyTotalPages); v != nil {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 0
}

func GetTotalItems(ctx context.Context) int64 {
	if v := ctx.Value(ctxKeyTotalItems); v != nil {
		if i, ok := v.(int64); ok {
			return i
		}
	}
	return 0
}
