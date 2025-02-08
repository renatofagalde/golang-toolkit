package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"strconv"
)

type Pagination[T any] struct {
	TotalPages int   `json:"totalPages"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalItems int64 `json:"totalItems"`
	Contents   []T   `json:"contents"`
}

func Paginate(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := getPage(ctx)
		pageSize := getPageSize(ctx)

		clone := db.Session(&gorm.Session{})
		var total int64
		clone.Count(&total)

		totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
		ctx.Set("page", page)
		ctx.Set("pageSize", pageSize)
		ctx.Set("totalPages", totalPages)
		offset := (page - 1) * pageSize

		return db.Offset(offset).Limit(pageSize)
	}
}

func getPage(ctx *gin.Context) int {
	p, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || p <= 0 {
		return 1
	}
	return p
}

func getPageSize(ctx *gin.Context) int {
	ps, err := strconv.Atoi(ctx.DefaultQuery("page_size", "100"))
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
